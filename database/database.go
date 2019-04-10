package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
)

var client *elastic.Client

func InitDB() error {
	var err error
	client, err = elastic.NewClient()

	if err != nil {
		log.Fatal("Could not connect to elastic")
		return err
	}
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		log.Fatal("Could not connect to elastic")
		return err
	}
	log.Printf("Elasticsearch version %s", esversion)

	return nil
}

func StoreElement(index string, typeString *string, entry interface{}, id *string) error {
	is := client.Index().BodyJson(entry)
	return store(is, index, typeString, id)
}

func StoreJSON(index string, typeString *string, entry string, id *string) error {
	is := client.Index().BodyString(entry)
	return store(is, index, typeString, id)
}

func store(is *elastic.IndexService, index string, typeString *string, id *string) error {
	is = is.Index(index)

	if typeString != nil {
		is = is.Type(*typeString)
	} else {
		is.Type(index)
	}

	if id != nil {
		is = is.Id(*id)
	}

	put1, err := is.Do(context.Background())
	if err != nil {
		// Handle error
		log.Println("Could not execute elastic search")
		return err
	}
	log.Printf("Indexed %s to index %s, type %s", put1.Id, put1.Index, put1.Type)
	return nil
}

func Exists(index string, id string) (bool, error, *elastic.GetResult) {
	get, err := client.Get().
		Index(index).
		Id(id).
		Do(context.Background())

	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			return false, nil, get
		case elastic.IsTimeout(err):
			log.Printf("Timeout retrieving document: %v", err)
			return false, err, get
		case elastic.IsConnErr(err):
			log.Printf("Connection problem: %v", err)
			return false, err, get
		default:
			log.Printf("Unknown error: %v", err)
			return false, err, get
		}
	}
	return true, nil, get
}

func GetByID(index string, id string, target interface{}) error {

	found, err, entry := Exists(index, id)
	if err != nil {
		log.Printf("Error fetching item: %v", err)
		return err
	}

	if !found {
		return fmt.Errorf("Element %s not found", id)
	}

	data, err := entry.Source.MarshalJSON()
	if err != nil {
		log.Printf("Could not get old entry from elastic: %v", err)
		return err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		log.Printf("Could unmarshall old entry from elastic: %v", err)
		return err
	}
	return nil
}

func AppendElement(index string, typeString *string, entry interface{}, id *string, newElement interface{}) error {
	ctx := context.Background()
	scriptString := fmt.Sprintf("ctx._source.%s += newEntry", newElement)
	script := elastic.NewScript(scriptString).
		Params(map[string]interface{}{"newEntry": newElement})

	_, err := client.Update().
		Index(index).
		Type(*typeString).
		Id(*id).
		Script(script).
		Do(ctx)
	if err != nil {
		log.Printf("Error while appending to %s", entry)
		return err
	}
	log.Printf("Appent %s to %s", newElement, entry)
	return nil
}

func UpdateElement(index string, typeString *string, entry interface{}, id *string) error {
	ctx := context.Background()
	scriptString := fmt.Sprintf("ctx._source.%s = newEntry", entry)
	script := elastic.NewScript(scriptString).
		Params(map[string]interface{}{"newEntry": entry})
	_, err := client.Update().
		Index(index).
		Type(*typeString).
		Id(*id).
		Script(script).
		Do(ctx)
	if err != nil {
		log.Printf("Error while updating %s", entry)
		return err
	}
	log.Printf("updated %s", entry)
	return nil
}

func Search(index string, terms map[string]string) error {
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index(index).
		Query(termQuery).
		Sort("user", true).
		From(0).Size(10).
		Pretty(false).
		Do(context.Background())

	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			return err
		case elastic.IsTimeout(err):
			log.Printf("Timeout retrieving document: %v", err)
			return err
		case elastic.IsConnErr(err):
			log.Printf("Connection problem: %v", err)
			return err
		default:
			log.Printf("Unknown error: %v", err)
			return err
		}
	}

	log.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	return nil
}
