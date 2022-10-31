package forge

import (
	"io"
	"net/http"
)

func HTTPSpaHandler(fileSystem http.FileSystem, entryPoint string, headerChanger func(http.Header)) *HTTPStatic {
	return &HTTPStatic{
		FileSystem: fileSystem,
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file, _ := fileSystem.Open(entryPoint)

			w.Header().Set("Cache-Control", "no-cache, no-store")
			headerChanger(w.Header())
			io.Copy(w, file)
		}),
	}
}
