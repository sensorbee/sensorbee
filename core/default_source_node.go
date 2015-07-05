package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/data"
)

type defaultSourceNode struct {
	*defaultNode
	source          Source
	dsts            *dataDestinations
	pausedOnStartup bool
	runErr          error
}

func (ds *defaultSourceNode) Type() NodeType {
	return NTSource
}

func (ds *defaultSourceNode) Source() Source {
	return ds.source
}

func (ds *defaultSourceNode) run() error {
	if err := ds.checkAndPrepareForRunning("source"); err != nil {
		return err
	}

	defer func() {
		defer ds.state.Set(TSStopped)
		ds.dsts.Close(ds.topology.ctx)
	}()

	err := func() error {
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
	if err != nil {
		return err
	}

	ds.runErr = ds.source.GenerateStream(ds.topology.ctx, newTraceWriter(ds.dsts, ETOutput, ds.name))
	return ds.runErr
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
		// without resuming it when it implements ResumableNode.
		if _, ok := ds.source.(ResumableNode); !ok {
			// When the default pause&resume implementation is used, dsts just
			// has to be closed to stop correctly. Because dsts.Write doesn't
			// return pipe closed errors, no error log will be written by
			// closing dsts.
			ds.dsts.Close(ds.topology.ctx)
		}
	}

	if err := ds.source.Stop(ds.topology.ctx); err != nil {
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
	if rn, ok := ds.source.(ResumableNode); ok {
		// prefer the implementation of the source to the default one.
		if err := rn.Pause(ds.topology.ctx); err != nil {
			return err
		}
		ds.state.setWithoutLock(TSPaused)
		return nil
	}

	// If the source doesn't implement ResumableNode, the default pause/resume
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

	if rn, ok := ds.source.(ResumableNode); ok {
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

func (ds *defaultSourceNode) Status() data.Map {
	st := ds.state.Get()
	var errMsg string
	if st == TSStopped {
		errMsg = ds.runErr.Error()
	}
	return data.Map{
		"state":        data.String(st.String()),
		"error":        data.String(errMsg),
		"output_stats": ds.dsts.status(),
	}
}

func (ds *defaultSourceNode) destinations() *dataDestinations {
	return ds.dsts
}
