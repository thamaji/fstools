package fstools

import (
	"io"
	"os"
)

func ReadDir(path string) ([]os.DirEntry, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	entries, err := file.ReadDir(-1)
	file.Close()
	return entries, err
}

func ReadDirFunc(path string, f func(os.DirEntry) error) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		entries, err := file.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for _, entry := range entries {
			if err := f(entry); err != nil {
				return err
			}
		}
	}
	return nil
}

func Readdir(path string) ([]os.FileInfo, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	infos, err := file.Readdir(-1)
	file.Close()
	return infos, err
}

func ReaddirFunc(path string, f func(os.FileInfo) error) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		infos, err := file.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for _, info := range infos {
			if err := f(info); err != nil {
				return err
			}
		}
	}
	return nil
}

func Readdirnames(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	names, err := file.Readdirnames(-1)
	file.Close()
	return names, err
}

func ReaddirnamesFunc(path string, f func(string) error) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		names, err := file.Readdirnames(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for _, name := range names {
			if err := f(name); err != nil {
				return err
			}
		}
	}
	return nil
}
