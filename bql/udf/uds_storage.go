package udf

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"io"
	"io/ioutil"
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
	//
	// A caller can assign a tag to the saved state so that multiple versions of
	// the UDS can be managed with unique names. When a tag is an empty string,
	// "default" will be used. The valid format of tags is same as node names,
	// which is validated by core.ValidateSymbol.
	Save(topology, state, tag string) (UDSStorageWriter, error)

	// Load loads the previously saved data of the state. io.ReadCloser.Close
	// has to be called when it gets unnecessary.
	//
	// Load can be called while a state is being saved. In such case, behavior
	// is up to each storage. A storage's Load can block until Save is done,
	// can return an error, or can even return a reader of the previously saved
	// data.
	//
	// Load returns core.NotExistError when the state doesn't exist.
	//
	// When a tag is an empty string, "default" will be used.
	Load(topology, state, tag string) (io.ReadCloser, error)

	// ListTopologies returns a list of topologies that have saved states.
	ListTopologies() ([]string, error)

	// List returns a list of names of saved states in a topology as a map
	// whose key is a name of a UDS. Each value contains tags assigned to
	// the state as an array.
	List(topology string) (map[string][]string, error)
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

func (s *inMemoryUDSStorage) Save(topology, state, tag string) (UDSStorageWriter, error) {
	if tag == "" {
		tag = "default"
	} else if err := core.ValidateSymbol(tag); err != nil {
		return nil, fmt.Errorf("tag is ill-formatted: %v", err)
	}

	s.m.Lock()
	defer s.m.Unlock()
	t, ok := s.topologies[topology]
	if !ok {
		t = &topologyUDSStorage{
			topologyName: topology,
			states:       map[string]map[string][]byte{},
		}
		s.topologies[topology] = t
	}
	return &inMemoryUDSStorageWriter{
		storage:   t,
		buf:       bytes.NewBuffer(nil),
		stateName: state,
		tag:       tag,
	}, nil
}

func (s *inMemoryUDSStorage) Load(topology, state, tag string) (io.ReadCloser, error) {
	if tag == "" {
		tag = "default"
	} else if err := core.ValidateSymbol(tag); err != nil {
		return nil, fmt.Errorf("tag is ill-formatted: %v", err)
	}

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
	data, ok := st[tag]
	if !ok {
		return nil, core.NotExistError(fmt.Errorf("a UDS '%v' doesn't have a tag '%v'", state, tag))
	}
	return ioutil.NopCloser(bytes.NewReader(data)), nil
}

func (s *inMemoryUDSStorage) ListTopologies() ([]string, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	res := make([]string, 0, len(s.topologies))
	for name := range s.topologies {
		res = append(res, name)
	}
	return res, nil
}

func (s *inMemoryUDSStorage) List(topology string) (map[string][]string, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	t, ok := s.topologies[topology]
	if !ok {
		return nil, core.NotExistError(fmt.Errorf("a topology '%v' was not found", topology))
	}
	return t.list(), nil
}

type topologyUDSStorage struct {
	m            sync.RWMutex
	topologyName string
	states       map[string]map[string][]byte
}

func (t *topologyUDSStorage) list() map[string][]string {
	t.m.RLock()
	defer t.m.RUnlock()
	res := map[string][]string{}
	for name, tags := range t.states {
		for tag := range tags {
			res[name] = append(res[name], tag)
		}
	}
	return res
}

type inMemoryUDSStorageWriter struct {
	storage   *topologyUDSStorage
	buf       *bytes.Buffer
	stateName string
	tag       string
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
	m := w.storage.states[w.stateName]
	if m == nil {
		m = map[string][]byte{}
	}
	m[w.tag] = w.buf.Bytes()
	w.storage.states[w.stateName] = m
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
