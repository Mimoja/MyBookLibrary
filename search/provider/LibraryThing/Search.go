package BuchhandelDE

import (
	"encoding/json"
	"fmt"
)

type LibraryThingSearch struct{}

func (s LibraryThingSearch) URL(isbn string) string {
	return fmt.Sprintf("", isbn)
}

func (s LibraryThingSearch) Unmarshall(data []byte) (LibraryThingReponse, error) {

	parsedResp := LibraryThingReponse{}

	err := json.Unmarshal(data, &parsedResp)

	if err != nil {
		return parsedResp, err
	}
	return parsedResp, nil
}
