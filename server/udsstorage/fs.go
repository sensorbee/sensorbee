package udsstorage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pfi/sensorbee/sensorbee/bql/udf"
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

func (s *fsUDSStorage) Save(topology, state string) (udf.UDSStorageWriter, error) {
	f, err := s.stateTempFile(topology, state)
	if err != nil {
		return nil, err
	}
	return &onDiskUDSStorageWriter{
		s:          s,
		f:          f,
		w:          bufio.NewWriter(f),
		targetPath: s.stateFilepath(topology, state),
	}, nil
}

func (s *fsUDSStorage) stateTempFile(topology, state string) (*os.File, error) {
	return ioutil.TempFile(s.tempDirPath, s.stateFilename(topology, state))
}

func (s *fsUDSStorage) Load(topology, state string) (io.ReadCloser, error) {
	f, err := os.Open(s.stateFilepath(topology, state))
	if err != nil {
		return nil, err
	}
	return f, nil
}

var (
	fsUDSStorageFilePathRegexp = regexp.MustCompile(`^(.+)-(.+).state$`)
)

func (s *fsUDSStorage) List() (map[string][]string, error) {
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
		res[m[1]] = append(res[m[1]], m[2])
	}
	return res, nil
}

func (s *fsUDSStorage) stateFilename(topology, state string) string {
	return fmt.Sprintf("%v-%v.state", topology, state)
}

func (s *fsUDSStorage) stateFilepath(topology, state string) string {
	return filepath.Join(s.dirPath, s.stateFilename(topology, state))
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
