package core

import (
	"fmt"
	"pfi/sensorbee/sensorbee/core/tuple"
	"reflect"
	"sync"
)

func newDynamicPipe(inputName string, capacity int) (*dynamicPipeReceiver, *dynamicPipeSender) {
	p := make(chan *tuple.Tuple, capacity) // TODO: the type should be chan []*tuple.Tuple

	r := &dynamicPipeReceiver{
		in: p,
	}

	s := &dynamicPipeSender{
		inputName: inputName,
		out:       p,
	}
	r.sender = s
	return r, s
}

type dynamicPipeReceiver struct {
	in     <-chan *tuple.Tuple
	sender *dynamicPipeSender
}

// close closes the channel from the receiver side. It doesn't directly close
// the channel. Instead, it sends a signal to the sender so that sender can
// close the channel.
func (r *dynamicPipeReceiver) close() {
	// Close the sender because close(r.in) isn't safe. This method uses a
	// goroutine because it might be dead-locked when the channel is full.
	// In that case dynamicPipeSender.Write blocks at s.out <- t having RLock.
	// Then, calling sender.close() is blocked until someone reads a tuple
	// from r.in.
	go func() {
		r.sender.close()
	}()
}

type dynamicPipeSender struct {
	inputName   string
	out         chan<- *tuple.Tuple
	inputClosed <-chan struct{}

	// rwm protects out from write-close conflicts.
	rwm    sync.RWMutex
	closed bool
}

// Write outputs the given tuple to the pipe. This method only returns
// errPipeClosed and never panics.
func (s *dynamicPipeSender) Write(ctx *Context, t *tuple.Tuple) error {
	s.rwm.RLock()
	defer s.rwm.RUnlock()

	if s.closed {
		return errPipeClosed
	}

	t.InputName = s.inputName
	s.out <- t
	return nil
}

// Close closes a channel. When multiple goroutines try to close the channel,
// only one goroutine can actually close it. Other goroutines don't wait until
// the channel is actually closed. Close never fails.
func (s *dynamicPipeSender) Close(ctx *Context) error {
	// Close method is provided to support WriteCloser interface.
	s.close()
	return nil
}

func (s *dynamicPipeSender) close() {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	close(s.out)
}

type dynamicDataSources struct {
	nodeName string

	// m protects state, recvs, and msgChs.
	m     sync.Mutex
	state *topologyStateHolder

	recvs map[string]*dynamicPipeReceiver

	// msgChs is a slice of channels which are connected to goroutines
	// pouring tuples. They receive controlling messages through this channel.
	msgChs []chan<- *dynamicDataSourcesMessage
}

func newDynamicDataSources(nodeName string) *dynamicDataSources {
	s := &dynamicDataSources{
		nodeName: nodeName,
		recvs:    map[string]*dynamicPipeReceiver{},
	}
	s.state = newTopologyStateHolder(&s.m)
	return s
}

type dynamicDataSourcesMessage struct {
	cmd dynamicDataSourcesCommand
	v   interface{}
}

type dynamicDataSourcesCommand int

const (
	ddscAddReceiver dynamicDataSourcesCommand = iota
	ddscStop
	ddscToggleGracefulStop
	ddscStopOnDisconnect
)

func (s *dynamicDataSources) add(name string, r *dynamicPipeReceiver) error {
	// Because dynamicDataSources is used internally and shouldn't return error
	// in most cases, there's no need to check s.recvs with RLock before
	// actually acquiring Lock.
	s.m.Lock()
	defer s.m.Unlock()
	if s.state.getWithoutLock() >= TSStopping {
		return fmt.Errorf("node '%v' already closed its input", s.nodeName)
	}

	// It's safe even if pouringThread panics while executing this method
	// because pour method will drain all channels in s.recvs. So, the sender
	// won't block even if pouringThread doesn't receive the new pipe receiver.
	if _, ok := s.recvs[name]; ok {
		return fmt.Errorf("node '%v' is already receiving tuples from '%v'", s.nodeName, name)
	}
	s.recvs[name] = r
	s.sendMessageWithoutLock(&dynamicDataSourcesMessage{
		cmd: ddscAddReceiver,
		v:   r,
	})
	return nil
}

func (s *dynamicDataSources) sendMessage(msg *dynamicDataSourcesMessage) {
	s.m.Lock()
	defer s.m.Unlock()
	s.sendMessageWithoutLock(msg)
}

