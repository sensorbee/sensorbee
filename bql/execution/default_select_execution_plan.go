package execution

import (
	"fmt"
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/bql/udf"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
)

type defaultSelectExecutionPlan struct {
	// TODO turn this into a list of structs to ensure same length
	// colHeaders stores the names of the result columns.
	colHeaders []string
	// selectors stores the evaluators of the result columns.
	selectors []Evaluator
	// filter stores the evaluator of the filter condition,
	// or nil if there is no WHERE clause.
	filter Evaluator
	// Window information (extracted from LogicalPlan):
	windowSize int64
	windowType parser.RangeUnit
	emitter    parser.Emitter
	// buffer holds data of a single stream window. It will be
	// updated (appended and possibly truncated) whenever
	// Process() is called with a new tuple.
	buffer []*tuple.Tuple
	// curResults holds results of a query over the buffer.
	curResults []tuple.Map
	// prevResults holds results of a query over the buffer
	// in the previous execution run.
	prevResults []tuple.Map
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
	projs := make([]Evaluator, len(lp.Projections))
	colHeaders := make([]string, len(lp.Projections))
	for i, proj := range lp.Projections {
		// compute evaluators for each column
		plan, err := ExpressionToEvaluator(proj, reg)
		if err != nil {
			return nil, err
		}
		projs[i] = plan
		// compute column name
		colHeader := fmt.Sprintf("col_%v", i+1)
		switch projType := proj.(type) {
		case parser.ColumnName:
			colHeader = projType.Name
		case parser.FuncAppAST:
			colHeader = string(projType.Function)
		}
		colHeaders[i] = colHeader
	}
	// compute evaluator for the filter
	var filter Evaluator
	if lp.Filter != nil {
		f, err := ExpressionToEvaluator(lp.Filter, reg)
		if err != nil {
			return nil, err
		}
		filter = f
	}
	// initialize buffer
	var buffer []*tuple.Tuple
	if lp.Unit == parser.Tuples {
		// we already know the required capacity of this buffer
		// if we work with absolute numbers
		buffer = make([]*tuple.Tuple, 0, lp.Value+1)
	}
	return &defaultSelectExecutionPlan{
		colHeaders:  colHeaders,
		selectors:   projs,
		filter:      filter,
		windowSize:  lp.Value,
		windowType:  lp.Unit,
		emitter:     lp.EmitterType,
		buffer:      buffer,
		curResults:  []tuple.Map{},
		prevResults: []tuple.Map{},
	}, nil
}

func (ep *defaultSelectExecutionPlan) Process(input *tuple.Tuple) ([]tuple.Map, error) {
	// stream-to-relation:
	// updates the internal buffer with correct window data
	if err := ep.addTupleToBuffer(input); err != nil {
		return nil, err
	}
	if err := ep.removeOutdatedTuplesFromBuffer(); err != nil {
		return nil, err
	}

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

// addTupleToBuffer appends the received tuple to the internal buffer.
// Note that after calling this function, the buffer may hold more
// items than matches the window specification, so a call to
// removeOutdatedTuplesFromBuffer is necessary afterwards.
func (ep *defaultSelectExecutionPlan) addTupleToBuffer(t *tuple.Tuple) error {
	// TODO maybe a slice is not the best implementation for a queue?
	ep.buffer = append(ep.buffer, t)
	return nil
}

// removeOutdatedTuplesFromBuffer removes tuples from the buffer that
// lie outside the current window as per the statement's window
// specification.
func (ep *defaultSelectExecutionPlan) removeOutdatedTuplesFromBuffer() error {
	curBufSize := int64(len(ep.buffer))
	if ep.windowType == parser.Tuples && curBufSize > ep.windowSize {
		// we just need to take the last `windowSize` items:
		// {a, b, c, d} => {b, c, d}
		ep.buffer = ep.buffer[curBufSize-ep.windowSize : curBufSize]

	} else if ep.windowType == parser.Seconds {
		// we need to remove all items older than `windowSize` seconds,
		// compared to the current tuple
		curTupTime := ep.buffer[curBufSize-1].Timestamp

		// copy "sufficiently new" tuples to new buffer
		newBuf := make([]*tuple.Tuple, 0, curBufSize)
		for _, tup := range ep.buffer {
			dur := curTupTime.Sub(tup.Timestamp)
			if dur.Seconds() <= float64(ep.windowSize) {
				newBuf = append(newBuf, tup)
			}
		}
		ep.buffer = newBuf
	}
	return nil
}

// performQueryOnBuffer executes a SELECT query on the data of the tuples
// currently stored in the buffer. The query results (which is a set of
// tuple.Value, not tuple.Tuple) is stored in ep.curResults. The data
// that was stored in ep.curResults before this method was called is
// moved to ep.prevResults.
//
// Currently performQueryOnBuffer can only perform SELECT ... WHERE ...
// queries on a single relation without aggregate functions, GROUP BY,
// JOIN etc. clauses.
func (ep *defaultSelectExecutionPlan) performQueryOnBuffer() error {
	if len(ep.colHeaders) != len(ep.selectors) {
		return fmt.Errorf("number of columns (%v) doesn't match selectors (%v)",
			len(ep.colHeaders), len(ep.selectors))
	}
	// reuse the allocated memory
	output := ep.prevResults[0:0]
	// remember the previous results
	ep.prevResults = ep.curResults
	for _, t := range ep.buffer {
		// evaluate filter condition and convert to bool
		filterResult, err := ep.filter.Eval(t.Data)
		if err != nil {
			return err
		}
		filterResultBool, err := tuple.ToBool(filterResult)
		if err != nil {
			return err
		}
		// if it evaluated to false, do not further process this tuple
		if !filterResultBool {
			continue
		}
		// otherwise, compute all the expressions
		result := tuple.Map(make(map[string]tuple.Value, len(ep.colHeaders)))
		for idx, selector := range ep.selectors {
			colName := ep.colHeaders[idx]
			value, err := selector.Eval(t.Data)
			if err != nil {
				return err
			}
			result[colName] = value
		}
		output = append(output, result)
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
	if ep.emitter == parser.Rstream {
		// emit all tuples
		output = ep.curResults
	} else if ep.emitter == parser.Istream {
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
	} else if ep.emitter == parser.Dstream {
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
	}
	return output, nil
}
