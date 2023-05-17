package forge

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type HTTPStatic struct {
	FileSystem      http.FileSystem
	SPAMode         bool
	NotFoundHandler func(w http.ResponseWriter, r *http.Request, httpStatic *HTTPStatic)
	Hook            func(w http.ResponseWriter, r *http.Request, fileInfo fs.FileInfo)
	Index           string
}

func HTTPStaticDefaultHook(w http.ResponseWriter, _ *http.Request, fileInfo fs.FileInfo) {
	fileExtension := filepath.Ext(fileInfo.Name())
	fileTypeHeader := mime.TypeByExtension(fileExtension)

	w.Header().Set("Content-Type", fileTypeHeader)

	switch fileExtension {
	case ".html":
		w.Header().Set("Cache-Control", "no-cache, no-store")
	default:
		w.Header().Set("Cache-Control", "max-age=0, must-revalidate")
	}
}

func (httpStatic *HTTPStatic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestedFileName := r.URL.Path

	isRequestingDirectory := strings.HasSuffix(requestedFileName, "/")
	if isRequestingDirectory {
		requestedFileName += httpStatic.Index
	}

	file, err := httpStatic.FileSystem.Open(requestedFileName)
	if err != nil {
		httpStatic.notFound(w, r)

		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	// Redirect to add forward slash
	if fileInfo.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusTemporaryRedirect)

			return
		}
	}

	httpStatic.ServeFile(w, r, file, fileInfo)
}

func (httpStatic *HTTPStatic) ServeFile(w http.ResponseWriter, r *http.Request, file http.File, fileInfo fs.FileInfo) {
	if httpStatic.Hook != nil {
		httpStatic.Hook(w, r, fileInfo)
	} else {
		HTTPStaticDefaultHook(w, r, fileInfo)
	}

	bodyBytes, _ := io.ReadAll(file)

	// ETAG Handling
	{
		h := sha256.New()
		h.Write(bodyBytes)
		w.Header().Set("etag", hex.EncodeToString(h.Sum(nil)))

		if r.Header.Get("if-none-match") == w.Header().Get("etag") {
			w.WriteHeader(http.StatusNotModified)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bodyBytes)
}

func (httpStatic *HTTPStatic) notFound(w http.ResponseWriter, r *http.Request) {
	if httpStatic.SPAMode {
		accept := r.Header.Get("Accept")
		if strings.HasPrefix(accept, "text/html") {
			file, err := httpStatic.FileSystem.Open(httpStatic.Index)
			if err != nil {
				// Use default not found handler
				http.NotFoundHandler().ServeHTTP(w, r)

				return
			}
			defer file.Close()
			fileInfo, _ := file.Stat()

			httpStatic.ServeFile(w, r, file, fileInfo)

			return
		}
	}

	if httpStatic.NotFoundHandler == nil {
		// Use default not found handler
		http.NotFoundHandler().ServeHTTP(w, r)

		return
	}

	httpStatic.NotFoundHandler(w, r, httpStatic)
}
