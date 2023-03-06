package forge

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type HTTPStatic struct {
	FileSystem      http.FileSystem
	NotFoundHandler http.Handler
	CacheControl    string
	Index           string
}

func (httpStatic *HTTPStatic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if httpStatic.Index == "" {
		httpStatic.Index = "index.html"
	}

	requestedFileName := r.URL.Path
	isRequestingDirectory := strings.HasSuffix(requestedFileName, "/")
	if isRequestingDirectory {
		requestedFileName += httpStatic.Index
	}

	file, err := httpStatic.FileSystem.Open(requestedFileName)
	if err != nil {
		correctNotFoundHandler(httpStatic.NotFoundHandler).ServeHTTP(w, r)
		return
	}
	defer file.Close()

	fileTypeHeader := mime.TypeByExtension(filepath.Ext(requestedFileName))

	w.Header().Set("Content-Type", fileTypeHeader)
	if httpStatic.CacheControl != "" {
		w.Header().Set("Cache-Control", httpStatic.CacheControl)
	}

	bodyBytes, _ := io.ReadAll(file)
	h := sha256.New()
	h.Write(bodyBytes)
	w.Header().Set("etag", hex.EncodeToString(h.Sum(nil)))

	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}
