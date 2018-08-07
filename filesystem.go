package pictures

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
}

type Tag string

type fileInfo struct {
	meta *Metadata
}

func (fi *fileInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(fi.meta)
}

// Name is the base name of the file
func (fi *fileInfo) Name() string { return fi.meta.Name }

// Size is the length in bytes for regular files; system-dependent for others
func (fi *fileInfo) Size() int64 { return fi.meta.Size }

// Mode are the file mode bits
func (fi *fileInfo) Mode() os.FileMode { return fi.meta.Mode }

// ModTime is the modification time
func (fi *fileInfo) ModTime() time.Time { return fi.meta.ModTime }

// IsDir is the abbreviation for Mode().IsDir()
func (fi *fileInfo) IsDir() bool { return fi.meta.dir }

type Metadata struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	Dir     bool        `json:"dir"`
	Tags    []Tag       `json:"tags"`
}

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type Filesystem interface {
	Open(filename string) (io.ReadWriteCloser, error)
	ReadFile(filename string) ([]byte, error)
	List(filename string) ([]FileInfo, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Exists(path string) bool
}

type OSFilesystem struct {
	rootpath string
}

func OpenOSFilesystem(rootpath string, createMode os.FileMode) (fs *OSFilesystem, err error) {
	fs = NewOSFilesystem(rootpath)

	if _, err = os.Stat(rootpath); err != nil {
		if os.IsNotExist(err) {
			err = nil
			err = os.MkdirAll(rootpath, createMode)
		}
	}
	return fs, err
}

func NewOSFilesystem(rootpath string) *OSFilesystem {
	return &OSFilesystem{
		rootpath: rootpath,
	}
}

func (osf *OSFilesystem) path(filename string) string {
	filename = filepath.Clean(filename)
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

func (osf *OSFilesystem) List(path string) (fis []FileInfo, err error) {
	files, err := ioutil.ReadDir(osf.path(path))
	if err == nil {
		for _, fi := range files {
			fis = append(fis, &fileInfo{
				&Metadata{
					Name:    fi.Name(),
					Size:    fi.Size(),
					Mode:    fi.Mode(),
					ModTime: fi.ModTime(),
					Dir:     fi.IsDir(),
				},
			})
		}
	}
	return
}
