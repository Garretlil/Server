package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getCatalog(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from products ")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	products := []map[string]interface{}{}
	for rows.Next() {
		var Id int
		var Name, Description, ImageResource string
		var Price, Old_price int
		var Amount int
		err := rows.Scan(&Id, &Name, &Description, &ImageResource, &Price, &Old_price, &Amount)

		if err != nil {
			fmt.Println(err)
			continue
		}
		product := map[string]interface{}{
			"Id":            Id,
			"Name":          Name,
			"Description":   Description,
			"ImageResource": ImageResource,
			"Price":         Price,
			"Old_price":     Old_price,
			"Amount":        Amount,
		}
		products = append(products, product)
	}

	jsonData, err := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, string(jsonData))
}
