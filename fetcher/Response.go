package main

import (
	"MyBookLibrary/model"
	"github.com/grokify/html-strip-tags-go"
	"log"
	"strings"
)

type DTVResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Extent struct {
				Text   string      `json:"text"`
				Type   string      `json:"type"`
				Number interface{} `json:"number"`
			} `json:"extent"`
			Illustrations     []interface{} `json:"illustrations"`
			MainLanguages     []string      `json:"mainLanguages"`
			BiographicalNotes []interface{} `json:"biographicalNotes"`
			RelatedProducts   []struct {
				Title        string      `json:"title"`
				Identifier   string      `json:"identifier"`
				ProductGroup string      `json:"productGroup"`
				FileFormat   interface{} `json:"fileFormat"`
			} `json:"relatedProducts"`
			Edition struct {
				Text   interface{} `json:"text"`
				Number interface{} `json:"number"`
			} `json:"edition"`
			MediaFiles           []interface{} `json:"mediaFiles"`
			Medium               interface{}   `json:"medium"`
			Title                string        `json:"title"`
			UnpricedItemCode     interface{}   `json:"unpricedItemCode"`
			HasAdvertising       interface{}   `json:"hasAdvertising"`
			ContributorNotes     []string      `json:"contributorNotes"`
			PublicationFrequency interface{}   `json:"publicationFrequency"`
			NumPages             int           `json:"numPages"`
			ZisSubjectGroups     interface{}   `json:"zisSubjectGroups"`
			Contributor          interface{}   `json:"contributor"`
			SubTitle             string   `json:"subTitle"`
			ContainedItem        []interface{} `json:"containedItem"`
			Collections          []interface{} `json:"collections"`
			SubLanguages         []interface{} `json:"subLanguages"`
			MainDescriptions     []struct {
				Description  string `json:"description"`
				ContainsHTML bool   `json:"containsHTML"`
			} `json:"mainDescriptions"`
			ProductIcon string      `json:"productIcon"`
			TitleShort  interface{} `json:"titleShort"`
			Prices      []struct {
				Value            float64     `json:"value"`
				Country          string      `json:"country"`
				Currency         string      `json:"currency"`
				State            string      `json:"state"`
				Type             string      `json:"type"`
				TaxRate          string      `json:"taxRate"`
				Description      interface{} `json:"description"`
				MinQuantity      interface{} `json:"minQuantity"`
				Provisional      bool        `json:"provisional"`
				TypeQualifier    interface{} `json:"typeQualifier"`
				PriceReference   bool        `json:"priceReference"`
				FixedRetailPrice bool        `json:"fixedRetailPrice"`
			} `json:"prices"`
			PublicationDate   string `json:"publicationDate"`
			ProductType       string `json:"productType"`
			Measurements      string `json:"measurements"`
			Identifier        string `json:"identifier"`
			ProductFileFormat string `json:"productFileFormat"`
			PricesAT          []struct {
				Value            float64     `json:"value"`
				Country          string      `json:"country"`
				Currency         string      `json:"currency"`
				State            string      `json:"state"`
				Type             string      `json:"type"`
				TaxRate          string      `json:"taxRate"`
				Description      interface{} `json:"description"`
				MinQuantity      interface{} `json:"minQuantity"`
				Provisional      bool        `json:"provisional"`
				TypeQualifier    interface{} `json:"typeQualifier"`
				PriceReference   bool        `json:"priceReference"`
				FixedRetailPrice bool        `json:"fixedRetailPrice"`
			} `json:"pricesAT"`
			OesbNr           interface{} `json:"oesbNr"`
			OriginalLanguage interface{} `json:"originalLanguage"`
			CoverURL         string      `json:"coverUrl"`
			SpecialPriceText interface{} `json:"specialPriceText"`
			ProductFormID    string      `json:"productFormId"`
			OriginalTitle    string      `json:"originalTitle"`
			Publisher        string      `json:"publisher"`
			Contributors     []struct {
				Name             string `json:"name"`
				Type             string `json:"type"`
				BiographicalNote string `json:"biographicalNote"`
			} `json:"contributors"`
		} `json:"attributes"`
		Relationships struct {
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
	Included []interface{} `json:"included"`
}

func (r DTVResponse) toModel() (model.MetaDataModel){
	a := r.Data.Attributes

	desc := strings.ReplaceAll(a.MainDescriptions[0].Description, "<br />","\n");
	desc = strip.StripTags(desc)
	desc = strings.Trim(desc, " \t\n\r")


	m := model.MetaDataModel{
		ISBN:          r.Data.ID,
		Languages:     a.MainLanguages,
		Title:         a.Title,
		SubTitle:      a.SubTitle,
		OriginalTitle: a.OriginalTitle,
		PageCount:     a.NumPages,
		Description:   desc,
		Published:     a.PublicationDate,
		Type:          a.ProductType,
		CoverURL:      a.CoverURL,
	}

	for _, contrib := range a.Contributors{

		note := strings.ReplaceAll(contrib.BiographicalNote, "<br />","\n");
		note = strip.StripTags(note)
		note = strings.Trim(note, " \t\n\r")

		name := contrib.Name

		c:= model.Contibutor{
			Name:  name,
			Notes: note,
		}

		if(contrib.Type[:1]=="A"){
			m.Contibutors = append(m.Contibutors, c)
		} else {
			m.Translators = append(m.Translators, c)
		}
	}

	if desc == ""{
		log.Fatal("No description found!")
	}

	return m;
}
