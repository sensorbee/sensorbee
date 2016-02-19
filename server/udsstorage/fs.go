package udsstorage

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// fsUDSStorage is a UDSStorage which store states as files on a filesystem.
// This simple storage doesn't provide any version controlling capability.
// It doesn't provide checksum, either. If such capability is required, another
// UDSStorage should be implemented.
type fsUDSStorage struct {
	dirPath     string
	tempDirPath string
}

var (
	_ udf.UDSStorage = &fsUDSStorage{}
)

func NewFS(dir, tempDir string) (udf.UDSStorage, error) {
	if err := validateDir(dir); err != nil {
		return nil, fmt.Errorf("dir (%v) isn't valid: %v", dir, err)
	}
	if tempDir == "" {
		tempDir = dir
	} else if err := validateDir(tempDir); err != nil {
		return nil, fmt.Errorf("temp_dir (%v) isn't valid: %v", tempDir, err)
	}

	return &fsUDSStorage{
		dirPath:     dir,
		tempDirPath: tempDir,
	}, nil
}

func validateDir(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("it isn't a directory")
	}
	return nil
}

func (s *fsUDSStorage) Save(topology, state, tag string) (udf.UDSStorageWriter, error) {
	if tag == "" {
		tag = "default"
	}
	if err := core.ValidateNodeName(tag); err != nil {
		return nil, err
	}

	f, err := s.stateTempFile(topology, state, tag)
	if err != nil {
		return nil, err
	}
	return &onDiskUDSStorageWriter{
		s:          s,
		f:          f,
		w:          bufio.NewWriter(f),
		targetPath: s.stateFilepath(topology, state, tag),
	}, nil
}

func (s *fsUDSStorage) stateTempFile(topology, state, tag string) (*os.File, error) {
	return ioutil.TempFile(s.tempDirPath, s.stateFilename(topology, state, tag))
}

func (s *fsUDSStorage) Load(topology, state, tag string) (io.ReadCloser, error) {
	if tag == "" {
		tag = "default"
	}
	if err := core.ValidateNodeName(tag); err != nil {
		return nil, err
	}

	f, err := os.Open(s.stateFilepath(topology, state, tag))
	if err != nil {
		return nil, err
	}
	return f, nil
}

var (
	fsUDSStorageFilePathRegexp = regexp.MustCompile(`^(.+)-(.+)-(.+).state$`)
)

func (s *fsUDSStorage) ListTopologies() ([]string, error) {
	fs, err := ioutil.ReadDir(s.dirPath)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, f := range fs {
		m := fsUDSStorageFilePathRegexp.FindStringSubmatch(f.Name())
		if m == nil {
			continue
		}
		res = append(res, m[1])
	}
	return res, nil
}

func (s *fsUDSStorage) List(topology string) (map[string][]string, error) {
	fs, err := ioutil.ReadDir(s.dirPath)
	if err != nil {
		return nil, err
	}

	res := map[string][]string{}
	for _, f := range fs {
		m := fsUDSStorageFilePathRegexp.FindStringSubmatch(f.Name())
		if m == nil {
			continue
		}
		if m[1] != topology {
			continue
		}
		res[m[2]] = append(res[m[2]], m[3])
	}
	if len(res) == 0 {
		return nil, core.NotExistError(fmt.Errorf("a topology '%v' was not found", topology))
	}
	return res, nil
}

func (s *fsUDSStorage) stateFilename(topology, state, tag string) string {
	return fmt.Sprintf("%v-%v-%v.state", topology, state, tag)
}

func (s *fsUDSStorage) stateFilepath(topology, state, tag string) string {
	return filepath.Join(s.dirPath, s.stateFilename(topology, state, tag))
}

type onDiskUDSStorageWriter struct {
	s *fsUDSStorage
	f *os.File
	w *bufio.Writer

	targetPath string
}

func (w *onDiskUDSStorageWriter) Write(data []byte) (int, error) {
	if w.w == nil {
		return 0, errors.New("writer is already closed")
	}
	return w.w.Write(data)
}

func (w *onDiskUDSStorageWriter) Commit() (err error) {
	if w.w == nil {
		return errors.New("writer is already closed")
	}
	defer func() {
		if w.f != nil {
			if e := w.f.Close(); e != nil && err == nil {
				err = e
			}
			w.f = nil
		}
	}()

	if e := w.w.Flush(); e != nil {
		return e
	}
	w.w = nil

	f := w.f
	// TODO: Name doesn't return an absolute path, so this code might fail with
	// some scenarios.
	fn := f.Name()
	w.f = nil
	if e := f.Close(); e != nil {
		return e
	}

	// This Rename might not always be atomic and application level locking
	// might be required.
	if e := os.Rename(fn, w.targetPath); e != nil {
		return e
	}
	return nil
}

func (w *onDiskUDSStorageWriter) Abort() error {
	if w.w == nil {
		return errors.New("writer is already closed")
	}
	w.w = nil

	f := w.f
	w.f = nil
	fn := f.Name()
	if err := f.Close(); err != nil {
		return err
	}
	return os.Remove(fn)
}
