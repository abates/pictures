package pictures

import (
	"path/filepath"
)

type OutputFilter struct {
	fs Filesystem
}

func NewOutputFilter(fs Filesystem) *OutputFilter {
	return &OutputFilter{
		fs: fs,
	}
}

func (of *OutputFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	dir := filepath.Dir(info.Path)
	err := of.fs.MkdirAll(dir, 0750)
	if err == nil {
		err = of.fs.WriteFile(info.Path, info.Buf, 0640)
	}
	return info, err
}
