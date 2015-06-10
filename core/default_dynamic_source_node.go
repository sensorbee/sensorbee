package core

import (
	"pfi/sensorbee/sensorbee/core/tuple"
)

type defaultDynamicSourceNode struct {
	*defaultDynamicNode
	source Source
	dsts   *dynamicDataDestinations
}

func (ds *defaultDynamicSourceNode) run() error {
	if err := ds.checkAndPrepareRunState("source"); err != nil {
		return err
	}

	defer func() {
		defer ds.state.Set(TSStopped)
		ds.dsts.Close(ds.topology.ctx)
	}()
	// TODO: Set TSStarting here and set TSRunning after GenerateStream
	// actually start writing tuples.
	ds.state.Set(TSRunning)
	return ds.source.GenerateStream(ds.topology.ctx, newTraceWriter(ds.dsts, tuple.Output, ds.name))
}

func (ds *defaultDynamicSourceNode) Stop() error {
	if stopped, err := ds.checkAndPrepareStopState("source"); err != nil {
		return err
	} else if stopped {
		return nil
	}

	ds.state.Set(TSStopping)
	if err := ds.source.Stop(ds.topology.ctx); err != nil {
		ds.dsts.Close(ds.topology.ctx) // never fails
		return err
	}
	ds.state.Wait(TSStopped)
	return nil
}

func (ds *defaultDynamicSourceNode) destinations() *dynamicDataDestinations {
	return ds.dsts
}