func (s *dynamicDataSources) sendMessageWithoutLock(msg *dynamicDataSourcesMessage) {
	for _, ch := range s.msgChs {
		ch <- msg
	}
}

func (s *dynamicDataSources) remove(name string) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.state.getWithoutLock() >= TSStopping {
		return
	}

	// Removing receiver here can result in a dead-lock when pouringThread
	// panics while the channel of receiver is full. In that case, there's no
	// reader who reads a tuple and unblock the sender. To avoid this problem,
	// pouringThread returns all input channels it had at the time the failure
	// occurred and pour method will read tuples from those channels.
	r, ok := s.recvs[name]
	if !ok {
		return
	}
	delete(s.recvs, name)
	r.close() // This eventually closes the channel.
}

// pour pours out tuples for the target Writer. The target must directly be
// connected to a Box or a Sink.
func (s *dynamicDataSources) pour(ctx *Context, w Writer, paralellism int) error {
	if paralellism == 0 {
		paralellism = 1
	}

	var (
		wg            sync.WaitGroup
		logOnce       sync.Once
		collectInputs sync.Once
		inputs        []reflect.SelectCase
		threadErr     error
	)

	err := func() error {
		s.m.Lock()
		defer s.m.Unlock()
		if st, err := s.state.checkAndPrepareForRunningWithoutLock(); err != nil {
			switch st {
			case TSRunning, TSPaused:
				return fmt.Errorf("'%v' already started to receive tuples", s.nodeName)
			case TSStopped:
				return fmt.Errorf("'%v' already stopped receiving tuples", s.nodeName)
			default:
				return fmt.Errorf("'%v' has invalid state: %v", s.nodeName, st)
			}
		}

		genCases := func(msgCh <-chan *dynamicDataSourcesMessage) []reflect.SelectCase {
			cs := make([]reflect.SelectCase, 0, len(s.recvs)+2)
			cs = append(cs, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(msgCh),
			})

			// This case is used as a default case after stopCh is closed.
			// Currently, it only has recv with nil channel so that
			// reflect.Select does nothing on it.
			cs = append(cs, reflect.SelectCase{
				Dir: reflect.SelectRecv,
			})

			for _, r := range s.recvs {
				cs = append(cs, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(r.in),
				})
			}
			return cs
		}

		for i := 0; i < paralellism; i++ {
			msgCh := make(chan *dynamicDataSourcesMessage)
			s.msgChs = append(s.msgChs, msgCh)

			wg.Add(1)
			go func() {
				defer wg.Done()
				ins, err := s.pouringThread(ctx, w, genCases(msgCh))
				collectInputs.Do(func() {
					// It's sufficient to collect input only once. The only
					// problem which might happen is that ins has old receivers.
					// However, they will simply be removed when calling
					// reflect.Select. There might be a case that only one
					// pouringThread has a newly added receiver but it isn't
					// assigned to inputs. To solve that problem, pour method
					// also reads tuples from s.recvs.
					//
					// In addition, when pouringThread panics while remove
					// method is called, s.recvs might not have receivers
					// which pour method should read tuples. However, inputs
					// returned from pouringThread has them. If the returned
					// inputs doesn't have them, that means they were already
					// removed successfully.
					//
					// In conclusion, by combining s.recvs and inputs, all
					// inputs can be drained and no sender will be blocked.
					inputs = ins
				})
				if err != nil {
					logOnce.Do(func() {
						threadErr = err // return only one error
						ctx.Logger.Log(Error, "'%v' stopped with a fatal error: %v", s.nodeName, err)
					})
				}
			}()
		}
		return nil
	}()
	if err != nil {
		return err
	}

	s.state.Set(TSRunning)
	wg.Wait()

	s.m.Lock()
	defer func() {
		s.state.setWithoutLock(TSStopped)
		s.m.Unlock()
	}()

	drainTargets := make([]reflect.SelectCase, 0, len(s.recvs)+len(inputs))
	drainTargets = append(drainTargets, inputs...)
	for _, r := range s.recvs {
		drainTargets = append(drainTargets, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(r.in),
		})
		r.close()
	}
	s.recvs = nil

	go func() {
		// drainTargets might have duplicated channels but it doesn't cause
		// a problem.
		for len(drainTargets) != 0 {
			i, _, ok := reflect.Select(drainTargets)
			if ok {
				continue
			}
			drainTargets[i] = drainTargets[len(drainTargets)-1]
			drainTargets = drainTargets[:len(drainTargets)-1]
		}
	}()

	for _, ch := range s.msgChs {
		close(ch)
	}
	s.msgChs = nil
	return threadErr
}

