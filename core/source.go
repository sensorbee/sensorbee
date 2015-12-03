package core

import (
	"errors"
	"pfi/sensorbee/sensorbee/data"
	"sync"
)

// A Source describes an entity that inserts data into a topology
// from the outside world (e.g., a fluentd instance).
type Source interface {
	// GenerateStream will start creating tuples and writing them to
	// the given Writer in a blocking way. (Start as a goroutine to start
	// the process in the background.) It will return when all tuples
	// have been written (in the case of a finite data source) or if
	// there was a severe error. The context that is passed in will be
	// used as a parameter to the Write method of the given Writer.
	GenerateStream(ctx *Context, w Writer) error

	// Stop will tell the Source to stop emitting tuples. After this
	// function returns, no more calls to Write shall be made on the
	// Writer passed in to GenerateStream.
	//
	// Stop could be called after GenerateStream returns. However,
	// it's guaranteed that Stop won't be called more than once by
	// components in SensorBee's core package.
	//
	// Stop won't be called if GenerateStream wasn't called.
	Stop(ctx *Context) error
}

// RewindableSource is a Source which can be rewound and generate the same
// stream from the beginning again (e.g. file based source).
//
// Until Stop is called, RewindableSource must not return from GenerateStream
// after it has generated all tuples.
//
// It must be resumable because it's impossible to provide rewinding without
// handling pause and resume appropriately.
type RewindableSource interface {
	Source
	Resumable

	// Rewind rewinds the stream the Source has. Rewind may be called
	// concurrently while GenerateStream is running. If Rewind is called while
	// the Source is being paused, the Source must not be resumed until
	// Resume is explicitly called. The Source doesn't have to care about
	// Pause/Resume if it doesn't implement Resumable.
	Rewind(ctx *Context) error
}

type rewindableSource struct {
	rwm              sync.RWMutex
	state            *topologyStateHolder
	rewind           bool
	waitingForRewind bool

	// rewindEnabled indicates if this rewindableSource actually supports
	// Rewind method. When this is false, Rewind method always returns an error.
	rewindEnabled bool

	forceStop chan struct{}
	source    Source

	// TODO: add methods to satisfy other important interfaces
}

var (
	_ RewindableSource = &rewindableSource{}
	_ Statuser         = &rewindableSource{}
)

var (
	// ErrSourceRewound is returned when the source is rewound and it has to
	// reproduce the stream again.
	ErrSourceRewound = errors.New("the source has been rewound")

	// ErrSourceStopped is returned when the source tries to write a tuple
	// after it's stopped. It's currently only returned from a writer passed
	// to a source created by NewRewindableSource.
	ErrSourceStopped = errors.New("the source has been stopped")
)

// NewRewindableSource creates a rewindable source from a non-rewindable source.
// The source passed to this function must satisfy the following requirements:
//
//	1. Its GenerateStream can safely be called multiple times.
//	2. Its GenerateStream must return when ErrSourceRewound or ErrSourceStopped
//	   is returned from the Writer. It must return the same err instance
//	   returned from the writer.
//
// It can be resumable, but its Pause and Resume won't be called. It doesn't
// have to implement Stop method (i.e. it can just return nil), either, although
// it has to provide it. Instead of implementing Stop method, it can return
// from GenerateStream when the writer returned ErrSourceStopped. If the Source
// has to clean up resources, it can implement Stop to do it. However,
// GenerateStream is still running while Stop is being called. Therefore, all
// resource allocation and deallocation should be done in GenerateStream rather
// than in an initialization function and Stop.
//
// The interface returned from this function will support following interfaces
// if the given source implements them:
//
//	* Statuser
//
// Known issue: There's one problem with NewRewindableSource. Stop method could
// block when the original source's GenerateStream doesn't generate any tuple
// (i.e. doesn't write any tuple) without returning from GenerateStream since
// whether the source is stopped is only determined by the error returned from
// Write.
func NewRewindableSource(s Source) RewindableSource {
	return newRewindableSource(s, true)
}

func newRewindableSource(s Source, rewindEnabled bool) RewindableSource {
	r := &rewindableSource{
		rewindEnabled: rewindEnabled,
		forceStop:     make(chan struct{}, 1),
		source:        s,
	}
	r.state = newTopologyStateHolder(&r.rwm)
	return r
}

