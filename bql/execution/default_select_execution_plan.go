package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"time"
)

type defaultSelectExecutionPlan struct {
	commonExecutionPlan
	emitterType  parser.Emitter
	emitterRules map[string]parser.IntervalAST
	emitCounters map[string]int64
	// store name->alias mapping
	relations []parser.AliasedStreamWindowAST
	// buffers holds data of a single stream window, keyed by the
	// alias (!) of the respective input stream. It will be
	// updated (appended and possibly truncated) whenever
	// Process() is called with a new tuple.
	buffers map[string]*inputBuffer
	// curResults holds results of a query over the buffer.
	curResults []tuple.Map
	// prevResults holds results of a query over the buffer
	// in the previous execution run.
	prevResults []tuple.Map
}

type inputBuffer struct {
	tuples     []*tuple.Tuple
	windowSize int64
	windowType parser.IntervalUnit
}

// CanBuildDefaultSelectExecutionPlan checks whether the given statement
// allows to use an defaultSelectExecutionPlan.
func CanBuildDefaultSelectExecutionPlan(lp *LogicalPlan, reg udf.FunctionRegistry) bool {
	// TODO check that there are no aggregate functions
	return len(lp.GroupList) == 0 &&
		lp.Having == nil
}

// defaultSelectExecutionPlan is a very simple plan that follows the
// theoretical processing model.
//
// After each tuple arrives,
// - compute the contents of the current window using the
//   specified window size/type,
// - perform a SELECT query on that data,
// - compute the data that need to be emitted by comparison with
//   the previous run's results.
func NewDefaultSelectExecutionPlan(lp *LogicalPlan, reg udf.FunctionRegistry) (ExecutionPlan, error) {
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
	// for compatibility with the old syntax, take the last RANGE
	// specification as valid for all buffers

	// initialize buffers (one per declared input relation)
	buffers := make(map[string]*inputBuffer, len(lp.Relations))
	for _, rel := range lp.Relations {
		var tuples []*tuple.Tuple
		rangeValue := rel.Value
		rangeUnit := rel.Unit
		if rangeUnit == parser.Tuples {
			// we already know the required capacity of this buffer
			// if we work with absolute numbers
			tuples = make([]*tuple.Tuple, 0, rangeValue+1)
		}
		// the alias of the relation is the key of the buffer
		buffers[rel.Alias] = &inputBuffer{
			tuples, rangeValue, rangeUnit,
		}
	}
	emitterRules := make(map[string]parser.IntervalAST, len(lp.EmitIntervals))
	emitCounters := make(map[string]int64, len(lp.EmitIntervals))
	if len(lp.EmitIntervals) == 0 {
		// set the default if not `EVERY ...` was given
		emitterRules["*"] = parser.IntervalAST{parser.NumericLiteral{1}, parser.Tuples}
		emitCounters["*"] = 0
	}
	for _, emitRule := range lp.EmitIntervals {
		// TODO implement time-based emitter as well
		if emitRule.Unit == parser.Seconds {
			return nil, fmt.Errorf("time-based emitter not implemented")
		}
		emitterRules[emitRule.Name] = emitRule.IntervalAST
		emitCounters["*"] = 0
	}
	return &defaultSelectExecutionPlan{
		commonExecutionPlan: commonExecutionPlan{
			projections: projs,
			filter:      filter,
		},
		emitterType:  lp.EmitterType,
		emitterRules: emitterRules,
		emitCounters: emitCounters,
		relations:    lp.Relations,
		buffers:      buffers,
		curResults:   []tuple.Map{},
		prevResults:  []tuple.Map{},
	}, nil
}

func (ep *defaultSelectExecutionPlan) Process(input *tuple.Tuple) ([]tuple.Map, error) {
	// stream-to-relation:
	// updates the internal buffer with correct window data
	if err := ep.addTupleToBuffer(input); err != nil {
		return nil, err
	}
	if err := ep.removeOutdatedTuplesFromBuffer(input.Timestamp); err != nil {
		return nil, err
	}

	if ep.shouldEmitNow(input) {
		// relation-to-relation:
		// performs a SELECT query on buffer and writes result
		// to temporary table
		if err := ep.performQueryOnBuffer(); err != nil {
			return nil, err
		}

		// relation-to-stream:
		// compute new/old/all result data and return it
		// TODO use an iterator/generator pattern instead
		return ep.computeResultTuples()
	}
	return nil, nil
}

