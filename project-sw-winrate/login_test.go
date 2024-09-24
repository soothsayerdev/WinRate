package main

import (
	"bytes"
	"database/sql"
	//"encoding/json"


	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Function to configurate a test environment with registered user
func setup() *sql.DB {
	db, _ := sql.Open("mysql", "root:20063020@tcp(localhost:3306)/WinRate")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (email, password) VALUES (?,?)", "test@example.com", string(passwordHash))
	return db
}

// Test to func loginUser
func TestLoginUser(t *testing.T) {
	db = setup()
	defer db.Close()

	var jsonStr = []byte(`{"email":"test@example.com", "password":"senha123"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginUser)

	handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

	expected := `{"message":"Login realizado com sucesso","user_id":1}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
}