package main

import (
	"encoding/json"
	"net/http"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	// Получаем из JSON данные о пользователе
	var user User
	//user := User{
	//	Id:    1,
	//	Name:  "John Doe",
	//	Email: "john.doe@example.com",
	//}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userID int
	_, err = database.Exec("INSERT INTO users (Name,Email) VALUES ($1, $2)", user.Name, user.Email)
	//rows1.Next()
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
