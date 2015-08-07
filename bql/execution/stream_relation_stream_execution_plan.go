package execution

import (
	"container/list"
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type inputBuffer struct {
	tuples     []tupleWithDerivedInputRows
	windowSize int64
	windowType parser.IntervalUnit
}

type tupleWithDerivedInputRows struct {
	tuple *core.Tuple
	rows  []*inputRowWithCachedResult
}

func (i *inputBuffer) isTimeBased() bool {
	return i.windowType == parser.Seconds ||
		i.windowType == parser.Milliseconds
}

type inputRowWithCachedResult struct {
	input     *data.Map
	output    *data.Map
	groupData []data.Value
}

// streamRelationStreamExecutionPlan provides methods for
// execution plans that follow the theoretical
// "stream-to-relation", "relation-to-relation", "relation-to-stream"
// model very closely:
//
// After each tuple arrives,
// - compute the contents of the current window using the
//   specified window size/type,
// - perform a SELECT query on that data,
// - compute the data that need to be emitted by comparison with
//   the previous run's results.
type streamRelationStreamExecutionPlan struct {
	commonExecutionPlan
	// store name->alias mapping
	relations []parser.AliasedStreamWindowAST
	// buffers holds data of a single stream window, keyed by the
	// alias (!) of the respective input stream. It will be
	// updated (appended and possibly truncated) whenever
	// Process() is called with a new tuple.
	buffers map[string]*inputBuffer
	// emitter configuration
	emitterType parser.Emitter
	// curResults holds results of a query over the buffer.
	curResults []data.Map
	// prevResults holds results of a query over the buffer
	// in the previous execution run.
	prevResults []data.Map
	// prevHashesForIstream is only for ISTREAM and holds the hashes
	// of the items from the previous run so that we can compute
	// the check "is current item in previous results?" quickly
	prevHashesForIstream map[data.HashValue]int
	// prevHashesForDstream is only for DSTREAM and holds the
	// hashes of the items from the previous run in the same order
	// as the data, so we need to compute them only once and also
	// preserve order
	prevHashesForDstream []data.HashValue
	// now holds the a time at the beginning of the execution of
	// a statement
	now time.Time
	// filteredInputRows holds data that serves as the input for
	// the relation-to-relation operation
	filteredInputRows *list.List
	// filteredInputRows holds data that serves as the input for
	// the relation-to-relation operation
	filteredInputRowsBuffer *list.List
	// lastTupleBuffers stores the names of the input buffers that
	// the last tuple was appended to. this is valid after
	// `addTupleToBuffer` has returned.
	lastTupleBuffers map[string]bool
}

func newStreamRelationStreamExecutionPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (*streamRelationStreamExecutionPlan, error) {
	// prepare projection components
	projs, err := prepareProjections(lp.Projections, reg)
	if err != nil {
		return nil, err
	}
	// compute evaluator for the filter
	filter, err := prepareFilter(lp.Filter, reg)
	if err != nil {
		return nil, err
	}
	// compute evaluators for the group clause
	groupList, err := prepareGroupList(lp.GroupList, reg)
	if err != nil {
		return nil, err
	}
	// for compatibility with the old syntax, take the last RANGE
	// specification as valid for all buffers

	// initialize buffers (one per declared input relation)
	buffers := make(map[string]*inputBuffer, len(lp.Relations))
	for _, rel := range lp.Relations {
		var tuples []tupleWithDerivedInputRows
		rangeValue := rel.Value
		rangeUnit := rel.Unit
		if rangeUnit == parser.Tuples {
			// we already know the required capacity of this buffer
			// if we work with absolute numbers
			tuples = make([]tupleWithDerivedInputRows, 0, rangeValue+1)
		}
		// the alias of the relation is the key of the buffer
		buffers[rel.Alias] = &inputBuffer{
			tuples, rangeValue, rangeUnit,
		}
	}

	return &streamRelationStreamExecutionPlan{
		commonExecutionPlan: commonExecutionPlan{
			projections: projs,
			groupList:   groupList,
			filter:      filter,
		},
		relations:            lp.Relations,
		buffers:              buffers,
		emitterType:          lp.EmitterType,
		curResults:           []data.Map{},
		prevResults:          []data.Map{},
		prevHashesForIstream: map[data.HashValue]int{},
		prevHashesForDstream: []data.HashValue{},
		filteredInputRows:    list.New(),
	}, nil
}

// relationKey computes the InputName that belongs to a relation.
// For a real stream this equals the stream's name (independent of)
// the alias, but for a UDSF we need to use the same method that
// was used in topologyBuilder.
func (ep *streamRelationStreamExecutionPlan) relationKey(rel *parser.AliasedStreamWindowAST) string {
	if rel.Type == parser.ActualStream {
		return rel.Name
	} else {
		return fmt.Sprintf("%s/%s", rel.Name, rel.Alias)
	}
}

// addTupleToBuffer appends the received tuple to all internal buffers that
// are associated to the tuple's input name (more than one on self-join).
// Note that after calling this function, these buffers may hold more
// items than allowed by the window specification, so a call to
// removeOutdatedTuplesFromBuffer is necessary afterwards.
func (ep *streamRelationStreamExecutionPlan) addTupleToBuffer(t *core.Tuple) error {
	// we need to append this tuple to all buffers where the input name
	// matches the relation name, so first we count the those buffers
	// (for `FROM a AS left, a AS right`, this tuple will be
	// appended to the two buffers for `left` and `right`)
	numAppends := 0
	for _, rel := range ep.relations {
		if t.InputName == ep.relationKey(&rel) {
			numAppends += 1
		}
	}
	// if the tuple's input name didn't match any known relation,
	// something is wrong in the topology and we should return an error
	if numAppends == 0 {
		knownRelNames := make([]string, 0, len(ep.relations))
		for _, rel := range ep.relations {
			knownRelNames = append(knownRelNames, rel.Name)
		}
		return fmt.Errorf("tuple has input name '%s' set, but we "+
			"can only deal with %v", t.InputName, knownRelNames)
	}
	ep.lastTupleBuffers = make(map[string]bool, numAppends)
	for _, rel := range ep.relations {
		if t.InputName == ep.relationKey(&rel) {
			// if we have numAppends > 1 (meaning: this tuple is used in a
			// self-join) we should work with a copy, otherwise we can use
			// the original item
			editTuple := t
			if numAppends > 1 {
				editTuple = t.Copy()
			}
			// nest the data in a one-element map using the alias as the key
			editTuple.Data = data.Map{rel.Alias: editTuple.Data}
			// wrap this in a container struct
			editTupleCont := tupleWithDerivedInputRows{
				tuple: editTuple,
			}
			// TODO maybe a slice is not the best implementation for a queue?
			bufferPtr := ep.buffers[rel.Alias]
			bufferPtr.tuples = append(bufferPtr.tuples, editTupleCont)
			ep.lastTupleBuffers[rel.Alias] = true
		}
	}

	return nil
}

// removeOutdatedTuplesFromBuffer removes tuples from the buffer that
// lie outside the current window as per the statement's window
// specification.
func (ep *streamRelationStreamExecutionPlan) removeOutdatedTuplesFromBuffer(curTupTime time.Time) error {
	expiredInputRows := map[*inputRowWithCachedResult]bool{}
	for _, buffer := range ep.buffers {
		curBufSize := int64(len(buffer.tuples))
		if buffer.windowType == parser.Tuples { // tuple-based window
			if curBufSize > buffer.windowSize {
				// mark input rows that are derived from outdated
				// tuples for deletion
				for _, tupCont := range buffer.tuples[:curBufSize-buffer.windowSize] {
					for _, inputRow := range tupCont.rows {
						expiredInputRows[inputRow] = true
					}
				}
				// we just need to take the last `windowSize` items:
				// {a, b, c, d} => {b, c, d}
				buffer.tuples = buffer.tuples[curBufSize-buffer.windowSize : curBufSize]
			}

		} else if buffer.isTimeBased() {
			windowSizeSeconds := float64(buffer.windowSize)
			if buffer.windowType == parser.Milliseconds {
				windowSizeSeconds = windowSizeSeconds / 1000
			}
			// copy all "sufficiently new" tuples to new buffer
			// TODO avoid the reallocation here
			newBuf := make([]tupleWithDerivedInputRows, 0, curBufSize)
			for _, tupCont := range buffer.tuples {
				dur := curTupTime.Sub(tupCont.tuple.Timestamp)
				if dur.Seconds() <= windowSizeSeconds {
					newBuf = append(newBuf, tupCont)
				} else {
					// mark input rows that are derived from outdated
					// tuples for deletion
					for _, inputRow := range tupCont.rows {
						expiredInputRows[inputRow] = true
					}
				}
			}
			buffer.tuples = newBuf
		} else {
			return fmt.Errorf("unknown window type: %+v", *buffer)
		}
	}
	// now delete all rows marked for deletion
	var next *list.Element
	for e := ep.filteredInputRows.Front(); e != nil; e = next {
		next = e.Next()
		itemPtr := e.Value.(*inputRowWithCachedResult)
		if toDelete := expiredInputRows[itemPtr]; toDelete {
			ep.filteredInputRows.Remove(e)
		}
	}

	return nil
}

// computeResultTuples compares the results of this run's query with
// the results of the previous run's query and returns the data to
// be emitted as per the Emitter specification (Rstream = new,
// Istream = new-old, Dstream = old-new).
//
// Currently there is no support for multiplicities, i.e., if an item
// is 3 times in `new` and 1 time in `old` it will *not* be contained
// in the result set.
func (ep *streamRelationStreamExecutionPlan) computeResultTuples() ([]data.Map, error) {
	// TODO turn this into an iterator/generator pattern
	var output []data.Map
	if ep.emitterType == parser.Rstream {
		// emit all tuples
		for _, res := range ep.curResults {
			output = append(output, res)
		}

	} else if ep.emitterType == parser.Istream {
		curHashes := make(map[data.HashValue]int, len(ep.curResults))
		// emit only new tuples
		for _, res := range ep.curResults {
			hash := data.Hash(res)
			identicalRows := curHashes[hash] + 1
			curHashes[hash] = identicalRows
			// check if this tuple is already present in the previous results
			if prevRows := ep.prevHashesForIstream[hash]; identicalRows <= prevRows {
				continue
			}
			// if we arrive here, `res` is not contained in prevResults
			// as often as in curResults
			output = append(output, res)
		}
		// the hashes computed for the current items will be reused
		// in the next run
		ep.prevHashesForIstream = curHashes

	} else if ep.emitterType == parser.Dstream {
		// we only access the current items' hashes, not their values.
		// however, in the next run, we *will* have to access their values.
		curHashMap := make(map[data.HashValue]int, len(ep.curResults))
		curHashList := make([]data.HashValue, len(ep.curResults))
		for i, res := range ep.curResults {
			hash := data.Hash(res)
			identicalRows := curHashMap[hash] + 1
			curHashMap[hash] = identicalRows
			curHashList[i] = hash
		}
		// emit only old tuples
		counts := map[data.HashValue]int{}
		for i, prevHash := range ep.prevHashesForDstream {
			identicalRows := counts[prevHash] + 1
			counts[prevHash] = identicalRows
			// check if this tuple is present in the current results
			if curRows := curHashMap[prevHash]; identicalRows <= curRows {
				continue
			}
			// if we arrive here, `prevRes` is not contained in curResults
			// as often as in prevResults
			prevRes := ep.prevResults[i]
			output = append(output, prevRes)
		}
		// the hashes computed for the current items will be reused
		// in the next run (we keep them in a list instead of only
		// a map to prevent the order)
		ep.prevHashesForDstream = curHashList

	} else {
		return nil, fmt.Errorf("emitter type '%s' not implemented", ep.emitterType)
	}
	return output, nil
}

// Process takes an input tuple, a function that represents the "subclassing"
// plan's core functionality and returns a slice of Map values that correspond
// to the results of the query represented by this execution plan. Note that the
// order of items in the returned slice is undefined and cannot be relied on.
func (ep *streamRelationStreamExecutionPlan) process(input *core.Tuple, performQueryOnBuffer func() error) ([]data.Map, error) {
	ep.now = time.Now().In(time.UTC)

	// stream-to-relation:
	// updates the internal buffer with correct window data
	if err := ep.addTupleToBuffer(input); err != nil {
		return nil, err
	}
	if err := ep.removeOutdatedTuplesFromBuffer(input.Timestamp); err != nil {
		return nil, err
	}

	// relation-to-relation:
	// performs a SELECT query on buffer and writes result
	// to temporary table
	if err := ep.filterInputTuples(); err != nil {
		return nil, err
	}
	if err := performQueryOnBuffer(); err != nil {
		return nil, err
	}

	// relation-to-stream:
	// compute new/old/all result data and return it
	return ep.computeResultTuples()

	return nil, nil
}

func (ep *streamRelationStreamExecutionPlan) filterInputTuples() error {
	// we need to make a cross product of the data in all buffers,
	// combine it to get an input like
	//  {"streamA": {data}, "streamB": {data}, "streamC": {data}}
	// and evalute the filter on each of these items

	dataHolder := data.Map{}

	// we append the filtered results to a separate buffer so that
	// we avoid having to rollback our actual buffer if something fails
	ep.filteredInputRowsBuffer = list.New()

	// Note: `ep.buffers` is a map, so iterating over its keys may yield
	// different results in every run of the program. We cannot expect
	// a consistent order in which evalItem is run on the items of the
	// cartesian product.
	allStreams := make(map[string][]tupleWithDerivedInputRows, len(ep.buffers))

	// we do not reevaluate the filter for all elements in the cartesian
	// product of the input buffers, but only for those elements that use
	// the newly added tuple.
	// so we have to compute the difference between "cartesian product of
	// buffers including new tuple" and "cartesian product of buffers
	// excluding new tuple". this becomes a bit combinatorial if
	// we have a self-join. for example, if we have a join over five
	// streams, three of which are identical, then we need to compute
	//  (A∪{t})×(B∪{t})×(C∪{t})×D×E \ A×B×C×D×E
	// which looks like
	//  ({t})×(B∪{t})×(C∪{t})×D×E ∪
	//  ( A )×( {t} )×(C∪{t})×D×E ∪
	//  ( A )×(  B  )×(C∪{t})×D×E
	buffersWithNewTuple := make([]string, 0, len(ep.lastTupleBuffers))
	buffersWithoutNewTuple := make([]string, 0, len(ep.buffers)-len(ep.lastTupleBuffers))
	for key := range ep.buffers {
		if justAppended := ep.lastTupleBuffers[key]; justAppended {
			buffersWithNewTuple = append(buffersWithNewTuple, key)
		} else {
			buffersWithoutNewTuple = append(buffersWithoutNewTuple, key)
		}
	}
	// we need as many runs as there are buffers that hold the new tuple
	for i := 0; i < len(buffersWithNewTuple); i++ {
		// buffers <i are taken without the new tuple
		for j := 0; j < i; j++ {
			key := buffersWithNewTuple[j]
			buffer := ep.buffers[key]
			allStreams[key] = buffer.tuples[:len(buffer.tuples)-1]
		}
		// buffer i uses just the new tuple
		{
			j := i
			key := buffersWithNewTuple[j]
			buffer := ep.buffers[key]
			allStreams[key] = buffer.tuples[len(buffer.tuples)-1 : len(buffer.tuples)]
		}
		// buffers >i are taken including the new tuple
		for j := i + 1; j < len(buffersWithNewTuple); j++ {
			key := buffersWithNewTuple[j]
			buffer := ep.buffers[key]
			allStreams[key] = buffer.tuples
		}
		// and all buffers that do not hold the new tuple are
		// always taken completely
		for _, key := range buffersWithoutNewTuple {
			buffer := ep.buffers[key]
			allStreams[key] = buffer.tuples
		}
		// write matching items to ep.filteredInputRowsBuffer
		if err := ep.preprocessCartesianProduct(dataHolder, allStreams); err != nil {
			return err
		}
	}
	// write only the items matching the filter to ep.filteredInputRows
	// (NB. the items appended here will be cleaned up in future
	// runs by `removeOutdatedTuplesFromBuffer`)
	ep.filteredInputRows.PushBackList(ep.filteredInputRowsBuffer)
	return nil
}

// preprocessCartesianProduct computes the cartesian product,
// applies this plan's filter/join condition to each item and
// appends it to `ep.filteredInputRows`
func (ep *streamRelationStreamExecutionPlan) preprocessCartesianProduct(dataHolder data.Map, remainingBuffers map[string][]tupleWithDerivedInputRows) error {
	return ep.preprocCartProdInt(dataHolder, remainingBuffers,
		map[string]*tupleWithDerivedInputRows{})
}

func (ep *streamRelationStreamExecutionPlan) preprocCartProdInt(dataHolder data.Map, remainingBuffers map[string][]tupleWithDerivedInputRows, origin map[string]*tupleWithDerivedInputRows) error {
	if len(remainingBuffers) > 0 {
		// not all buffers have been visited yet
		var myKey string
		for key := range remainingBuffers {
			myKey = key
			break
		}
		myBuffer := remainingBuffers[myKey]
		// compile a dictionary with the rest of the unvisited streams
		// (do NOT modify remainingBuffers directly!)
		rest := map[string][]tupleWithDerivedInputRows{}
		for key, buffer := range remainingBuffers {
			if key != myKey {
				rest[key] = buffer
			}
		}
		for i, t := range myBuffer {
			// add the data of this tuple to dataHolder and recurse
			dataHolder[myKey] = t.tuple.Data[myKey]
			origin[myKey] = &myBuffer[i]
			setMetadata(dataHolder, myKey, t.tuple)
			if err := ep.preprocCartProdInt(dataHolder, rest, origin); err != nil {
				return err
			}
		}

	} else {
		// all tuples have been visited and we should now have the data
		// of one cartesian product item in dataHolder

		// add the information accessed by the now() function
		// to each item
		dataHolder[":meta:NOW"] = data.Timestamp(ep.now)

		// evaluate filter condition and convert to bool
		if ep.filter != nil {
			filterResult, err := ep.filter.Eval(dataHolder)
			if err != nil {
				return err
			}
			filterResultBool, err := data.ToBool(filterResult)
			if err != nil {
				return err
			}
			// if it evaluated to false, do not further process this tuple
			// (ToBool also evalutes the NULL value to false, so we don't
			// need to treat this specially)
			if !filterResultBool {
				return nil
			}
		}

		// if we arrive here, this item of the cartesian product fulfills
		// the filter/join condition, so we make a shallow copy (that should
		// be fine) and add it to the list of input items
		item := make(data.Map, len(dataHolder))
		for key, val := range dataHolder {
			item[key] = val
		}
		itemWithCachedResult := &inputRowWithCachedResult{
			input: &item,
		}
		// also write the address of this item to all tuples
		// it originates from
		for _, tupHolder := range origin {
			tupHolder.rows = append(tupHolder.rows, itemWithCachedResult)
		}
		ep.filteredInputRowsBuffer.PushBack(itemWithCachedResult)
	}
	return nil
}
