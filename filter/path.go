package filter

import (
	"fmt"
	"path/filepath"

	"github.com/abates/pictures/filesystem"
)

type PathFilter struct {
	fs filesystem.Filesystem
}

func NewPathFilter(fs filesystem.Filesystem) *PathFilter {
	return &PathFilter{fs}
}

func (pf *PathFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	ext := filepath.Ext(info.Path)
	year := info.Properties["year"]
	month := info.Properties["month"]
	day := info.Properties["day"]
	time := info.Properties["time"]

	path := ""
	if year == "" || year == "Unknown" {
		filename := filepath.Base(info.Path)
		path = fmt.Sprintf("/Unknown/%s", filename[0:len(filename)-len(ext)])
	} else {
		path = fmt.Sprintf("/%s/%s/%s.%s.%s-%s", year, month, year, month, day, time)
	}

	postfix := ""
	for i := 1; ; i++ {
		filename := fmt.Sprintf("%s%s%s", path, postfix, ext)
		if pf.fs.Exists(filename) {
			postfix = fmt.Sprintf("_%03d", i)
		} else {
			info.Path = filename
			break
		}
	}
	return info, nil
}
