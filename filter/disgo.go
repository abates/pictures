package filter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abates/disgo"
)

type DisgoFilter struct {
	db    *disgo.DB
	dirty bool
}

func NewDisgoFilter() *DisgoFilter {
	return &DisgoFilter{
		db:    disgo.New(),
		dirty: false,
	}
}

func (df *DisgoFilter) LoadDB(filename string) error {
	file, err := os.Open(filename)
	if err == nil {
		err = df.db.Load(file)
		if err == nil {
			df.dirty = false
		}
	}
	return err
}

func (df *DisgoFilter) SaveDB(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err == nil {
		err = df.db.Save(file)
		if err == nil {
			df.dirty = false
		}
	}
	return err
}

func (df *DisgoFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	var hash disgo.PHash
	var err error

	if hashString, found := info.Properties["disgo"]; found {
		var phash uint64
		phash, err = strconv.ParseUint(hashString, 10, 64)
		hash = disgo.PHash(phash)
	} else {
		hash, err = disgo.Hash(info.Img)
		info.Properties["disgo"] = fmt.Sprintf("%d", hash)
	}

	if err == nil {
		duplicates, _ := df.db.SearchByHash(hash, 3)
		found := false
		for _, h := range duplicates {
			if h == hash {
				found = true
			}
		}
		if !found {
			df.db.AddHash(hash)
			df.dirty = true
		}
		if len(duplicates) == 0 {
			return info, nil
		}
		err = &NonfatalError{fmt.Sprintf("duplicate image %v", info.Path)}
	} else {
		err = &NonfatalError{fmt.Sprintf("Failed to hash %s: %v\n", info.Path, err)}
	}
	return nil, err
}

func (df *DisgoFilter) DirtyDB() bool {
	return df.dirty
}
