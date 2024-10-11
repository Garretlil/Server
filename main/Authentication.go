package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id    int    `json:"Id"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
}
type ResponseUserId struct {
	Id int `json:"Id"`
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	// Получаем из JSON данные о пользователе
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//проверка на существование такого пользователя и возврат Id
	var userID int

	rows, err := database.Query("select id from users WHERE Name=$1 and Email=$2", user.Name, user.Email)
	rows.Next()
	err = rows.Scan(&userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Формируем ответ
	response := ResponseUserId{Id: userID}
	json.NewEncoder(w).Encode(response)
}
