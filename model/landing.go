package model

//Landing describes the format of the root page
type Landing struct{

	Title string `json:"title"`
	Description string `json:"description"`
	Links *[]Link `json:"links"`
}