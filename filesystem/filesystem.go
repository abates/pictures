package filesystem

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Filesystem interface {
	Open(filename string) (io.ReadWriteCloser, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Exists(path string) bool
}

type OSFilesystem struct {
	rootpath string
}

func NewOSFilesystem(rootpath string) *OSFilesystem {
	return &OSFilesystem{
		rootpath: rootpath,
	}
}

func (osf *OSFilesystem) path(filename string) string {
	return path.Join(osf.rootpath, filename)
}

func (osf *OSFilesystem) Open(filename string) (io.ReadWriteCloser, error) {
	return os.Open(osf.path(filename))
}

func (osf *OSFilesystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(osf.path(filename))
}

func (osf *OSFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(osf.path(filename), data, perm)
}

func (osf *OSFilesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(osf.path(path), perm)
}

func (osf *OSFilesystem) Exists(path string) bool {
	if _, err := os.Stat(osf.path(path)); err == nil {
		return true
	}
	return false
}
