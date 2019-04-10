package main

import (
	"MyBookLibrary/coordinator"
	"MyBookLibrary/database"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	err := database.InitDB()

	//ebookParser.StartDatabaseViewer()

	if err != nil {
		fmt.Println("Could{ not create DB: %v", err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/metadata/isbn/{id}", coordinator.ISBNHandler)
	r.HandleFunc("/upload", coordinator.UploadHandler())
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	err = http.ListenAndServe(":8080", loggedRouter)

	if err != nil {
		fmt.Println("Could not create server: %v", err)
		return
	}
}
