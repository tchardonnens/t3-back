package models

type Poi struct {
	Id          int     `json:"id" gorm:"primary_key"`
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Type        string  `json:"type"`
	Theme       string  `json:"theme"`
	Postcode    string  `json:"postcode"`
	Region      string  `json:"region"`
	Department  string  `json:"department"`
	City        string  `json:"city"`
	Street      string  `json:"address"`
	Website     string  `json:"website"`
	Description string  `json:"description"`
}
