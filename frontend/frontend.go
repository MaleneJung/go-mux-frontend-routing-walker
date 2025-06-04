package frontend

import (
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func MuxFrontendWalker(serveMux *http.ServeMux, baseRoute string, baseDir string, logging bool) error {
	return filepath.Walk(baseDir, func(filePath string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		filePath = strings.ReplaceAll(filePath, "\\", "/")

		fileRoute, _ := strings.CutPrefix(filePath, baseDir)

		fileHandler := func(w http.ResponseWriter, r *http.Request) {

			fileBuffer, err := os.ReadFile(filePath)
			if err != nil {
				w.WriteHeader(404)
				if logging {
					log.Fatal(err)
				}
				return
			}

			mimeType := mime.TypeByExtension(fileRoute)
			w.Header().Add("Content-Type", mimeType)

			w.WriteHeader(200)
			w.Write(fileBuffer)

		}

		if indexRoute, isIndexMarkup := strings.CutSuffix(fileRoute, "/index.html"); isIndexMarkup {
			if logging {
				log.Println("Registered Route \"" + baseRoute + "/" + indexRoute + "/" + "\" for \"" + filePath + "\"")
			}
			serveMux.HandleFunc(baseRoute+indexRoute+"/", fileHandler)
		}
		if logging {
			log.Println("Registered Route \"" + baseRoute + "/" + fileRoute + "\" for \"" + filePath + "\"")
		}
		serveMux.HandleFunc(baseRoute+fileRoute, fileHandler)

		return nil

	})
}
