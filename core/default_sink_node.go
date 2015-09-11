package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/data"
)

type defaultSinkNode struct {
	*defaultNode
	config *SinkConfig
	srcs   *dataSources
	sink   Sink

	gracefulStopEnabled     bool
	stopOnDisconnectEnabled bool
	runErr                  error
}

func (ds *defaultSinkNode) Type() NodeType {
	return NTSink
}

func (ds *defaultSinkNode) Sink() Sink {
	return ds.sink
}

func (ds *defaultSinkNode) Name() string {
	return ds.name
}

func (ds *defaultSinkNode) State() TopologyStateHolder {
	return ds.state
}

func (ds *defaultSinkNode) Input(refname string, config *SinkInputConfig) error {
	s, err := ds.topology.dataSource(refname)
	if err != nil {
		return err
	}

	if config == nil {
		config = defaultSinkInputConfig
	}

	recv, send := newPipe("output", config.capacity())
	send.dropMode = config.DropMode
	if err := s.destinations().add(ds.name, send); err != nil {
		return err
	}
	if err := ds.srcs.add(s.Name(), recv); err != nil {
		s.destinations().remove(ds.name)
		return err
	}

	if !s.destinations().isDroppedTupleReportingEnabled() {
		// disable dropped tuple reporting to avoid infinite reporting loop
		ds.srcs.disableDroppedTupleReporting()
	}
	return nil
}

func (ds *defaultSinkNode) run() (runErr error) {
	if err := ds.checkAndPrepareForRunning("sink"); err != nil {
		return err
	}

	defer ds.state.Set(TSStopped)
	defer func() {
		defer func() {
			if e := recover(); e != nil {
				if ds.runErr == nil {
					ds.runErr = fmt.Errorf("the box couldn't be terminated due to panic: %v", e)
				} else {
					ds.topology.ctx.ErrLog(fmt.Errorf("%v", e)).WithFields(nodeLogFields(NTBox, ds.name)).
						Error("Cannot terminate the box due to panic")
				}
			}
			runErr = ds.runErr
		}()
		if err := ds.sink.Close(ds.topology.ctx); err != nil {
			ds.runErr = err
			ds.topology.ctx.ErrLog(err).WithFields(nodeLogFields(NTSink, ds.name)).
				Error("Cannot stop the sink")
		}
	}()
	ds.state.Set(TSRunning)
	ds.runErr = ds.srcs.pour(ds.topology.ctx, newTraceWriter(ds.sink, ETInput, ds.name), 1)
	return
}

func (ds *defaultSinkNode) Stop() error {
	ds.stop()
	return nil
}

func (ds *defaultSinkNode) EnableGracefulStop() {
	ds.stateMutex.Lock()
	ds.gracefulStopEnabled = true
	ds.stateMutex.Unlock()
	ds.srcs.enableGracefulStop()
}

func (ds *defaultSinkNode) StopOnDisconnect() {
	ds.stateMutex.Lock()
	ds.stopOnDisconnectEnabled = true
	ds.stateMutex.Unlock()
	ds.srcs.stopOnDisconnect()
}

func (ds *defaultSinkNode) stop() {
	if stopped, err := ds.checkAndPrepareForStopping("sink"); stopped || err != nil {
		return
	}
	ds.srcs.stop(ds.topology.ctx)
	ds.state.Wait(TSStopped)
}

func (ds *defaultSinkNode) Status() data.Map {
	ds.stateMutex.Lock()
	st := ds.state.getWithoutLock()
	gstop := ds.gracefulStopEnabled
	stopOnDisconnect := ds.stopOnDisconnectEnabled
	removeOnStop := ds.config.RemoveOnStop
	ds.stateMutex.Unlock()

	m := data.Map{
		"state":       data.String(st.String()),
		"input_stats": ds.srcs.status(),
		"behaviors": data.Map{
			"stop_on_disconnect": data.Bool(stopOnDisconnect),
			"graceful_stop":      data.Bool(gstop),
			"remove_on_stop":     data.Bool(removeOnStop),
		},
	}
	if st == TSStopped && ds.runErr != nil {
		m["error"] = data.String(ds.runErr.Error())
	}
	if s, ok := ds.sink.(Statuser); ok {
		m["sink"] = s.Status()
	}
	return m
}

func (ds *defaultSinkNode) RemoveOnStop() {
	ds.stateMutex.Lock()
	ds.config.RemoveOnStop = true
	st := ds.state.getWithoutLock()
	ds.stateMutex.Unlock()

	if st == TSStopped {
		ds.topology.Remove(ds.name)
	}
}
