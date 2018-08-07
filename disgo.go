package pictures

import (
	"fmt"
	"os"

	"github.com/abates/disgo"
)

type DisgoFilter struct {
	db *DisgoDB
}

func NewDisgoFilter(db *DisgoDB) *DisgoFilter {
	return &DisgoFilter{
		db: db,
	}
}

func (df *DisgoFilter) LoadDB(filename string) error {
	file, err := os.Open(filename)
	if err == nil {
		err = df.db.Load(file)
	}
	return err
}

func (df *DisgoFilter) SaveDB(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err == nil {
		err = df.db.Save(file)
	}
	return err
}

func (df *DisgoFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	var err error

	if info.Hash == 0 {
		info.Hash, err = disgo.Hash(info.Img)
	}

	if err == nil {
		duplicates, _ := df.db.SearchByHash(info.Hash, 3)
		found := false
		for _, h := range duplicates {
			if h == info.Hash {
				found = true
			}
		}

		if !found {
			df.db.AddHash(info.Hash, info.Path)
		} else {
			err = &NonFatalError{fmt.Sprintf("duplicate image %v", info.Path), false}
			info = nil
		}
	}
	return info, err
}
