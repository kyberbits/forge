package main

import (
	"net/http"
	"time"

	"github.com/kyberbits/forge/forge"
)

func main() {
	fileSystem := http.Dir("examples/static/public")

	static := &forge.HTTPStatic{
		FileSystem: fileSystem,
		SPAMode:    true,
		Index:      "index.html",
	}

	server := http.Server{
		Addr:              "127.0.0.1:8000",
		ReadHeaderTimeout: time.Second * 5,
		Handler:           static,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
