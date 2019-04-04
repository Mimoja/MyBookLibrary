package BuchhandelDE

import (
	"encoding/json"
	"fmt"
)

type ArchiveOrgSearch struct{}

func (s ArchiveOrgSearch) URL(isbn string) string {
	return fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json", isbn)

}

func (s ArchiveOrgSearch) Unmarshall(data []byte) (ArchiveOrgResponse, error) {

	parsedResp := ArchiveOrgResponse{}

	err := json.Unmarshal(data, &parsedResp)

	if err != nil {
		return parsedResp, err
	}
	return parsedResp, nil
}