// addTupleToBuffer appends the received tuple to all internal buffers that
// are associated to the tuple's input name (more than one on self-join).
// Note that after calling this function, these buffers may hold more
// items than allowed by the window specification, so a call to
// removeOutdatedTuplesFromBuffer is necessary afterwards.
func (ep *defaultSelectExecutionPlan) addTupleToBuffer(t *tuple.Tuple) error {
	// we need to append this tuple to all buffers where the input name
	// matches the relation name, so first we count the those buffers
	// (for `FROM a AS left, a AS right`, this tuple will be
	// appended to the two buffers for `left` and `right`)
	numAppends := 0
	for _, rel := range ep.relations {
		if t.InputName == rel.Name {
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
		if t.InputName == rel.Name {
			// if we have numAppends > 1 (meaning: this tuple is used in a
			// self-join) we should work with a copy, otherwise we can use
			// the original item
			editTuple := t
			if numAppends > 1 {
				editTuple = t.Copy()
			}
			// nest the data in a one-element map using the alias as the key
			editTuple.Data = tuple.Map{rel.Alias: editTuple.Data}
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
func (ep *defaultSelectExecutionPlan) removeOutdatedTuplesFromBuffer(curTupTime time.Time) error {
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
			newBuf := make([]*tuple.Tuple, 0, curBufSize)
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

// shouldEmitNow returns true if the input tuple should trigger
// computation of output values.
func (ep *defaultSelectExecutionPlan) shouldEmitNow(t *tuple.Tuple) bool {
	// first check if we have a stream-independent rule
	// (e.g., `RSTREAM` or `RSTREAM [EVERY 2 TUPLES]`)
	if interval, ok := ep.emitterRules["*"]; ok {
		counter := ep.emitCounters["*"]
		nextCounter := counter + 1
		if nextCounter%interval.Value == 0 {
			ep.emitCounters["*"] = 0
			return true
		}
		ep.emitCounters["*"] = nextCounter
		return false
	}
	// if there was no such rule, check if there is a
	// rule for the input stream the tuple came from
	if interval, ok := ep.emitterRules[t.InputName]; ok {
		counter := ep.emitCounters[t.InputName]
		nextCounter := counter + 1
		if nextCounter%interval.Value == 0 {
			ep.emitCounters[t.InputName] = 0
			return true
		}
		ep.emitCounters[t.InputName] = nextCounter
		return false
	}
	// there is no general rule and no rule for the input
	// stream of the tuple, so don't do anything
	return false
}

// performQueryOnBuffer executes a SELECT query on the data of the tuples
// currently stored in the buffer. The query results (which is a set of
// tuple.Value, not tuple.Tuple) is stored in ep.curResults. The data
// that was stored in ep.curResults before this method was called is
// moved to ep.prevResults.
//
// In case of an error the contents of ep.curResults will still be
// the same as before the call (so that the next run performs as
// if no error had happened), but the contents of ep.curResults are
// undefined.
//
// Currently performQueryOnBuffer can only perform SELECT ... WHERE ...
// queries without aggregate functions, GROUP BY, or HAVING clauses.
func (ep *defaultSelectExecutionPlan) performQueryOnBuffer() error {
	// reuse the allocated memory
	output := ep.prevResults[0:0]
	// remember the previous results
	ep.prevResults = ep.curResults

	// we need to make a cross product of the data in all buffers,
	// combine it to get an input like
	//  {"streamA": {data}, "streamB": {data}, "streamC": {data}}
	// and then run filter/projections on each of this items

	dataHolder := tuple.Map{}

	// function to compute cartesian product and do something on each
	// resulting item
	var procCartProd func([]string, func(tuple.Map) error) error

	procCartProd = func(remainingKeys []string, processItem func(tuple.Map) error) error {
		if len(remainingKeys) > 0 {
			// not all buffers have been visited yet
			myKey := remainingKeys[0]
			myBuffer := ep.buffers[myKey].tuples
			rest := remainingKeys[1:]
			for _, t := range myBuffer {
				// add the data of this tuple to dataHolder and recurse
				dataHolder[myKey] = t.Data[myKey]
				if err := procCartProd(rest, processItem); err != nil {
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

	// function to evaluate filter on the input data and -- if the filter does
	// not exist or evaluates to true -- compute the projections and store
	// the result in the `output` slice
	evalItem := func(data tuple.Map) error {
		// evaluate filter condition and convert to bool
		if ep.filter != nil {
			filterResult, err := ep.filter.Eval(data)
			if err != nil {
				return err
			}
			filterResultBool, err := tuple.ToBool(filterResult)
			if err != nil {
				return err
			}
			// if it evaluated to false, do not further process this tuple
			if !filterResultBool {
				return nil
			}
		}
		// otherwise, compute all the expressions
		result := tuple.Map(make(map[string]tuple.Value, len(ep.projections)))
		for _, proj := range ep.projections {
			value, err := proj.evaluator.Eval(data)
			if err != nil {
				return err
			}
			if err := assignOutputValue(result, proj.alias, value); err != nil {
				return err
			}
		}
		output = append(output, result)
		return nil
	}

	allStreams := make([]string, 0, len(ep.buffers))
	for key := range ep.buffers {
		allStreams = append(allStreams, key)
	}
	if err := procCartProd(allStreams, evalItem); err != nil {
		// NB. ep.prevResults currently points to an slice with
		//     results from the previous run. ep.curResults points
		//     to the same slice. output points to a different slice
		//     with a different underlying array.
		//     in the next run, output will be reusing the underlying
		//     storage of the current ep.prevResults to hold results.
		//     therefore when we leave this function we must make
		//     sure that ep.prevResults and ep.curResults have
		//     different underlying arrays or ISTREAM/DSTREAM will
		//     return wrong results.
		ep.prevResults = output
		return err
	}

	ep.curResults = output
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
func (ep *defaultSelectExecutionPlan) computeResultTuples() ([]tuple.Map, error) {
	// TODO turn this into an iterator/generator pattern
	var output []tuple.Map
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
