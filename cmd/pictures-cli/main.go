package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/abates/pictures/db"
	"github.com/abates/pictures/filesystem"
	"github.com/abates/pictures/filter"
)

func mkdir(path string) {
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0750)
		}

		if err != nil {
			log.Fatalf("Failed to create output directory %q: %v", path, err)
		}
	} else if !fi.IsDir() {
		log.Fatalf("%q is not a directory", path)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <input path> <output path>\n", os.Args[0])
	}

	outputPath := os.Args[2]
	mkdir(outputPath)

	dbPath := path.Join(outputPath, "db")
	mkdir(dbPath)

	imgdb, err := db.OpenBadger(dbPath)
	if err != nil {
		log.Fatalf("Failed to open image database: %v", err)
	}
	defer imgdb.Close()

	disgoDb, err := db.OpenDisgo(imgdb)
	if err != nil {
		log.Fatalf("Failed to open disgo database: %v", err)
	}
	defer disgoDb.Close()

	fs := filesystem.NewOSFilesystem(outputPath)

	processingChain := filter.NewProcessingChain()
	processingChain.
		Append(&filter.ImageDecoderFilter{}).
		Append(&filter.ExifFilter{}).
		Append(filter.NewPathFilter(fs)).
		Append(filter.NewDisgoFilter(disgoDb)).
		Append(&filter.IPTCInputFilter{}).
		Append(&filter.IPTCOutputFilter{}).
		AppendLast(filter.NewOutputFilter(fs))

	filepath.Walk(os.Args[1], func(path string, fi os.FileInfo, err error) error {
		if err == nil {
			if !fi.IsDir() {
				var buf []byte
				buf, err = ioutil.ReadFile(path)
				if err == nil {
					info := filter.NewImageInfo()
					info.Buf = buf
					info.Path = path
					processingChain.Push(info)
				}
			}
		}
		return err
	})

	processingChain.Close()
}
