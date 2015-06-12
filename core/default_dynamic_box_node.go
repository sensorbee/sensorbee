package core

type defaultDynamicBoxNode struct {
	*defaultDynamicNode
	srcs *dynamicDataSources
	box  Box
	dsts *dynamicDataDestinations
}

func (db *defaultDynamicBoxNode) Type() NodeType {
	return NTBox
}

func (db *defaultDynamicBoxNode) Input(refname string, config *BoxInputConfig) error {
	s, err := db.topology.dataSource(refname)
	if err != nil {
		return err
	}

	// TODO: schema validation

	if config == nil {
		config = defaultBoxInputConfig
	}

	if err := checkBoxInputName(db.box, db.name, config.inputName()); err != nil {
		return err
	}

	recv, send := newDynamicPipe(config.inputName(), config.capacity())
	if err := s.destinations().add(db.name, send); err != nil {
		return err
	}
	if err := db.srcs.add(s.Name(), recv); err != nil {
		s.destinations().remove(db.name)
		return err
	}
	return nil
}

func (db *defaultDynamicBoxNode) run() error {
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

func (db *defaultDynamicBoxNode) Stop() error {
	db.stop()
	return nil
}

func (db *defaultDynamicBoxNode) EnableGracefulStop() {
	db.srcs.enableGracefulStop()
}

func (db *defaultDynamicBoxNode) StopOnDisconnect() {
	db.srcs.stopOnDisconnect()
}

func (db *defaultDynamicBoxNode) stop() {
	if stopped, err := db.checkAndPrepareForStopping("box"); stopped || err != nil {
		return
	}

	db.state.Set(TSStopping)
	db.srcs.stop(db.topology.ctx) // waits until all tuples get processed.
	db.state.Wait(TSStopped)
}

func (db *defaultDynamicBoxNode) destinations() *dynamicDataDestinations {
	return db.dsts
}
