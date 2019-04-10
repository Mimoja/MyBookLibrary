package coordinator

import (
	"MyBookLibrary/ebookParser"
	"MyBookLibrary/model"
	"MyBookLibrary/search"
	"fmt"
	"log"
)

func HandleFile(path string) (model.MetaDataModel,error){
	isbns, err:= ebookParser.ParseFile(path)

	if err != nil {
		log.Printf("Could not parse file %s: %v",path,  err)
		return model.MetaDataModel{}, err
	}

	for _, isbn := range isbns{
		model, err := search.Query(isbn)
		if err != nil {
			continue;
		}
		log.Printf("Found model: %v", model)
		return model, nil
	}
	return  model.MetaDataModel{}, fmt.Errorf("Could not find Metadata")
}
