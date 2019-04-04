package model

type MetaDataModel struct {
	ISBN          string
	Languages     []string
	Title         string
	SubTitle      string
	OriginalTitle string
	Contibutors   []Contibutor
	Translators   []Contibutor
	PageCount     int
	Description   string
	Published     string
	Type          string
	Cover         string
}

type Contibutor struct {
	Name  string
	Notes string
}
