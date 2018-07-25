package api

import (
	"log"
	"net/http"

	"github.com/abates/pictures/filesystem"
	"github.com/abates/pictures/filter"
	"github.com/gin-gonic/gin"
)

var prefix = "/api"

type api struct {
	fs              filesystem.Filesystem
	processingChain *filter.ProcessingChain
}

func (api *api) listPictures(context *gin.Context) {
	log.Printf("Listing pictures")
	context.String(http.StatusOK, "OK")
}

func (api *api) submitPictures(context *gin.Context) {
	log.Printf("Saving pictures")
	// receive picture
	fh, err := context.FormFile()
	if err == nil {
		var file multipart.File
		file, err = fh.Open()
		if err == nil {
			var buf []byte
			buf, err = ioutil.ReadAll(file)
			if err == nil {
				picture := filter.NewImageInfo()
				picture.Buf = buf
				// ingest picture
				api.processingChain.Push(picture)
			} else {
				log.Printf("Failed to read file: %v", err)
			}
		} else {
			log.Printf("Failed to open file: %v", err)
		}
	} else {
		log.Printf("Failed to retrieve file header: %v", err)
	}

	if err == nil {
		context.String(http.StatusOK, "OK")
	} else {
		context.Error(fmt.Errorf("Failed to receive file"))
		context.String(http.StatusInternalServerError, "Failed to receive file")
	}
}

func New(fs filesystem.Filesystem) http.Handler {
	chain := filter.NewProcessingChain()
	api := &api{fs, chain}

	handler := gin.New()
	handler.GET(prefix+"/pictures", api.listPictures)
	handler.POST(prefix+"/pictures", api.submitPictures)

	return handler
}
