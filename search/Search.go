package search

import (
	"MyBookLibrary/database"
	"MyBookLibrary/model"
	"MyBookLibrary/search/isbnprovider/BuchhandelDE"
	"fmt"
	"github.com/moraes/isbn"
	"gopkg.in/resty.v1"
	"log"
	"strings"
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

	if (len(input) == 13 && input[12:]=="X"){
		checkDigit13, _ := isbn.CheckDigit13(input)
		input = input[:12] + checkDigit13;
	}

	if (len(input) == 10 && input[9:]=="X"){
		checkDigit10 , _ := isbn.CheckDigit10(input)
		input = input[:9] + checkDigit10;
		input, _ = isbn.To13(input)
	}

	if (len(input) == 10 ){
		input, _ = isbn.To13(input)
	}

	i := model.ReformatISBN(input)

	existingModel := model.MetaDataModel{}
	err := database.GetByID("search", i, &existingModel)

	if err == nil {
		log.Println("Entry already exists:", existingModel)
		return existingModel, nil
	}

	search := BuchhandelDE.BuchhandelDESearch{}

	log.Println(search.ProvierName())
	query := search.URL(i)
	log.Println(query)

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

	if err != nil {
		return model.MetaDataModel{}, err
	}

	err = database.StoreElement(fmt.Sprintf("reponse_%s", strings.ToLower(search.ProvierName())), nil, response, &i)
	if err != nil {
		log.Printf("Could not store orignal response: %v", err)
	}

	err = database.StoreElement("search", nil, response.Model(), &i)
	if err != nil {
		log.Printf("Could not store orignal model: %v", err)
	}

	return response.Model(), nil
}
