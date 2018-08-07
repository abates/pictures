package api

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/abates/pictures"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

var prefix = "/api"

type httpError struct {
	code int
	msg  string
}

func (he *httpError) Error() string { return fmt.Sprintf("%d: %s", he.code, he.msg) }

type api struct {
	server *pictures.Server
}

func (api *api) returnError(context *gin.Context, err error) {
	if pictures.IsNonFatalError(err) {
		context.String(http.StatusConflict, err.Error())
	} else if he, ok := err.(*httpError); ok {
		context.String(he.code, he.Error())
	} else if os.IsNotExist(err) {
		context.String(http.StatusNotFound, "Not Found")
	} else {
		context.Error(err)
		context.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func (api *api) getThumbnail(context *gin.Context) {
	path := context.Params.ByName("path")
	reader, err := api.server.Open(path)
	if err == nil {
		var img image.Image
		img, _, err = image.Decode(reader)
		if err == nil {
			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, imaging.Resize(img, 0, 200, imaging.Lanczos), &jpeg.Options{100})
			if err == nil {
				context.Data(http.StatusOK, "image/jpeg", buf.Bytes())
			}
		}
	}

	if err != nil {
		api.returnError(context, err)
	}
}

func (api *api) listPictures(context *gin.Context) {
	var err error

	path := context.Params.ByName("path")
	if strings.HasSuffix(path, "index.json") {
		// get directory listing
		path := strings.TrimSuffix(path, "/index.json")
		var list []pictures.FileInfo
		list, err = api.server.List(path)
		if err == nil {
			context.JSON(http.StatusOK, list)
		}
	} else {
		var buf []byte
		buf, err = api.server.ReadFile(path)
		if err == nil {
			contentType := http.DetectContentType(buf)
			context.Data(http.StatusOK, contentType, buf)
		}
	}

	if err != nil {
		api.returnError(context, err)
	}
}

func (api *api) submitPictures(context *gin.Context) {
	// receive picture
	status := http.StatusOK
	contentString := "OK"

	fh, err := context.FormFile("file")
	if err == nil {
		var file multipart.File
		file, err = fh.Open()
		if err == nil {
			var buf []byte
			buf, err = ioutil.ReadAll(file)
			contentType := http.DetectContentType(buf)
			if strings.HasPrefix(contentType, "image") {
				// ingest picture
				err = api.server.Ingest(buf)
			} else if err == nil {
				err = &httpError{http.StatusUnsupportedMediaType, fmt.Sprintf("Unsupported file type %s", contentType)}
			}
		}
	}

	if err == nil {
		context.String(status, contentString)
	} else {
		api.returnError(context, err)
	}
}

func New(server *pictures.Server) http.Handler {
	api := &api{server}

	handler := gin.New()
	handler.GET(prefix+"/thumbs/*path", api.getThumbnail)
	handler.GET(prefix+"/pictures/*path", api.listPictures)
	handler.POST(prefix+"/pictures", api.submitPictures)

	return handler
}
