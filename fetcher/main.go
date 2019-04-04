package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/resty.v1"
	"log"
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
	id := vars["id"]


	query := fmt.Sprintf("https://www.buchhandel.de/jsonapi/productDetails/%s", id)
	fmt.Println(query)


	// GET request
	resp, err := resty.R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0",
		"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}).Get(query)

	// explore response object
	/*log.Printf("Error: %v", err)
	log.Printf("Response Status Code: %v", resp.StatusCode())
	log.Printf("Response Status: %v", resp.Status())
	log.Printf("Response Time: %v", resp.Time())
	log.Printf("Response Received At: %v", resp.ReceivedAt())
	log.Printf("Response Body: %v", resp)     // or res
	*/

	parsedResp := DTVResponse{}

	err = json.Unmarshal(resp.Body(), &parsedResp)

	if err != nil {
		log.Printf("Could not parse DTV response!: %v", err)
		w.WriteHeader(404)
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	//enc.Encode(parsedResp)
	//enc.Encode(parsedResp.toModel())

	byteRepres, err := json.Marshal(parsedResp.toModel())
	w.Write(byteRepres)
}