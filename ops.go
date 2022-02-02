package fstools

import (
	"io"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Move(oldpath, newpath string) error {
	if err := os.Rename(oldpath, newpath); err == nil {
		return nil
	}

	oldfile, err := os.OpenFile(oldpath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	oldinfo, err := oldfile.Stat()
	if err != nil {
		oldfile.Close()
		return err
	}

	if oldinfo.IsDir() {
		if err := os.Mkdir(newpath, oldinfo.Mode().Perm()); err != nil {
			oldfile.Close()
			return err
		}

		names, err := oldfile.Readdirnames(-1)
		oldfile.Close()
		if err != nil {
			return err
		}

		for _, name := range names {
			if err := Move(filepath.Join(oldpath, name), filepath.Join(newpath, name)); err != nil {
				return err
			}
		}

		return os.Remove(oldpath)
	}

	newfile, err := os.OpenFile(newpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, oldinfo.Mode().Perm())
	if err != nil {
		oldfile.Close()
		return err
	}

	_, err = io.Copy(newfile, oldfile)
	oldfile.Close()
	if err1 := newfile.Sync(); err == nil {
		err = err1
	}
	if err1 := newfile.Close(); err == nil {
		err = err1
	}
	if err != nil {
		return err
	}

	return os.Remove(oldpath)
}

func Copy(srcpath, dstpath string) error {
	srcfile, err := os.OpenFile(srcpath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	srcinfo, err := srcfile.Stat()
	if err != nil {
		srcfile.Close()
		return err
	}

	if srcinfo.IsDir() {
		if err := os.Mkdir(dstpath, srcinfo.Mode().Perm()); err != nil {
			srcfile.Close()
			return err
		}

		names, err := srcfile.Readdirnames(-1)
		srcfile.Close()
		if err != nil {
			return err
		}

		for _, name := range names {
			if err := Copy(filepath.Join(srcpath, name), filepath.Join(dstpath, name)); err != nil {
				return err
			}
		}

		return nil
	}

	dstfile, err := os.OpenFile(dstpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcinfo.Mode().Perm())
	if err != nil {
		srcfile.Close()
		return err
	}

	_, err = io.Copy(dstfile, srcfile)
	srcfile.Close()
	if err1 := dstfile.Sync(); err == nil {
		err = err1
	}
	if err1 := dstfile.Close(); err == nil {
		err = err1
	}
	return err
}

func Remove(path string) error {
	return os.RemoveAll(path)
}
