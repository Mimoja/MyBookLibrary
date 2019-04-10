package ebookParser

import (
	"MyBookLibrary/model"
	"archive/zip"
	"github.com/moraes/isbn"
	"github.com/readium/r2-streamer-go/models"
	"github.com/readium/r2-streamer-go/parser"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StartDatabaseViewer() {

	go func() {
		t := time.NewTimer(5 * time.Second)
		for true {
			<-t.C
			scanDirectory()
			t = time.NewTimer(5 * time.Second)
		}
	}()
}

func scanDirectory() {
	err := filepath.Walk("ebookParser/testbooks",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			log.Println(path, info.Size())

			if info.IsDir() {
				return nil
			}

			_, err = ParseFile(path)
			return err
		})
	if err != nil {
		log.Println(err)
	}
}

func ParseFile(path string) ([]string, error) {
	publication, err := parser.Parse(path)
	if err != nil {
		log.Printf("Error while parsing epub!")
		return []string{}, err
	}

	//enc.Encode(publication)

	isbns := model.FindISBN(publication.Metadata.Identifier)

	if len(isbns) > 0 && isbn.Validate(model.ReformatISBN(isbns[0])) {
		log.Println("metadata-ISBN:", isbns[0])
		return isbns, nil
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		log.Printf("Error while parsing epub!")
		return isbns, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return isbns, err
		}

		entryBytes, err := ioutil.ReadAll(rc)

		for _, i := range model.FindISBN(string(entryBytes)) {
			i = model.ReformatISBN(i)
			valid := isbn.Validate(i)
			if valid {
				isbns = append(isbns, i)
			}
		}
	}

	if len(isbns) == 0 {
		log.Println("Could not find isbn!")
	}

	return removeDuplicatISBNs(isbns), nil
}


func removeDuplicatISBNs(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
			log.Println("ISBN:", elements[v])
		}
	}
	// Return the new slice.
	return result
}


func containsString(s []models.Link, e string) bool {
	for _, a := range s {
		if a.Href == e {
			return true
		}
	}
	return false
}
