package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type defaultDynamicSinkNode struct {
	*defaultDynamicNode
	srcs *dynamicDataSources
	sink Sink
}

func (ds *defaultDynamicSinkNode) Type() NodeType {
	return NTSink
}

func (ds *defaultDynamicSinkNode) Name() string {
	return ds.name
}

func (ds *defaultDynamicSinkNode) State() TopologyStateHolder {
	return ds.state
}

func (ds *defaultDynamicSinkNode) Input(refname string, config *SinkInputConfig) error {
	s, err := ds.topology.dataSource(refname)
	if err != nil {
		return err
	}

	// TODO: schema validation

	if config == nil {
		config = defaultSinkInputConfig
	}

	recv, send := newDynamicPipe("output", config.capacity())
	if err := s.destinations().add(ds.name, send); err != nil {
		return err
	}
	if err := ds.srcs.add(s.Name(), recv); err != nil {
		s.destinations().remove(ds.name)
		return err
	}
	return nil
}

func (ds *defaultDynamicSinkNode) run() error {
	if err := ds.checkAndPrepareRunState("sink"); err != nil {
		return err
	}

	defer ds.state.Set(TSStopped)
	defer func() {
		if err := ds.sink.Close(ds.topology.ctx); err != nil {
			ds.topology.ctx.Logger.Log(Error, "Sink '%v' in topology '%v' failed to stop: %v",
				ds.name, ds.topology.name, err)
		}
	}()
	ds.state.Set(TSRunning)
	return ds.srcs.pour(ds.topology.ctx, newTraceWriter(ds.sink, tuple.Input, ds.name), 1)
}

func (ds *defaultDynamicSinkNode) Stop() error {
	ds.stop()
	return nil
}

func (ds *defaultDynamicSinkNode) stop() {
	if stopped, err := ds.checkAndPrepareStopState("sink"); stopped || err != nil {
		return
	}
	ds.state.Set(TSStopping)
	ds.srcs.stop(ds.topology.ctx)
	ds.state.Wait(TSStopped)
}
