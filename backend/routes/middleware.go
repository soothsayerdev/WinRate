package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soothsayerdev/WinRate/backend/handlers"
)

func RegisterRoutes(router *mux.Router) {
	// configuration of middleware CORS
	router.Use(MiddlewareCORS)

	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/decks", handlers.GetDecks).Methods("GET", "OPTIONS")
	router.HandleFunc("/decks", handlers.CreateDeck).Methods("POST", "OPTIONS")

	router.HandleFunc("/matches", handlers.CreateMatch).Methods("POST", "OPTIONS")
	router.HandleFunc("/matches", handlers.GetMatches).Methods("GET", "OPTIONS")
	router.HandleFunc("/matches/{id}", handlers.UpdateMatch).Methods("PUT", "OPTIONS")
	//router.HandleFunc("/status", testResponse).Methods("GET")

}

func MiddlewareCORS(next http.Handler) http.Handler {
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
