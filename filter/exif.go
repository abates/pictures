package filter

import (
	"bytes"
	"fmt"

	"github.com/rwcarlsen/goexif/exif"
)

type ExifFilter struct {
}

func (ef *ExifFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	ex, err := exif.Decode(bytes.NewReader(info.Buf))
	if err == nil {
		t, err := ex.DateTime()
		if err == nil {
			info.Properties["year"] = fmt.Sprintf("%4d", t.Year())
			info.Properties["month"] = fmt.Sprintf("%02d", t.Month())
			info.Properties["day"] = fmt.Sprintf("%02d", t.Day())
			info.Properties["time"] = t.Format("15:04")
			return info, nil
		} else {
			err = &NonfatalError{fmt.Sprintf("%s failed exif date/time parse: %v", info.Path, err)}
		}
	} else {
		err = &NonfatalError{fmt.Sprintf("%s failed to decode exif: %v", info.Path, err)}
	}
	info.Properties["year"] = "Unknown"
	info.Properties["month"] = "Unknown"
	info.Properties["day"] = "Unknown"
	info.Properties["time"] = "00:00:00"
	return info, err
}
