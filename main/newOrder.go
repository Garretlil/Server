package main

import (
	"encoding/json"
	"net/http"
)

func newOrder(w http.ResponseWriter, r *http.Request) {
	// Создаем данные для JSON объекта
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//start transaction
	tx, err := database.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//проверка на законченность транзакции
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	//создание в бд и получение id нового заказа
	var orderID int

	stmt, err := tx.Prepare("INSERT INTO orders DEFAULT VALUES RETURNING id")
	err = stmt.QueryRow().Scan(&orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// вставка заказа в бд
	for _, product := range order.Products {
		_, err = tx.Exec("INSERT INTO order_items (order_id, product_id) VALUES ($1, $2)", orderID, product.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Формируем ответ
	response := Responses{OrderID: orderID}
	json.NewEncoder(w).Encode(response)

}
