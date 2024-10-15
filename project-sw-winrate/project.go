// Host: localhost
// Port : 3306
// User: root

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt" // To compare password with hash

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORS MiddleWare triggered")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite requisições de qualquer origem
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Para requisições OPTIONS (preflight CORS), retorna sem passar ao próximo handler
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	// conect to mysql
	// "root:20063020soothSAYER#@tcp(127.0.0.1:3306)/WinRate"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	//defer db.Close()

	// verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão: ", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso!")

	// configuration the routes
	router := mux.NewRouter()

	// configuration of middleware CORS
	router.Use(middlewareCORS)

	router.HandleFunc("/decks", getDecks).Methods("GET", "OPTIONS")
	router.HandleFunc("/register", registerUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", loginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/decks", createDeck).Methods("POST", "OPTIONS")
	router.HandleFunc("/matches", createMatch).Methods("POST", "OPTIONS")
	router.HandleFunc("/matches/{id}", updateMatch).Methods("PUT", "OPTIONS")
	//router.HandleFunc("/status", testResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// func middlewareIPWhiteList(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r*http.Request){
// 		allowedIPS := map[string]bool{
// 			"127.0.0.1": true, // Allow localhost for testing
//             "localhost:8080": true, // Replace with actual IPs to allow
// 		}

// 		clientIP := r.RemoteAddr
// 		if _, allowed := allowedIPS[clientIP]; !allowed {
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

type User struct {
	ID       int    `json:"userID"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Deck struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	DeckName string `json:"deck_name"`
}

type Match struct {
	ID             int `json:"matchsID"`
	UserDeckID     int `json:"user_deck_id"`
	OpponentDeckID int `json:"opponent_deck_id"`
	Victories      int `json:"victories"`
	Defeats        int `json:"defeats"`
}

// func testResponse(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	response := map[string]string{
// 		"message": "Test response",
// 		"status":  "sucesso",
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	// Decodify JSON requested to struct User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Dados invalidos", http.StatusBadRequest)
		return
	}

	// Check if email already exists in database
	var existingUserID int
	err = db.QueryRow("SELECT userID FROM users WHERE email = ?", user.Email).Scan(&existingUserID)
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
	_, err = db.Exec(query, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response of sucess in register
	response := map[string]string{
		"message": "Usuario registrado com sucesso!",
	}
	json.NewEncoder(w).Encode(response)
	// fmt.Fprintf(w, "Usario registrado com sucesso!")
}

// Request email and password of user
// Verify existente of user with this email in database
// Compare to password with the password in database
// Return response of sucess/fail of login
func loginUser(w http.ResponseWriter, r *http.Request) {
	// define type of response JSON
	w.Header().Set("Content-Type", "application/json")

	var user User
	var dbPassword string

	// decodify JSON requested in requisition to struct User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"Dados invalidos"}`, http.StatusBadRequest)
		return
	}

	// consult in database to verify if user exists
	query := "SELECT userID, password FROM user WHERE email = ?"
	err = db.QueryRow(query, user.Email).Scan(&user.ID, &dbPassword)
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

func createDeck(w http.ResponseWriter, r *http.Request) {
	var deck Deck
	err := json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if deck.UserID == 0 {
		http.Error(w, "ID do usúario é necessário", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO decks (user_id, deck_name) VALUES (?, ?)", deck.UserID, deck.DeckName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Deck criado com sucesso!")

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Deck criado com sucesso!",
	}
	json.NewEncoder(w).Encode(response)
}

func calculateWinRate(victories int, defeats int) float64 {
	totalGames := victories + defeats
	if totalGames == 0 {
		return 0
	}
	return (float64(victories) / float64(totalGames)) * 100
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	winRate := calculateWinRate(match.Victories, match.Defeats)

	_, err = db.Exec("INSERT INTO matches (user_deck_id, opponent_deck_id, victories, defeats) VALUES (?, ?, ?, ?)",
		match.UserDeckID, match.OpponentDeckID, match.Victories, match.Defeats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Match criada com sucesso!",
		"winRate": winRate,
	}
	// fmt.Fprintf(w, "Partida criada com sucesso!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func updateMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	matchID := vars["id"]

	winRate := calculateWinRate(match.Victories, match.Defeats)

	_, err = db.Exec("UPDATE matches SET victories = ?, defeats = ? WHERE id = ?", match.Victories, match.Defeats, matchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Partida atualizada com sucesso!")
	response := map[string]interface{}{
		"message":  "Match atualizada com sucesso!",
		"win_rate": winRate,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getDecks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	rows, err := db.Query("SELECT id, deck_name FROM decks")
	if err != nil {
		http.Error(w, "Erro ao buscar decks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var decks []Deck
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(&deck.ID, &deck.DeckName); err != nil {
			http.Error(w, "Erro ao ler dados do deck", http.StatusInternalServerError)
			return
		}
		decks = append(decks, deck)
	}

	if err := json.NewEncoder(w).Encode(decks); err != nil {
		http.Error(w, "Erro ao serializar dados do deck", http.StatusInternalServerError)
		return
	}
}
