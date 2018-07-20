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

func (of *OutputFilter) mkdir(path string) error {
	path = filepath.Dir(path)
	return of.fs.MkdirAll(path, 0750)
}

func (of *OutputFilter) save(info *ImageInfo, path string) error {
	extension := filepath.Ext(info.FI.Name())
	postfix := ""
	for i := 1; ; i++ {
		filename := fmt.Sprintf("%s%s%s", path, postfix, extension)
		if !of.fs.Exists(filename) {
			info.Path = filename
			return of.fs.WriteFile(filename, info.Buf, 0640)
		}
		postfix = fmt.Sprintf("_%03d", i)
	}
	return nil
}

func (of *OutputFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	year := info.Properties["year"]
	month := info.Properties["month"]
	day := info.Properties["day"]
	time := info.Properties["time"]

	path := ""
	if year == "Unknown" {
		ext := filepath.Ext(info.FI.Name())
		filename := info.FI.Name()
		path = fmt.Sprintf("/Unknown/%s", filename[0:len(filename)-len(ext)])

	} else {
		path = fmt.Sprintf("/%s/%s/%s.%s.%s-%s", year, month, year, month, day, time)
	}

	err := of.mkdir(path)
	if err == nil {
		of.save(info, path)
	} else {
		err = &FatalError{fmt.Sprintf("Failed creating %s: %v\n", info.Path, err)}
	}
	return info, err
}
