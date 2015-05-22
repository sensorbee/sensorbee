package bql

import (
	"pfi/sensorbee/sensorbee/bql/parser"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/core/tuple"
)

type bqlBox struct {
	stmt *parser.CreateStreamStmt
}

func NewBqlBox(stmt *parser.CreateStreamStmt) *bqlBox {
	return &bqlBox{stmt}
}

func (b *bqlBox) Init(ctx *core.Context) error {
	// TODO initialize box
	return nil
}

func (b *bqlBox) Process(ctx *core.Context, t *tuple.Tuple, s core.Writer) error {
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
	b.emitResultTuples(s)
	b.removeOutdatedResultBuffers()

	return nil
}

func (b *bqlBox) addTupleToBuffer(t *tuple.Tuple) {
	// TODO
}

func (b *bqlBox) removeOutdatedTuplesFromBuffer() {
	// TODO
}

func (b *bqlBox) performQueryOnBuffer() {
	// TODO
}

func (b *bqlBox) emitResultTuples(s core.Writer) {
	// TODO
}

func (b *bqlBox) removeOutdatedResultBuffers() {
	// TODO
}

func (b *bqlBox) Terminate(ctx *core.Context) error {
	// TODO cleanup
	return nil
}

func (b *bqlBox) InputConstraints() (*core.BoxInputConstraints, error) {
	return nil, nil
}

func (b *bqlBox) OutputSchema(s map[string]*core.Schema) (*core.Schema, error) {
	return nil, nil
}
