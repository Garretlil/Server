package main

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"regexp"
	_ "regexp"
	"testing"
	_ "testing"

	_ "github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta("select id from users WHERE Name=$1 and Email=$2")).
		WithArgs("Test User2", "test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	userJSON := `{"Name": "Test User2", "Email": "test@example.com"}`
	req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer([]byte(userJSON)))
	if err != nil {
		t.Fatal(err)
	}
	// создаем тестовый ответ
	w := httptest.NewRecorder()

	Authentication(w, req, db)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response ResponseUserId
	err = json.NewDecoder(w.Body).Decode(&response)

	if err != nil {
		t.Fatal(err)
	}
	if response.Id != 1 {
		t.Errorf("Expected ID %d, got %d", 1, response.Id)
	}

	// проверяем, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestRegistration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (Name,Email) VALUES ($1, $2)")).
		WithArgs("Test User", "test@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(regexp.QuoteMeta("select id from users WHERE Name=$1 and Email=$2")).
		WithArgs("Test User", "test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// создаем тестовый запрос
	userJSON := `{"Name": "Test User", "Email": "test@example.com"}`
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(userJSON)))
	if err != nil {
		t.Fatal(err)
	}

	// создаем тестовый ответ
	w := httptest.NewRecorder()

	Registration(w, req, db)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response ResponseUserId
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Id != 1 {
		t.Errorf("Expected ID %d, got %d", 123, response.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
