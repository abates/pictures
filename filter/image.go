package filter

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/abates/pictures/filesystem"
)

type ImageFileFilter struct {
	fs filesystem.Filesystem
}

func NewImageFileFilter(fs filesystem.Filesystem) *ImageFileFilter {
	return &ImageFileFilter{
		fs: fs,
	}
}

func (ff *ImageFileFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	var err error
	if info.FI.IsDir() {
		return nil, &NonfatalError{fmt.Sprintf("%q is a directory, not a file", info.Path)}
	}

	info.Buf, err = ff.fs.ReadFile(info.Path)
	if err == nil {
		contentType := http.DetectContentType(info.Buf)
		if strings.HasPrefix(contentType, "image") {
			info.Img, _, err = image.Decode(bytes.NewReader(info.Buf))
		} else {
			err = &NonfatalError{fmt.Sprintf("%v is not an image", info.Path)}
		}
	} else {
		err = &NonfatalError{fmt.Sprintf("Error reading %s: %v", info.Path, err)}
	}
	return info, err
}
