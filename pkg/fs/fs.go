package fs

import (
	"bytes"
	"io"
	"io/fs"
	iofs "io/fs"
	"os"
	"path/filepath"
)

type FullStdFS interface {
	iofs.ReadDirFS
	iofs.StatFS
	iofs.ReadFileFS
}

type FS interface {
	FullStdFS

	WriteFile(string, []byte, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
}

type advancedFS struct {
	FullStdFS

	path string
}

func ForOS(path string) FS {
	osFS := os.DirFS(path)
	return &advancedFS{
		FullStdFS: osFS.(FullStdFS),
		path:      path,
	}
}

func (afs *advancedFS) WriteFile(path string, data []byte, perm fs.FileMode) error {
	fullpath := filepath.Join(afs.path, path)
	return os.WriteFile(fullpath, data, perm)
}

func (afs *advancedFS) MkdirAll(path string, perm fs.FileMode) error {
	fullpath := filepath.Join(afs.path, path)
	return os.MkdirAll(fullpath, perm)
}

type writableFile struct {
	fs   FS
	path string
	buf  bytes.Buffer
}

func Create(dir FS, path string) (io.WriteCloser, error) {
	return &writableFile{
		fs:   dir,
		path: path,
	}, nil
}

func (wf *writableFile) Write(data []byte) (int, error) {
	return wf.buf.Write(data)
}

func (wf *writableFile) Close() error {
	return wf.fs.WriteFile(wf.path, wf.buf.Bytes(), 0644)
}
