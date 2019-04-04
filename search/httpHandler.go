package search

import (
	"MyBookLibrary/search/isbnprovider"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func ISBNHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["id"]

	model, err := isbnprovider.Query(isbn)

	if(err != nil ){
		w.WriteHeader(404)
		return
	}

	byteRepres, _ := json.Marshal(model)
	w.Write(byteRepres)
}