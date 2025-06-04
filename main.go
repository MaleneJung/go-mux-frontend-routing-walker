package main

import (
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func MuxFrontendWalker(serveMux *http.ServeMux, baseRoute string, baseDir string) error {
	return filepath.Walk(baseDir, func(filePath string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		fileRoute, _ := strings.CutPrefix(filePath, baseDir)

		fileHandler := func(w http.ResponseWriter, r *http.Request) {

			fileBuffer, err := os.ReadFile(filePath)
			if err != nil {
				w.WriteHeader(404)
				return
			}

			mimeType := mime.TypeByExtension(fileRoute)
			w.Header().Add("Content-Type", mimeType)

			w.WriteHeader(200)
			w.Write(fileBuffer)

		}

		if indexRoute, isIndexMarkup := strings.CutSuffix(fileRoute, "/index.html"); isIndexMarkup {
			serveMux.HandleFunc(baseRoute+"/"+indexRoute+"/", fileHandler)
		}
		serveMux.HandleFunc(baseRoute+"/"+fileRoute, fileHandler)

		return nil

	})
}
