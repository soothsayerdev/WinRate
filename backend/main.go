package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/soothsayerdev/WinRate/backend/config"
	"github.com/soothsayerdev/WinRate/backend/handlers"
	"github.com/soothsayerdev/WinRate/backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Conectar ao banco
	config.ConnectDB()
	defer config.DB.Close()

	fmt.Println("Conectado ao banco de dados com sucesso!")

	// Carregar as rotas
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// User routes
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// Deck routes
	router.HandleFunc("/decks", handlers.GetDecks).Methods("GET")
	router.HandleFunc("/decks", handlers.CreateDeck).Methods("POST")

	// Match routes
	router.HandleFunc("/matches", handlers.CreateMatch).Methods("POST")
	router.HandleFunc("/matches", handlers.UpdateMatch).Methods("PUT")

	// Iniciar servidor
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
