package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/abates/pictures/filesystem"
	"github.com/abates/pictures/filter"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input path> <output path>\n", os.Args[0])
		os.Exit(1)
	}

	disgoDbFile := path.Join(os.Args[2], "disgo.db")
	disgoFilter := filter.NewDisgoFilter()
	if _, err := os.Stat(disgoDbFile); err == nil {
		err = disgoFilter.LoadDB(disgoDbFile)
		if err != nil {
			log.Fatalf("Failed to load disgo DB: %v\n", err)
		}
	}

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			if disgoFilter.DirtyDB() {
				err := disgoFilter.SaveDB(disgoDbFile)
				if err != nil {
					log.Fatalf("Failed to save disgo DB: %v\n", err)
				}
			}
		}
	}()

	inputFilter := filter.NewImageFileFilter(filesystem.NewOSFilesystem("/"))
	outputFilter := filter.NewOutputFilter(filesystem.NewOSFilesystem(os.Args[2]))

	processingChain := filter.NewProcessingChain()
	processingChain.
		Append(inputFilter).
		Append(disgoFilter).
		Append(&filter.ExifFilter{}).
		Append(&filter.IPTCInputFilter{}).
		Append(&filter.IPTCOutputFilter{}).
		AppendLast(outputFilter)

	filepath.Walk(os.Args[1], func(path string, fi os.FileInfo, err error) error {
		if err == nil {
			info := filter.NewImageInfo()
			info.FI = fi
			info.Path = path
			processingChain.Input() <- info
		}
		return err
	})

	processingChain.Close()
	ticker.Stop()
}