func (s *dynamicDataSources) pouringThread(ctx *Context, w Writer, cs []reflect.SelectCase) (inputs []reflect.SelectCase, retErr error) {
	const (
		message = iota
		defaultCase

		// maxControlIndex has the max index of special channels used to
		// control this method.
		maxControlIndex = defaultCase
	)

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				if !IsFatalError(err) {
					err = FatalError(err)
				}
				retErr = err
			} else {
				retErr = fmt.Errorf("'%v' got an unknown error through panic: %v", s.nodeName, e)
			}
		}

		if len(cs) > maxControlIndex+1 {
			// Return all inputs this method still has. pour method will read
			// tuples from those input channels so that senders don't block.
			inputs = cs[maxControlIndex+1:]
		}

		// drain channels specific to pouringThread
		go func() {
			cs = cs[:1] // exclude defaultCase
			for len(cs) != 0 {
				i, _, ok := reflect.Select(cs)
				if ok {
					continue
				}
				cs[i] = cs[len(cs)-1]
				cs = cs[:len(cs)-1]
			}
		}()
	}()

	gracefulStopEnabled := false
	stopOnDisconnect := false

receiveLoop:
	for {
		if stopOnDisconnect && len(cs) == maxControlIndex+1 {
			// When stopOnDisconnect is enabled, this loop breaks if the data
			// source doesn't have any input channel. Otherwise, it keeps
			// running because a new input could dynamically be added.
			break
		}

		i, v, ok := reflect.Select(cs) // all cases are receive direction
		if !ok && i != defaultCase {
			if i <= maxControlIndex {
				retErr = FatalError(fmt.Errorf("a controlling channel (%v) of '%v' has been closed", i, s.nodeName))
				return
			}

			// remove the closed channel by swapping it with the last element.
			cs[i], cs[len(cs)-1] = cs[len(cs)-1], cs[i]
			cs = cs[:len(cs)-1]
			continue
		}

		switch i {
		case message:
			msg, ok := v.Interface().(*dynamicDataSourcesMessage)
			if !ok {
				ctx.Logger.Log(Warning, "Received an invalid control message in dynamicDataSources: %v", v.Interface())
				continue
			}

			switch msg.cmd {
			case ddscAddReceiver:
				c, ok := msg.v.(*dynamicPipeReceiver)
				if !ok {
					ctx.Logger.Log(Warning, "Cannot add a new receiver due to a type error")
					break
				}
				cs = append(cs, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(c.in),
				})

			case ddscStop:
				if !gracefulStopEnabled {
					break receiveLoop
				}
				cs[defaultCase].Dir = reflect.SelectDefault // activate the default case

			case ddscToggleGracefulStop:
				gracefulStopEnabled = true

			case ddscStopOnDisconnect:
				stopOnDisconnect = true
			}

		case defaultCase:
			// stop has been called and there's no additional input.
			break receiveLoop

		default:
			t, ok := v.Interface().(*tuple.Tuple)
			if !ok {
				ctx.Logger.Log(Error, "Cannot receive a tuple from a receiver due to a type error")
				break
			}

			err := w.Write(ctx, t)
			switch {
			case IsFatalError(err):
				// logging is done by pour method
				retErr = err
				return

			case IsTemporaryError(err):
				// TODO: retry

			default:
				// Skip this tuple
			}
		}
	}
	return // return values will be set by the defered function.
}

// enableGracefulStop enables graceful stop mode. If the mode is enabled, the
// suorce automatically stops when it doesn't receive any input after stop is
// called.
func (s *dynamicDataSources) enableGracefulStop() {
	// Perhaps this function should be something like 'toggle', but it wasn't
	// necessary at the time of this writing.
	s.sendMessage(&dynamicDataSourcesMessage{
		cmd: ddscToggleGracefulStop,
	})
}

// stopOnDisconnect activates automatic stop when the dynamicDataSources has
// no receiver.
func (s *dynamicDataSources) stopOnDisconnect() {
	s.sendMessage(&dynamicDataSourcesMessage{
		cmd: ddscStopOnDisconnect,
	})
}

