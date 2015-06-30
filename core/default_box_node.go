package core

type defaultBoxNode struct {
	*defaultNode
	srcs *dataSources
	box  Box
	dsts *dataDestinations

	stopOnOutboundDisconnect bool
}

func (db *defaultBoxNode) Type() NodeType {
	return NTBox
}

func (db *defaultBoxNode) Box() Box {
	return db.box
}

func (db *defaultBoxNode) Input(refname string, config *BoxInputConfig) error {
	s, err := db.topology.dataSource(refname)
	if err != nil {
		return err
	}

	if config == nil {
		config = defaultBoxInputConfig
	}

	if err := checkBoxInputName(db.box, db.name, config.inputName()); err != nil {
		return err
	}

	recv, send := newPipe(config.inputName(), config.capacity())
	if err := s.destinations().add(db.name, send); err != nil {
		return err
	}
	if err := db.srcs.add(s.Name(), recv); err != nil {
		s.destinations().remove(db.name)
		return err
	}
	return nil
}

func (db *defaultBoxNode) run() error {
	if err := db.checkAndPrepareForRunning("box"); err != nil {
		return err
	}

	defer func() {
		defer func() {
			db.dsts.Close(db.topology.ctx)
			db.state.Set(TSStopped)
		}()
		if sb, ok := db.box.(StatefulBox); ok {
			if err := sb.Terminate(db.topology.ctx); err != nil {
				db.topology.ctx.Logger.Log(Error, "Box '%v' in topology '%v' failed to terminate: %v",
					db.name, db.topology.name, err)
			}
		}
	}()
	db.state.Set(TSRunning)
	w := newBoxWriterAdapter(db.box, db.name, db.dsts)
	return db.srcs.pour(db.topology.ctx, w, 1) // TODO: make parallelism configurable
}

func (db *defaultBoxNode) Stop() error {
	db.stop()
	return nil
}

func (db *defaultBoxNode) EnableGracefulStop() {
	db.srcs.enableGracefulStop()
}

func (db *defaultBoxNode) StopOnDisconnect(dir ConnDir) {
	if dir&Inbound != 0 {
		db.srcs.stopOnDisconnect()
	} else if dir&Outbound != 0 {
		db.stateMutex.Lock()
		db.stopOnOutboundDisconnect = true
		db.stateMutex.Unlock()

		if db.dsts.len() == 0 {
			db.stop()
		}
	}
}

func (db *defaultBoxNode) stop() {
	if stopped, err := db.checkAndPrepareForStopping("box"); stopped || err != nil {
		return
	}

	db.state.Set(TSStopping)
	db.srcs.stop(db.topology.ctx) // waits until all tuples get processed.
	db.state.Wait(TSStopped)
}

func (db *defaultBoxNode) destinations() *dataDestinations {
	return db.dsts
}

func (db *defaultBoxNode) dstCallback(e ddEvent) {
	switch e {
	case ddeDisconnect:
		db.stateMutex.Lock()
		shouldStop := db.stopOnOutboundDisconnect
		db.stateMutex.Unlock()

		if shouldStop {
			db.stop()
		}
	}
}
