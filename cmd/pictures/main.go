//go:generate ./bundle_app.sh
package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/abates/pictures"
	"github.com/abates/pictures/api"
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
	debug := flag.Bool("d", false, "enable debugging output")
	dir := flag.String("o", "pictures", "directory to store pictures")
	port := flag.String("p", "8100", "port to serve on")
	flag.Parse()

	server, err := pictures.DefaultServer(*dir)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	server.SetDebug(*debug)

	apiHandler := api.New(server)

	http.Handle("/api/", apiHandler)
	http.HandleFunc("/", fileHandler)
	log.Printf("Listening on HTTP port: %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
