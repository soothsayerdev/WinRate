package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
    "github.com/soothsayerdev/WinRate/backend/config"
	"github.com/soothsayerdev/WinRate/backend/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	// Decodify JSON requested to struct User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Dados invalidos", http.StatusBadRequest)
		return
	}

	// Check if email already exists in database
	var existingUserID int
	err = config.DB.QueryRow("SELECT userID FROM users WHERE email = ?", user.Email).Scan(&existingUserID)
	if err == nil { // If email exists, return an error message
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	// Generate the hash of password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Erro ao gerar senha", http.StatusInternalServerError)
		return
	}

	// Insert new user into database
	query := "INSERT INTO user (email, password) VALUES (?, ?)"
	_, err = config.DB.Exec(query, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response of sucess in register
	response := map[string]string{
		"message": "Usuario registrado com sucesso!",
	}
	json.NewEncoder(w).Encode(response)
}

// Request email and password of user
// Verify existente of user with this email in database
// Compare to password with the password in database
// Return response of sucess/fail of login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// define type of response JSON
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var dbPassword string

	// decodify JSON requested in requisition to struct User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"Dados invalidos"}`, http.StatusBadRequest)
		return
	}

	// consult in database to verify if user exists
	query := "SELECT userID, password FROM user WHERE email = ?"
	err = config.DB.QueryRow(query, user.Email).Scan(&user.ID, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"Usuario não encontrado"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"Erro ao buscar o usuario"}`, http.StatusInternalServerError)
		}
		return
	}

	// Verify if password provided matches to hash in database
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	// Return a response of login sucess
	response := map[string]interface{}{
		"message": "Login realizado com sucesso",
		"userID":  user.ID,
		"success": true,
	}
	json.NewEncoder(w).Encode(response)
}
