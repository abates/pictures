package filter

import (
	"fmt"
	"path/filepath"

	"github.com/abates/pictures/filesystem"
)

type OutputFilter struct {
	fs filesystem.Filesystem
}

func NewOutputFilter(fs filesystem.Filesystem) *OutputFilter {
	return &OutputFilter{
		fs: fs,
	}
}

func (of *OutputFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	dir := filepath.Dir(info.Path)
	err := of.fs.MkdirAll(dir, 0750)
	if err == nil {
		err = of.fs.WriteFile(info.Path, info.Buf, 0640)
	} else {
		err = &FatalError{fmt.Sprintf("Failed creating %s: %v\n", info.Path, err)}
	}
	return info, err
}