func (r *rewindableSource) GenerateStream(ctx *Context, w Writer) error {
	defer r.state.Set(TSStopped)

	// Create a wrapper writer to handle pause/resume and rewind.
	rewindWriter := WriterFunc(func(ctx *Context, t *Tuple) error {
		err := func() error {
			r.rwm.RLock()
			defer r.rwm.RUnlock()
		resumeLoop:
			for !r.rewind { // wait for resume
				switch r.state.getWithoutLock() {
				case TSStopping, TSStopped:
					return ErrSourceStopped

				case TSPaused:
					r.rwm.RUnlock()
					r.rwm.Lock()
					// Don't use state.waitWithoutLock. Wait manually because
					// this loop has to check r.rewind, too.
					r.state.cond.Wait()
					r.rwm.Unlock()
					r.rwm.RLock() // for defer above

				default:
					break resumeLoop
				}
			}
			if r.rewind {
				// The source must return after receiving this error.
				return ErrSourceRewound
			}
			return nil
		}()
		if err != nil {
			return err
		}
		return w.Write(ctx, t) // pass the tuple to the original writer
	})

	ch := make(chan error, 1)
	go func() {
		for {
			if err := r.source.GenerateStream(ctx, rewindWriter); err != nil {
				if err != ErrSourceRewound {
					r.rwm.Lock()
					r.rewind = false
					r.state.cond.Broadcast()
					r.rwm.Unlock()
					if err == ErrSourceStopped {
						ch <- nil
						return
					}
					ch <- err
					return
				}
			}

			// wait for rewind or stop
			shouldReturn := func() bool {
				r.rwm.Lock()
				defer r.rwm.Unlock()
				if !r.rewindEnabled {
					return true
				}

				defer func() {
					r.rewind = false
					r.waitingForRewind = false
					r.state.cond.Broadcast()
				}()
				r.waitingForRewind = true

				for {
					if r.state.getWithoutLock() >= TSStopping {
						return true
					}
					if r.rewind {
						return false
					}
					r.state.cond.Wait() // wait until the status is updated.
				}
			}()
			if shouldReturn {
				ch <- nil
				return
			}

			// rewindableSource must not stop (i.e. return) until Stop is called.
		}
	}()

	select {
	case err := <-ch:
		return err
	case <-r.forceStop:
		go func() {
			// Wait until the source
			if err := <-ch; err == nil {
				ctx.ErrLog(err).Info("The source which has forcibly been stopped finished generating stream")
			} else {
				ctx.ErrLog(err).Warn("The source which has forcibly been stopped returned an error")
			}
		}()
		return errors.New("the source has been stopped forcibly")
	}
}

func (r *rewindableSource) Stop(ctx *Context) error {
	if r.state.Get() >= TSStopping { // just in case
		r.state.Wait(TSStopped)
		return nil
	}

	r.state.Set(TSStopping)
	defer r.state.Wait(TSStopped)
	if err := r.source.Stop(ctx); err != nil {
		r.forceStop <- struct{}{}
		return err
	}
	return nil
}

func (r *rewindableSource) Pause(ctx *Context) error {
	return r.state.Set(TSPaused)
}

func (r *rewindableSource) Resume(ctx *Context) error {
	return r.state.Set(TSRunning)
}

func (r *rewindableSource) Rewind(ctx *Context) error {
	r.rwm.Lock()
	defer r.rwm.Unlock()
	if !r.rewindEnabled {
		return errors.New("this source doesn't support rewind")
	}

	r.rewind = true
	r.state.cond.Broadcast()

	for r.rewind {
		r.state.cond.Wait()
	}
	return nil
}

func (r *rewindableSource) Status() data.Map {
	r.rwm.RLock()
	waiting := r.waitingForRewind
	enabled := r.rewindEnabled
	r.rwm.RUnlock()
	m := data.Map{
		"rewindable":         data.Bool(enabled),
		"waiting_for_rewind": data.Bool(waiting),
	}
	if s, ok := r.source.(Statuser); ok {
		m["internal_source"] = s.Status()
	}
	return m
}

// ImplementSourceStop implements Stop method of a Source in a thread-safe
// manner on behalf of the given Source. Source passed to this function must
// follow the rule described in NewRewindableSource with one exception that
// the Writer doesn't return ErrSourceRewound.
func ImplementSourceStop(s Source) Source {
	// This is implemented as a rewindableSource with rewind disabled.
	return newRewindableSource(s, false)
}
