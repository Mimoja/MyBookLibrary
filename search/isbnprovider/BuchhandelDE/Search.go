package BuchhandelDE

import (
	"encoding/json"
	"fmt"
)

type BuchhandelDESearch struct{}

func (s BuchhandelDESearch) URL(isbn string) string {
	return fmt.Sprintf("https://www.buchhandel.de/jsonapi/productDetails/%s", isbn)

}

func (s BuchhandelDESearch) Unmarshall(data []byte) (BuchhandelDEResponse, error) {

	parsedResp := BuchhandelDEResponse{}

	err := json.Unmarshal(data, &parsedResp)

	if err != nil {
		return parsedResp, err
	}
	return parsedResp, nil
}

func (s BuchhandelDESearch) ProvierName() string {
	return "Buchhandel.de"
}
