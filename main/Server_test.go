package main

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	_ "regexp"
	"testing"
	_ "testing"

	_ "github.com/stretchr/testify/assert"
)

func TestGetCatalog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Name", "Description", "ImageResource", "Price", "Old_price", "Amount"}).
		AddRow(1, "Product 1", "Description 1", "image1.jpg", 100, 120, 10).
		AddRow(2, "Product 2", "Description 2", "image2.jpg", 200, 250, 5).
		AddRow(3, "Product 3", "Description 3", "image3.jpg", 700, 600, 50)

	mock.ExpectQuery(regexp.QuoteMeta("select * from products")).
		WillReturnRows(rows)

	req, err := http.NewRequest(http.MethodGet, "/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetCatalog(w, req, db)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var products []map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&products)
	if err != nil {
		t.Fatal(err)
	}

	expectedProducts := []map[string]interface{}{
		{
			"Id":            1,
			"Name":          "Product 1",
			"Description":   "Description 1",
			"ImageResource": "image1.jpg",
			"Price":         100,
			"Old_price":     120,
			"Amount":        10,
		},
		{
			"Id":            2,
			"Name":          "Product 2",
			"Description":   "Description 2",
			"ImageResource": "image2.jpg",
			"Price":         200,
			"Old_price":     250,
			"Amount":        5,
		},
		{
			"Id":            3,
			"Name":          "Product 3",
			"Description":   "Description 3",
			"ImageResource": "image3.jpg",
			"Price":         700,
			"Old_price":     600,
			"Amount":        50,
		},
	}

	if !reflect.DeepEqual(products, expectedProducts) {
		t.Errorf("Expected products: %v, got: %v", expectedProducts, products)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

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

	userJSON := `{"Name": "Test User", "Email": "test@example.com"}`
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(userJSON)))
	if err != nil {
		t.Fatal(err)
	}

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
