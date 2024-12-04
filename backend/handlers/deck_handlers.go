package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soothsayerdev/WinRate/backend/config"
	"github.com/soothsayerdev/WinRate/backend/models"
)

func CreateDeck(w http.ResponseWriter, r *http.Request) {
	var deck models.Deck
	err := json.NewDecoder(r.Body).Decode(&deck)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if deck.UserID == 0 {
		http.Error(w, "ID do usúario é necessário", http.StatusBadRequest)
		return
	}

  	_, err = config.DB.Exec("INSERT INTO decks (user_id, deck_name) VALUES (?, ?)", deck.UserID, deck.DeckName)
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

func GetDecks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	rows, err := config.DB.Query("SELECT deckID, deck_name, user_id FROM decks") // Update name column if necessary
	if err != nil {
		log.Printf("Erro ao buscar decks: %v", err)
		http.Error(w, "Erro ao buscar decks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var decks []models.Deck
	for rows.Next() {
		var deck models.Deck
		if err := rows.Scan(&deck.ID, &deck.DeckName, &deck.UserID); err != nil {
			log.Printf("Erro ao escanear dados do deck: %v", err)
			http.Error(w, "Erro ao ler dados do deck", http.StatusInternalServerError)
			return
		}
		decks = append(decks, deck)
	}

	if err := json.NewEncoder(w).Encode(decks); err != nil {
		log.Printf("Erro ao serializar dados do deck: %v", err)
		http.Error(w, "Erro ao serializar dados do deck", http.StatusInternalServerError)
		return
	}
}
