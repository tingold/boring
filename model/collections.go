package model


type Collections struct {
	Links       *[]Link       `json:"links"`
	Collections *[]Collection `json:"collections"`
}

type Collection struct {

	Id string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Links *[]Link `json:"links"`
	Extent *Extent `json:"extent"`
	ItemType string `json:"itemType"`
	CRS *[]string `json:"crs"` //must support http://www.opengis.net/def/crs/OGC/1.3/CRS84

}
