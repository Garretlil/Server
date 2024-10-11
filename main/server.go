package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "html/template"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "door"
	dbname   = "testServer"
)

var database *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	http.HandleFunc("/", getCatalog)
	http.HandleFunc("/orders", newOrder)
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		Authentication(w, r, db)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		Registration(w, r, db) // db передается как замыкание
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "About Page")
	})
	//start server
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
