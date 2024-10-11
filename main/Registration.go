package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/DATA-DOG/go-sqlmock"
	"net/http"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Registration(w http.ResponseWriter, r *http.Request, database DB) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.Exec("INSERT INTO users (Name,Email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error executing query:", err) // Отладочный вывод
		fmt.Println("Executing query:", "INSERT INTO users (Name,Email) VALUES ($1, $2)", user.Name, user.Email)
		return
	}

	var userID int
	rows, err := database.Query("select id from users WHERE Name=$1 and Email=$2", user.Name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	response := ResponseUserId{Id: userID}
	json.NewEncoder(w).Encode(response)
}
