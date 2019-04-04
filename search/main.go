package main

import (
	"MyBookLibrary/search/provider"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search/{id}", productsHandler)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":8080", loggedRouter)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["id"]

	// explore response object
	/*log.Printf("Error: %v", err)
	log.Printf("Response Status Code: %v", resp.StatusCode())
	log.Printf("Response Status: %v", resp.Status())
	log.Printf("Response Time: %v", resp.Time())
	log.Printf("Response Received At: %v", resp.ReceivedAt())
	log.Printf("Response Body: %v", resp)     // or res
	*/

	model := provider.Query(isbn)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	//enc.Encode(parsedResp)
	//enc.Encode(parsedResp.toModel())

	byteRepres, _ := json.Marshal(model)
	w.Write(byteRepres)
}
