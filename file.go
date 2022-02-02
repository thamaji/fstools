package fstools

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func OpenFile(path string, flag int) (*os.File, error) {
	file, err := os.OpenFile(path, flag, 0644)
	if err == nil {
		return file, nil
	}

	if flag&os.O_CREATE == 0 || !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(path, flag, 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadFileFunc(path string, f func(io.Reader) error) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	err = f(file)
	file.Close()
	return err
}

func WriteFile(path string, b []byte) error {
	return WriteFileFunc(path, func(w io.Writer) error {
		_, err := io.Copy(w, bytes.NewReader(b))
		return err
	})
}

func WriteFileFunc(path string, f func(io.Writer) error) error {
	dir, name := filepath.Split(path)

	temp, err := os.CreateTemp(dir, name)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		temp, err = os.CreateTemp(dir, name)
		if err != nil {
			return err
		}
	}

	if err := temp.Chmod(0644); err != nil {
		os.Remove(temp.Name())
		return err
	}

	err = f(temp)
	if err1 := temp.Sync(); err == nil {
		err = err1
	}
	if err1 := temp.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(temp.Name())
		return err
	}

	return os.Rename(temp.Name(), path)
}
