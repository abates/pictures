package api

import "net/http"
import "github.com/gin-gonic/gin"

var prefix = "api"

type api struct {
	fs filesystem.Filesystem
}

func (api *api) listPictures(context *gin.Context) {
}

func New(fs filesystem.Filesystem) http.Handler {
	api := &api{fs}
	handler := gin.New()
	handler.GET(prefix+"/pictures", api.listPictures)

	return handler
}
