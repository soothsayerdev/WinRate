package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.db

func main() {
	var err error
	// conect to mysql
	// "root:20063020soothSAYER#@tcp(127.0.0.1:3306)/WinRate"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
        os.Getenv("DB_NAME"))

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	defer db.Close()

	// verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão: ", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso!")

	// configuration the routes
	router := mux.NewRouter()

	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("login", loginUser).Methods("POST")
	router.HandleFunc("/decks", createDeck).Methods("POST")
	router.HandleFunc("matches", createMatch).Methods("POST")
	router.HandleFunc("/matches/{id}", updateMatch).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Deck struct {
	ID 					int `json:"id"`
	UserDeckID 			int `json:"user_deck_id"`
	OpponentDeckID 		int `json:"opponent_deck_id"`
	Victories         	int `json:"victories"`
	Defeats 	 		int `json:"defeats"`
}


func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO user(email, password) VALUES (?, ?)", user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error)
		return
	}
	
	fmt.Fprintf(w, "Usario registrado com sucesso!")
}

func createDeck(w http.ResponseWriter, r *http.Request) {
	var deck Deck
	err := json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO decks (user_id, deck_name) VALUES (?, ?)", deck.UserDeckID, deck.DeckName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deck criado com sucesso!")
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO matches (user_deck_id, opponent_deck_id, victories, defeats) VALUES (?, ?, ?, ?)",
					match.UserDeckID, match.OpponentDeckID, match.Victories, match.Defeats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Partida criada com sucesso!")

}

func updateMatch(w, http.ResponseWriter, r *http.Request) {
	var match Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	matchID := vars["id"]

	_, err = db.Exec("UPDATE matches SET victories = ?, defeats = ? WHERE id = ?", match.Victories, match.Defeats, matchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Partida atualizada com sucesso!")
}

