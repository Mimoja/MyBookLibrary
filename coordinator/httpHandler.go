package coordinator

import (
	"MyBookLibrary/search"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 20 * 1024 * 1024 // 2 mb
const uploadPath = "."

func ISBNHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["id"]

	model, err := search.Query(isbn)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	byteRepres, _ := json.Marshal(model)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteRepres)
}

func UploadHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest, err)
			return
		}

		// parse and validate file and post parameters
		file, header, err := r.FormFile("data")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest, err)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest, err)
			return
		}

		fileName := randToken(12)

		newPath := filepath.Join(uploadPath, fileName+filepath.Ext(header.Filename))

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError, err)
			return
		}

		defer newFile.Close() // idempotent, okay to call twice
		defer os.Remove(newPath)
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError, err)
			return
		}


		model, err :=  HandleFile(newPath)
		if err != nil {
			renderError(w, "METADATA_NOT_FOUND", http.StatusGone, err)
			return
		}

		byteRepres, _ := json.Marshal(model)

		w.Header().Set("Content-Type", "application/json")
		w.Write(byteRepres)


		return
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int, e error) {
	w.WriteHeader(statusCode)
	log.Printf("Error: %v", e)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
