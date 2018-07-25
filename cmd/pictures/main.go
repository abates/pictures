//go:generate ./bundle_app.sh
package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/abates/pictures/api"
	"github.com/abates/pictures/filesystem"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	file := strings.TrimLeft(r.URL.Path, "/")
	if file == "" {
		file = "index.html"
	}

	var data []byte
	fi, err := AssetInfo(file)
	if err == nil {
		data, _ = Asset(file)
	} else {
		fi, _ = AssetInfo("index.html")
		data, _ = Asset("index.html")
	}
	http.ServeContent(w, r, file, fi.ModTime(), bytes.NewReader(data))
}

func main() {
	dir := flag.String("d", "pictures", "directory to store pictures")
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	apiHandler := api.New(filesystem.NewOSFilesystem(*dir))

	http.Handle("/api/", apiHandler)
	http.HandleFunc("/", fileHandler)
	log.Printf("Listening on HTTP port: %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
