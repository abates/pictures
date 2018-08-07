package pictures

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"strings"
)

type ImageDecoderFilter struct {
}

func (ff *ImageDecoderFilter) Process(info *ImageInfo) (i *ImageInfo, err error) {
	contentType := http.DetectContentType(info.Buf)
	if strings.HasPrefix(contentType, "image") {
		info.Img, _, err = image.Decode(bytes.NewReader(info.Buf))
	} else {
		err = &NonFatalError{fmt.Sprintf("%v is not an image", info.Path), false}
	}
	return info, err
}
