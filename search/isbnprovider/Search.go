package isbnprovider

import (
	"MyBookLibrary/model"
	"MyBookLibrary/search/isbnprovider/BuchhandelDE"
	"fmt"
	"gopkg.in/resty.v1"
	"log"
)

type Search interface {
	Url(isbn string)
	ProviderName() string
	Unmarshall(data []byte) (SearchReponse, error)
}

type SearchReponse interface {
	Model() model.MetaDataModel
}

func Query(input string) (model.MetaDataModel, error) {
	isbn := model.ReformatISBN(input)

	search := BuchhandelDE.BuchhandelDESearch{}

	fmt.Println(search.ProvierName())
	query := search.URL(isbn)
	fmt.Println(query)

	// GET request
	resp, err := resty.R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0",
		"Accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}).Get(query)

	if err != nil {
		log.Fatalf("Could not query! %v", err)
	}
	response, err := search.Unmarshall(resp.Body())

	if err != nil{
		return model.MetaDataModel{}, err
	}
	return response.Model(), nil
}
