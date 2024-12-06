package models

type Song struct {
	Id    int    `json:"id"`
	Song  string `json:"song"`
	Group string `json:"group"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}

type Filters struct {
	Song   string
	Group  string
	Text   string
	Link   string
	Date   string
	Limit  int
	Offset int
}
