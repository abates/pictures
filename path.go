package pictures

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
)

type PathFilter struct {
	fs Filesystem
}

func NewPathFilter(fs Filesystem) *PathFilter {
	return &PathFilter{fs}
}

func (pf *PathFilter) Process(info *ImageInfo) (ii *ImageInfo, err error) {
	contentType := http.DetectContentType(info.Buf)
	extensions, _ := mime.ExtensionsByType(contentType)
	if len(extensions) > 0 {
		ext := extensions[0]

		year := info.Time.Year()
		month := info.Time.Month()
		day := info.Time.Day()
		time := info.Time.Format("15:04")
		hashBuf, _ := info.Hash.MarshalBinary()
		//hash := strings.TrimSuffix(base64.StdEncoding.EncodeToString(hashBuf), "=")
		hash := base64.RawURLEncoding.EncodeToString(hashBuf)

		info.Path = fmt.Sprintf("/%04d/%02d/%04d.%02d.%02d-%s_%v%s", year, month, year, month, day, time, hash, ext)
	} else {
		err = &NonFatalError{fmt.Sprintf("No known extension for %v", contentType), false}
	}
	return info, err
}
