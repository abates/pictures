package pictures

import (
	"image"
	"time"

	"github.com/abates/disgo"
)

type ImageFilter interface {
	Process(*ImageInfo) (*ImageInfo, error)
}

type ImageInfo struct {
	Path string      `json:"path"`
	Buf  []byte      `json:"-"`
	Img  image.Image `json:"-"`
	Time time.Time   `json:"time"`
	Tags []string    `json:"tags"`
	Hash disgo.PHash `json:"hash"`
}

func NewImageInfo() *ImageInfo {
	return &ImageInfo{}
}
