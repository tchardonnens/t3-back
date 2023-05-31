package models

type Site struct {
	Id          int    `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	Type        string `json:"type"`
	Postcode    string `json:"postcode"`
	Region      string `json:"region"`
	Department  string `json:"department"`
	City        string `json:"city"`
	Street      string `json:"address"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Visited     bool
	Neighbours  []*Site `gorm:"-";`
}

/*
//MARK: TO FIX
-- got the error when I parse string to float, so setting lat and lng to String --
var lat, long float64
_, err := fmt.Sscan(record[2], &lat)
if err != nil {
	log.Fatal(err)
	continue
}
*/
