package common

type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID    string `json:"id"`
	Items []Item `json:"items"`
}
