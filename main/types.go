package main

type Product struct {
	Id            int    `json:"Id" `
	Name          string `json:"Name"`
	Description   string `json:"Description"`
	ImageResource string `json:"ImageResource"`
	Price         int    `json:"Price"`
	Old_price     int    `json:"Old_price"`
	Amount        int    `json:"Amount"`
}
type Responses struct {
	OrderID int `json:"OrderID"`
}
type Order struct {
	Products []Product `json:"products"`
}
