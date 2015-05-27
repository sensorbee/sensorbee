package bql

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"sync"
)

type bqlBox struct {
	// stmt is the BQL statement executed by this box
	stmt *parser.CreateStreamAsSelectStmt
	// windowSize is the size of the window as given
	// in the RANGE clause -- TODO store an int in the stmt
	windowSize int64
	// mutex protects access to shared state
	mutex sync.Mutex
	// buffer holds data of a stream window
	buffer []*tuple.Tuple
	// curResults holds results of a query over the buffer
	curResults []tuple.Map
	// prevResults holds results of a query over the buffer
	// in the previous execution run
	prevResults []tuple.Map
}

func NewBqlBox(stmt *parser.CreateStreamAsSelectStmt) *bqlBox {
	return &bqlBox{stmt: stmt}
}

func (b *bqlBox) Init(ctx *core.Context) error {
	// take the int from the RANGE clause
	b.windowSize = b.stmt.Value
	// initialize buffer
	if b.stmt.Unit == parser.Tuples {
		// we already know the required capacity of this buffer
		// if we work with absolute numbers
		b.buffer = make([]*tuple.Tuple, 0, b.windowSize+1)
	}
	// initialize result tables
	b.curResults = []tuple.Map{}
	b.prevResults = []tuple.Map{}
	return nil
}

func (b *bqlBox) Process(ctx *core.Context, t *tuple.Tuple, s core.Writer) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// stream-to-relation:
	// updates an internal buffer with correct window data
	b.addTupleToBuffer(t)
	b.removeOutdatedTuplesFromBuffer()

	// relation-to-relation:
	// performs a SELECT query on buffer and writes result
	// to temporary table
	b.performQueryOnBuffer()

	// relation-to-stream:
	// emit new/old/all result tuples and then purge older
	// result buffers
	mkOutTuple := func(data tuple.Map) *tuple.Tuple {
		// capture the current input tuple's metadata
		return &tuple.Tuple{
			Data:          data,
			Timestamp:     t.Timestamp,
			ProcTimestamp: t.ProcTimestamp,
			BatchID:       t.BatchID,
		}
	}
	b.emitResultTuples(mkOutTuple, ctx, s)
	b.removeOutdatedResultBuffers()

	return nil
}

func (b *bqlBox) addTupleToBuffer(t *tuple.Tuple) {
	b.buffer = append(b.buffer, t)
}

func (b *bqlBox) removeOutdatedTuplesFromBuffer() {
	curBufSize := int64(len(b.buffer))
	if b.stmt.Unit == parser.Tuples && curBufSize > b.windowSize {
		// we just need to take the last `windowSize` items:
		// {a, b, c, d} => {b, c, d}
		b.buffer = b.buffer[curBufSize-b.windowSize : curBufSize]

	} else if b.stmt.Unit == parser.Seconds {
		// we need to remove all items older than `windowSize` seconds,
		// compared to the current tuple
		curTupTime := b.buffer[curBufSize-1].Timestamp

		// copy "sufficiently new" tuples to new buffer
		newBuf := make([]*tuple.Tuple, 0, curBufSize)
		for _, tup := range b.buffer {
			dur := curTupTime.Sub(tup.Timestamp)
			if dur.Seconds() <= float64(b.windowSize) {
				newBuf = append(newBuf, tup)
			}
		}
		b.buffer = newBuf
	}
}

func (b *bqlBox) performQueryOnBuffer() {
	// TODO don't use a dummy filter, but the actual statement
	results := []tuple.Map{}
	for _, tup := range b.buffer {
		value, _ := tup.Data.Get("int")
		intVal, _ := value.AsInt()
		// filter
		if intVal%2 == 0 {
			res := tuple.Map{}
			// projection
			res["int"] = value
			results = append(results, res)
		}
	}
	b.curResults = results
}

func (b *bqlBox) emitResultTuples(mkTuple func(data tuple.Map) *tuple.Tuple,
	ctx *core.Context, s core.Writer) {

	if b.stmt.EmitterType == parser.Rstream {
		// emit all tuples
		for _, res := range b.curResults {
			outTup := mkTuple(res)
			// TODO error handling
			s.Write(ctx, outTup)
		}
	} else if b.stmt.EmitterType == parser.Istream {
		// emit only new tuples
		for _, res := range b.curResults {
			// check if this tuple is already present in the previous results
			found := false
			for _, prevRes := range b.prevResults {
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
			// if we arrive here, `res` is not contained in prevRes
			outTup := mkTuple(res)
			// TODO error handling
			s.Write(ctx, outTup)
		}
	} else if b.stmt.EmitterType == parser.Dstream {
		// emit only old tuples
		for _, prevRes := range b.prevResults {
			// check if this tuple is present in the current results
			found := false
			for _, res := range b.curResults {
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
			// if we arrive here, `res` is not contained in prevRes
			outTup := mkTuple(prevRes)
			// TODO error handling
			s.Write(ctx, outTup)
		}
	}
}

func (b *bqlBox) removeOutdatedResultBuffers() {
	b.prevResults = b.curResults
	b.curResults = nil
}

func (b *bqlBox) Terminate(ctx *core.Context) error {
	// TODO cleanup
	return nil
}
