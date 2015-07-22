package core

import (
	"errors"
	"fmt"
	"pfi/sensorbee/sensorbee/data"
)

type defaultSourceNode struct {
	*defaultNode
	config                  *SourceConfig
	source                  Source
	dsts                    *dataDestinations
	pausedOnStartup         bool
	stopOnDisconnectEnabled bool
	runErr                  error
}

func (ds *defaultSourceNode) Type() NodeType {
	return NTSource
}

func (ds *defaultSourceNode) Source() Source {
	return ds.source
}

func (ds *defaultSourceNode) run() (runErr error) {
	if err := ds.checkAndPrepareForRunning("source"); err != nil {
		return err
	}

	defer func() {
		defer ds.state.Set(TSStopped)
		if e := recover(); e != nil {
			// ds.runErr is always nil here
			ds.runErr = fmt.Errorf("the source failed to generate a stream due to panic: %v", e)
		}
		runErr = ds.runErr
		ds.dsts.Close(ds.topology.ctx)
	}()

	ds.runErr = func() error {
		ds.stateMutex.Lock()
		defer ds.stateMutex.Unlock()
		if ds.pausedOnStartup {
			if err := ds.pause(); err != nil {
				return err
			}
			// TSPaused is set by pause method
			return nil
		}
		ds.state.setWithoutLock(TSRunning)
		return nil
	}()
	if ds.runErr != nil {
		return
	}

	ds.runErr = ds.source.GenerateStream(ds.topology.ctx, newTraceWriter(ds.dsts, ETOutput, ds.name))
	return
}

func (ds *defaultSourceNode) Stop() error {
	ds.stateMutex.Lock()
	defer ds.stateMutex.Unlock()
	paused := false
	if ds.state.getWithoutLock() == TSPaused {
		paused = true
	}

	if stopped, err := ds.checkAndPrepareForStoppingWithoutLock("source"); err != nil {
		return err
	} else if stopped {
		return nil
	}

	if paused {
		// The source doesn't have to be resumed since Stop must stop the source
		// without resuming it when it implements Resumable.
		if _, ok := ds.source.(Resumable); !ok {
			// When the default pause&resume implementation is used, dsts just
			// has to be closed to stop correctly. Because dsts.Write doesn't
			// return pipe closed errors, no error log will be written by
			// closing dsts.
			ds.dsts.Close(ds.topology.ctx)
		}
	}

	err := func() (err error) {
		defer func() {
			if e := recover(); e != nil {
				if er, ok := e.(error); ok {
					err = er
				} else {
					err = fmt.Errorf("cannot stop the source due to panic: %v", e)
				}
			}
		}()
		return ds.source.Stop(ds.topology.ctx)
	}()
	if err != nil {
		ds.dsts.Close(ds.topology.ctx) // never fails
		return err
	}
	ds.state.waitWithoutLock(TSStopped)
	return nil
}

func (ds *defaultSourceNode) Pause() error {
	ds.stateMutex.Lock()
	defer ds.stateMutex.Unlock()

	// Because defaultSourceNode will be returned after run method is
	// called by defaultTopology, the possible states are limited.
	switch ds.state.getWithoutLock() {
	case TSRunning:
	case TSPaused:
		return nil
	default:
		return fmt.Errorf("source '%v' is already stopped", ds.name)
	}
	return ds.pause()
}

func (ds *defaultSourceNode) pause() error {
	// pause doesn't acquire lock
	if rn, ok := ds.source.(Resumable); ok {
		// prefer the implementation of the source to the default one.
		err := func() (err error) {
			defer func() {
				if e := recover(); e != nil {
					if er, ok := e.(error); ok {
						err = er
					} else {
						err = fmt.Errorf("the source couldn't stop due to panic: %v", e)
					}
				}
			}()
			return rn.Pause(ds.topology.ctx)
		}()
		if err != nil {
			return err
		}
		ds.state.setWithoutLock(TSPaused)
		return nil
	}

	// If the source doesn't implement Resumable, the default pause/resume
	// implementation in dataDestinations is used.
	ds.dsts.pause()
	ds.state.setWithoutLock(TSPaused)
	return nil
}

func (ds *defaultSourceNode) Resume() error {
	ds.stateMutex.Lock()
	defer ds.stateMutex.Unlock()

	switch ds.state.getWithoutLock() {
	case TSRunning:
		return nil
	case TSPaused:
	default:
		return fmt.Errorf("source '%v' is already stopped", ds.name)
	}

	if rn, ok := ds.source.(Resumable); ok {
		// prefer the implementation of the source to the default one.
		if err := rn.Resume(ds.topology.ctx); err != nil {
			return err
		}
		ds.state.setWithoutLock(TSRunning)
		return nil
	}

	ds.dsts.resume()
	ds.state.setWithoutLock(TSRunning)
	return nil
}

func (ds *defaultSourceNode) Rewind() error {
	rs, ok := ds.source.(RewindableSource)
	if !ok {
		return errors.New("the source doesn't support rewinding")
	}

	ds.stateMutex.Lock()
	defer ds.stateMutex.Unlock()
	if ds.state.getWithoutLock() >= TSStopping {
		return errors.New("the source is stopped")
	}
	return rs.Rewind(ds.topology.ctx)
}

func (ds *defaultSourceNode) Status() data.Map {
	ds.stateMutex.Lock()
	st := ds.state.getWithoutLock()
	stopOnDisconnect := ds.stopOnDisconnectEnabled
	removeOnStop := ds.config.RemoveOnStop
	ds.stateMutex.Unlock()

	m := data.Map{
		"state":        data.String(st.String()),
		"output_stats": ds.dsts.status(),
		"behaviors": data.Map{
			"stop_on_disconnect": data.Bool(stopOnDisconnect),
			"remove_on_stop":     data.Bool(removeOnStop),
		},
	}
	if st == TSStopped && ds.runErr != nil {
		m["error"] = data.String(ds.runErr.Error())
	}
	if s, ok := ds.source.(Statuser); ok {
		m["source"] = s.Status()
	}
	return m
}

func (ds *defaultSourceNode) destinations() *dataDestinations {
	return ds.dsts
}

func (ds *defaultSourceNode) StopOnDisconnect() {
	ds.stateMutex.Lock()
	ds.stopOnDisconnectEnabled = true
	ds.stateMutex.Unlock()

	if ds.dsts.len() == 0 {
		ds.Stop()
	}
}

func (ds *defaultSourceNode) dstCallback(e ddEvent) {
	switch e {
	case ddeDisconnect:
		ds.stateMutex.Lock()
		shouldStop := ds.stopOnDisconnectEnabled
		ds.stateMutex.Unlock()

		if shouldStop {
			ds.Stop()
		}
	}
}

func (ds *defaultSourceNode) RemoveOnStop() {
	ds.stateMutex.Lock()
	ds.config.RemoveOnStop = true
	st := ds.state.getWithoutLock()
	ds.stateMutex.Unlock()

	if st == TSStopped {
		ds.topology.Remove(ds.name)
	}
}
