package pictures

import (
	"bytes"

	"github.com/rwcarlsen/goexif/exif"
)

type ExifFilter struct {
}

func (ef *ExifFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	ex, err := exif.Decode(bytes.NewReader(info.Buf))
	if err == nil {
		t, err := ex.DateTime()
		if err == nil {
			info.Time = t
			return info, nil
		}
	}

	if err != nil {
		err = &NonFatalError{err.Error(), true}
	}
	return info, err
}
