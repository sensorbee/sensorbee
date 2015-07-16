package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"reflect"
	"time"
)

type inputBuffer struct {
	tuples     []*core.Tuple
	windowSize int64
	windowType parser.IntervalUnit
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
	emitterType  parser.Emitter
	emitterRules map[string]parser.IntervalAST
	emitCounters map[string]int64
	// curResults holds results of a query over the buffer.
	curResults []data.Map
	// prevResults holds results of a query over the buffer
	// in the previous execution run.
	prevResults []data.Map
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
		var tuples []*core.Tuple
		rangeValue := rel.Value
		rangeUnit := rel.Unit
		if rangeUnit == parser.Tuples {
			// we already know the required capacity of this buffer
			// if we work with absolute numbers
			tuples = make([]*core.Tuple, 0, rangeValue+1)
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
		relations:   lp.Relations,
		buffers:     buffers,
		emitterType: lp.EmitterType,
		curResults:  []data.Map{},
		prevResults: []data.Map{},
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
			// TODO maybe a slice is not the best implementation for a queue?
			bufferPtr := ep.buffers[rel.Alias]
			bufferPtr.tuples = append(bufferPtr.tuples, editTuple)
		}
	}

	return nil
}

// removeOutdatedTuplesFromBuffer removes tuples from the buffer that
// lie outside the current window as per the statement's window
// specification.
func (ep *streamRelationStreamExecutionPlan) removeOutdatedTuplesFromBuffer(curTupTime time.Time) error {
	for _, buffer := range ep.buffers {
		curBufSize := int64(len(buffer.tuples))
		if buffer.windowType == parser.Tuples { // tuple-based window
			if curBufSize > buffer.windowSize {
				// we just need to take the last `windowSize` items:
				// {a, b, c, d} => {b, c, d}
				buffer.tuples = buffer.tuples[curBufSize-buffer.windowSize : curBufSize]
			}

		} else if buffer.windowType == parser.Seconds { // time-based window
			// copy all "sufficiently new" tuples to new buffer
			// TODO avoid the reallocation here
			newBuf := make([]*core.Tuple, 0, curBufSize)
			for _, tup := range buffer.tuples {
				dur := curTupTime.Sub(tup.Timestamp)
				if dur.Seconds() <= float64(buffer.windowSize) {
					newBuf = append(newBuf, tup)
				}
			}
			buffer.tuples = newBuf
		} else {
			return fmt.Errorf("unknown window type: %+v", *buffer)
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
		// emit only new tuples
		for _, res := range ep.curResults {
			// check if this tuple is already present in the previous results
			found := false
			for _, prevRes := range ep.prevResults {
				if reflect.DeepEqual(res, prevRes) {
					// yes, it is, do not emit
					// TODO we may want to delete the found item from prevRes
					//      so that item counts are considered for "new items"
					found = true
					break
				}
			}
			if found {
				continue
			}
			// if we arrive here, `res` is not contained in prevResults
			output = append(output, res)
		}
	} else if ep.emitterType == parser.Dstream {
		// emit only old tuples
		for _, prevRes := range ep.prevResults {
			// check if this tuple is present in the current results
			found := false
			for _, res := range ep.curResults {
				if reflect.DeepEqual(res, prevRes) {
					// yes, it is, do not emit
					// TODO we may want to delete the found item from curRes
					//      so that item counts are considered for "new items",
					//      but can we do this safely with regard to the next run?
					found = true
					break
				}
			}
			if found {
				continue
			}
			// if we arrive here, `prevRes` is not contained in curResults
			output = append(output, prevRes)
		}
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
	if err := performQueryOnBuffer(); err != nil {
		return nil, err
	}

	// relation-to-stream:
	// compute new/old/all result data and return it
	return ep.computeResultTuples()

	return nil, nil
}

// processCartesianProduct computes the cartesian product and executes the
// given function on each resulting item
func (ep *streamRelationStreamExecutionPlan) processCartesianProduct(dataHolder data.Map, remainingKeys []string, processItem func(data.Map) error) error {
	if len(remainingKeys) > 0 {
		// not all buffers have been visited yet
		myKey := remainingKeys[0]
		myBuffer := ep.buffers[myKey].tuples
		rest := remainingKeys[1:]
		for _, t := range myBuffer {
			// add the data of this tuple to dataHolder and recurse
			dataHolder[myKey] = t.Data[myKey]
			setMetadata(dataHolder, myKey, t)
			if err := ep.processCartesianProduct(dataHolder, rest, processItem); err != nil {
				return err
			}
		}

	} else {
		// all tuples have been visited and we should now have the data
		// of one cartesian product item in dataHolder
		if err := processItem(dataHolder); err != nil {
			return err
		}
	}
	return nil
}
