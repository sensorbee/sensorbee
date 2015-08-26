package udf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/core"
	"sync"
)

// UDSStorage is an interface to support saving and loading UDSs.
type UDSStorage interface {
	// Save returns a writer to write the state data. Save doesn't discard the
	// previously saved data until UDSStorageWriter.Commit is called.
	//
	// Save can be called while a state is being loaded. In such case, behavior
	// is up to each storage. Some possible implementations are: (1) Save
	// creates a new entry for the state and conflicting Loads continue to read
	// the previous data, (2) Save blocks until conflicting Loads finish, and
	// (3) Save returns an error. Implementation (3) isn't recommended because
	// it might result in starvation in highly concurrent workload.
	//
	// Either Commit or Abort of a UDSStorageWriter returned from this method
	// has to be called. When Commit is called, the data is persisted. When
	// Abort is called, the data is discarded and the previous data remains.
	//
	// Save can write header information or other data such as a space for
	// storing checksum later to UDSStorageWriter before returning it. Save can
	// also manipulate the written data as long as the data can be loaded again.
	Save(topology, state string) (UDSStorageWriter, error)

	// Load loads the previously saved data of the state. io.ReadCloser.Close
	// has to be called when it gets unnecessary.
	//
	// Load can be called while a state is being saved. In such case, behavior
	// is up to each storage. A storage's Load can block until Save is done,
	// can return an error, or can even return a reader of the previously saved
	// data.
	//
	// Load returns core.NotExistError when the state doesn't exist.
	Load(topology, state string) (io.ReadCloser, error)

	// List returns a list of names of saved states as a map whose key is a
	// name of a topology. Each value contains names of states in a topology.
	List() (map[string][]string, error)
}

// UDSStorageWriter is used to save a state. An instance of UDSStorageWriter
// doesn't have to be thread-safe. It means that an instance may not be able to
// be used from multiple goroutines. However, different instances can be used
// concurrently so that multiple states can be saved simultaneously.
type UDSStorageWriter interface {
	io.Writer

	// Commit persists the data written to the writer so far and closes it.
	// Write cannot be called once the data is committed.
	Commit() error

	// Abort discard the data written to the writer. Write cannot be called
	// after calling Abort.
	Abort() error
}

type inMemoryUDSStorage struct {
	m          sync.RWMutex
	topologies map[string]*topologyUDSStorage
}

// NewInMemoryUDSStorage creates a new UDSStorage which store all data in
// memory. This storage should only be used for experiment or test purpose.
func NewInMemoryUDSStorage() UDSStorage {
	return &inMemoryUDSStorage{
		topologies: map[string]*topologyUDSStorage{},
	}
}

func (s *inMemoryUDSStorage) Save(topology, state string) (UDSStorageWriter, error) {
	s.m.Lock()
	defer s.m.Unlock()
	t, ok := s.topologies[topology]
	if !ok {
		t = &topologyUDSStorage{
			topologyName: topology,
			states:       map[string][]byte{},
		}
		s.topologies[topology] = t
	}
	return &inMemoryUDSStorageWriter{
		storage:   t,
		buf:       bytes.NewBuffer(nil),
		stateName: state,
	}, nil
}

func (s *inMemoryUDSStorage) Load(topology, state string) (io.ReadCloser, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	t, ok := s.topologies[topology]
	if !ok {
		return nil, core.NotExistError(fmt.Errorf("a topology '%v' was not found", topology))
	}
	st, ok := t.states[state]
	if !ok {
		return nil, core.NotExistError(fmt.Errorf("a UDS '%v' was not found", state))
	}
	return ioutil.NopCloser(bytes.NewReader(st)), nil
}

func (s *inMemoryUDSStorage) List() (map[string][]string, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	res := make(map[string][]string, len(s.topologies))
	for name, t := range s.topologies {
		res[name] = t.list()
	}
	return res, nil
}

type topologyUDSStorage struct {
	m            sync.RWMutex
	topologyName string
	states       map[string][]byte
}

func (t *topologyUDSStorage) list() []string {
	t.m.RLock()
	defer t.m.RUnlock()
	res := make([]string, len(t.states))
	i := 0
	for name := range t.states {
		res[i] = name
		i++
	}
	return res
}

type inMemoryUDSStorageWriter struct {
	storage   *topologyUDSStorage
	buf       *bytes.Buffer
	stateName string
}

func (w *inMemoryUDSStorageWriter) Write(data []byte) (int, error) {
	if w.buf == nil {
		return 0, errors.New("writer is already closed")
	}
	return w.buf.Write(data)
}

func (w *inMemoryUDSStorageWriter) Commit() error {
	if w.buf == nil {
		return errors.New("writer is already closed")
	}

	w.storage.m.Lock()
	defer w.storage.m.Unlock()
	w.storage.states[w.stateName] = w.buf.Bytes()
	w.buf = nil
	return nil
}

func (w *inMemoryUDSStorageWriter) Abort() error {
	if w.buf == nil {
		return errors.New("writer is already closed")
	}
	w.buf = nil
	return nil
}