// stop stops the source after processing tuples which it currently has.
func (s *dynamicDataSources) stop(ctx *Context) {
	s.m.Lock()
	defer s.m.Unlock()
	if stopped, err := s.state.checkAndPrepareForStoppingWithoutLock(false); stopped || err != nil {
		return
	}

	for _, r := range s.recvs {
		// This eventually closes the channels and remove edges from
		// data sources.
		r.close()
	}
	s.recvs = nil

	s.sendMessageWithoutLock(&dynamicDataSourcesMessage{
		cmd: ddscStop,
	})
	s.state.waitWithoutLock(TSStopped)
}

// dynamicDataDestinations have writers connected to multiple destination nodes and
// distributes tuples to them.
type dynamicDataDestinations struct {
	// nodeName is the name of the node which writes tuples to
	// destinations (i.e. the name of a Source or a Box).
	nodeName string
	rwm      sync.RWMutex
	cond     *sync.Cond
	dsts     map[string]*dynamicPipeSender
	paused   bool
}

func newDynamicDataDestinations(nodeName string) *dynamicDataDestinations {
	d := &dynamicDataDestinations{
		nodeName: nodeName,
		dsts:     map[string]*dynamicPipeSender{},
	}
	d.cond = sync.NewCond(&d.rwm)
	return d
}

func (d *dynamicDataDestinations) add(name string, s *dynamicPipeSender) error {
	d.rwm.Lock()
	defer d.rwm.Unlock()
	if d.dsts == nil {
		return fmt.Errorf("node '%v' already closed its output", d.nodeName)
	}

	if _, ok := d.dsts[name]; ok {
		return fmt.Errorf("node '%v' already has the destination '%v'", d.nodeName, name)
	}
	d.dsts[name] = s
	return nil
}

func (d *dynamicDataDestinations) remove(name string) {
	d.rwm.Lock()
	defer d.rwm.Unlock()
	dst, ok := d.dsts[name]
	if !ok {
		return
	}
	delete(d.dsts, name)
	dst.close()
}

// Write writes tuples to destinations. It doesn't return any error including
// errPipeClosed.
func (d *dynamicDataDestinations) Write(ctx *Context, t *tuple.Tuple) error {
	d.rwm.RLock()
	shouldUnlock := true
	defer func() {
		if shouldUnlock {
			d.rwm.RUnlock()
		}
	}()

	// RLock will be acquired again by the end of the loop.
	for d.paused {
		d.rwm.RUnlock()

		// assuming d.cond.Wait doesn't panic.
		d.rwm.Lock()
		for d.paused {
			d.cond.Wait()
		}
		d.rwm.Unlock()

		// Because paused can be changed in this Unlock -> RLock interval,
		// its value has to be checked again in the next iteration.

		d.rwm.RLock()
	}
	// It's safe even if Close method is called while waiting in the loop above.

	needsCopy := len(d.dsts) > 1
	var closed []string
	for name, dst := range d.dsts {
		// TODO: recovering from panic here instead of using RWLock in
		// dynamicPipeSender might be faster.

		s := t
		if needsCopy {
			s = t.Copy()
		}

		if err := dst.Write(ctx, s); err != nil { // never panics
			// err is always errPipeClosed when it isn't nil.
			// Because the closed destination doesn't do anything harmful,
			// it'll be removed later for performance reason.

			closed = append(closed, name)
		}
	}

	if closed != nil {
		shouldUnlock = false
		d.rwm.RUnlock()
		d.rwm.Lock()
		defer d.rwm.Unlock()
		for _, n := range closed {
			delete(d.dsts, n)
		}
	}
	return nil
}

func (d *dynamicDataDestinations) pause() {
	d.rwm.Lock()
	defer d.rwm.Unlock()
	d.setPaused(true)
}

func (d *dynamicDataDestinations) resume() {
	d.rwm.Lock()
	defer d.rwm.Unlock()
	d.setPaused(false)
}

func (d *dynamicDataDestinations) setPaused(p bool) {
	if d.paused == p {
		return
	}
	d.paused = p
	d.cond.Broadcast()
}

func (d *dynamicDataDestinations) Close(ctx *Context) error {
	d.rwm.Lock()
	defer d.rwm.Unlock()
	for _, dst := range d.dsts {
		dst.close()
	}
	d.dsts = nil
	d.setPaused(false)
	return nil
}
