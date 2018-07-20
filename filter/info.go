package filter

import (
	"image"
	"os"
)

type ImageFilter interface {
	Process(*ImageInfo) (*ImageInfo, error)
}

type ImageInfo struct {
	FI         os.FileInfo
	Path       string
	Buf        []byte
	Img        image.Image
	Properties map[string]string
	properties map[string]interface{}
}

func NewImageInfo() *ImageInfo {
	return &ImageInfo{
		Properties: make(map[string]string),
		properties: make(map[string]interface{}),
	}
}
