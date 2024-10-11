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

func Authentication(w http.ResponseWriter, r *http.Request, database DB) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var userID int

	rows, err := database.Query("select id from users WHERE Name=$1 and Email=$2", user.Name, user.Email)

	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	response := ResponseUserId{Id: userID}
	json.NewEncoder(w).Encode(response)
}
