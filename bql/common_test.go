package bql

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
)

func newTestTopology() core.Topology {
	t, _ := core.NewDefaultTopology(core.NewContext(nil), "testTopology")
	return t
}

func addBQLToTopology(tb *TopologyBuilder, bql string) error {
	p := parser.New()
	// execute all parsed statements
	stmts, err := p.ParseStmts(bql)
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		_, err := tb.AddStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

type dummyUDS struct {
	num int64
}

func newDummyUDS(ctx *core.Context, params data.Map) (core.SharedState, error) {
	s := &dummyUDS{}
	if err := s.setNum(params); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *dummyUDS) setNum(params data.Map) error {
	v, ok := params["num"]
	if !ok {
		return nil
	}

	n, err := data.ToInt(v)
	if err != nil {
		return err
	}
	s.num = n
	return nil
}

func (s *dummyUDS) Terminate(ctx *core.Context) error {
	return nil
}

type dummyUpdatableUDS struct {
	dummyUDS
}

var (
	_ core.Updater = &dummyUpdatableUDS{}
)

func (s *dummyUpdatableUDS) Update(ctx *core.Context, params data.Map) error {
	return s.setNum(params)
}

func (s *dummyUpdatableUDS) Save(ctx *core.Context, w io.Writer, params data.Map) error {
	return binary.Write(w, binary.LittleEndian, s.num)
}

type dummyUpdatableUDSCreator struct {
}

func (*dummyUpdatableUDSCreator) CreateState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	state, _ := newDummyUDS(ctx, params)
	uds, _ := state.(*dummyUDS)
	s := &dummyUpdatableUDS{
		dummyUDS: *uds,
	}
	return s, nil
}

func (*dummyUpdatableUDSCreator) LoadState(ctx *core.Context, r io.Reader, params data.Map) (core.SharedState, error) {
	s := &dummyUpdatableUDS{}
	if err := binary.Read(r, binary.LittleEndian, &s.num); err != nil {
		return nil, err
	}
	return s, nil
}

type dummySelfLoadableUDS struct {
	dummyUpdatableUDS
}

func (s *dummySelfLoadableUDS) Load(ctx *core.Context, r io.Reader, params data.Map) error {
	return binary.Read(r, binary.LittleEndian, &s.num)
}

type dummySelfLoadableUDSCreator struct {
}

func (*dummySelfLoadableUDSCreator) CreateState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	state, _ := newDummyUDS(ctx, params)
	uds, _ := state.(*dummyUDS)
	s := &dummySelfLoadableUDS{
		dummyUpdatableUDS: dummyUpdatableUDS{
			dummyUDS: *uds,
		},
	}
	return s, nil
}

func (*dummySelfLoadableUDSCreator) LoadState(ctx *core.Context, r io.Reader, params data.Map) (core.SharedState, error) {
	s := &dummySelfLoadableUDS{}
	if err := s.Load(ctx, r, params); err != nil {
		return nil, err
	}
	return s, nil
}

func init() {
	udf.MustRegisterGlobalUDSCreator("dummy_uds", udf.UDSCreatorFunc(newDummyUDS))
	udf.MustRegisterGlobalUDSCreator("dummy_updatable_uds", &dummyUpdatableUDSCreator{})
	udf.MustRegisterGlobalUDSCreator("dummy_self_loadable_uds", &dummySelfLoadableUDSCreator{})
}

type duplicateUDSF struct {
	dup int
}

func (d *duplicateUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	for i := 0; i < d.dup; i++ {
		w.Write(ctx, t.Copy())
	}
	return nil
}

func (d *duplicateUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createDuplicateUDSF(decl udf.UDSFDeclarer, stream string, dup int) (udf.UDSF, error) {
	if err := decl.Input(stream, &udf.UDSFInputConfig{
		InputName: "test",
		Capacity:  999,
		DropMode:  core.DropLatest,
	}); err != nil {
		return nil, err
	}

	return &duplicateUDSF{
		dup: dup,
	}, nil
}

func failingUDSFCreator(decl udf.UDSFDeclarer, stream string, dup int) (udf.UDSF, error) {
	return nil, errors.New("test UDSF creation failed")
}

func init() {
	udf.MustRegisterGlobalUDSFCreator("duplicate", udf.MustConvertToUDSFCreator(createDuplicateUDSF))
	udf.MustRegisterGlobalUDSFCreator("failing_duplicate", udf.MustConvertToUDSFCreator(failingUDSFCreator))
}

type sequenceUDSF struct {
	cnt int
}

// TODO: remove this WaitGroup after supporting pause/resume of streams
var (
	wgSequenceUDSF sync.WaitGroup
)

func (s *sequenceUDSF) Process(ctx *core.Context, t *core.Tuple, w core.Writer) error {
	wgSequenceUDSF.Wait()
	for i := 1; i <= s.cnt; i++ {
		if err := w.Write(ctx, core.NewTuple(data.Map{"int": data.Int(i)})); err != nil {
			if err == core.ErrSourceStopped {
				return err
			}
		}
	}
	return nil
}

func (s *sequenceUDSF) Terminate(ctx *core.Context) error {
	return nil
}

func createSequenceUDSF(decl udf.UDSFDeclarer, cnt int) (udf.UDSF, error) {
	return &sequenceUDSF{
		cnt: cnt,
	}, nil
}

func init() {
	udf.MustRegisterGlobalUDSFCreator("test_sequence", udf.MustConvertToUDSFCreator(createSequenceUDSF))
}
