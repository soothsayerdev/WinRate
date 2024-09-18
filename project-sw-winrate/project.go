package main

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"fmt"
)

// Estrutura do jogador

type Player struct {
	ID int `json:"id"`
	Name string `json:"name"`
	DeckName string `json:"deck_name"`
	WinRate float64 `json:"winrate"`
	UpdateAt time.Time `json:"update_at"`
}

// Slice para armazenar os jogadores
var players []Player
var nextID int = 1

// Função para obter todos os jogadores
func getPlayers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(players)
}

// Função para adicionar um novo jogador

func AddPlayer(name, deckName string) {
    // player := Player{ID: nextID, Name: name, DeckName: deckName, WinRate: 0.5, UpdateAt: time.Now()}
    // players = append(players, player)
    // nextID++
    // fmt.Printf("Jogador '%s' adicionado com sucesso.\n", name)

	var newPlayer Player
	json.NewDecoder(r.Body).Decode(&newPlayer)
	newPlayer.ID = nextID
	newPlayer.UpdateAt = time.Now()
	players = append(players, newPlayer)
	nextID++
	json.NewEncoder(w).Encode(newPlayer)
}

// Função para atualizar o WinRate
func updateWinRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, player := range players {
		if player.ID == id {
			json.NewDecoder(r.Body).Decode(&player)
			player.UpdateAt = time.Now()
			players[index] = player
			json.NewEncoder(w).Encode(player)
			return
		}
	}
	http.Error(w, "Player not found", http.StatusNotFound)
}

// Função para remover um jogador
func deletePlayer(w http.ResponseWriter, r*http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, player := range players {
		if player.ID == id {
			players = append(players[:index], players[index +1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(players)
}

func main() {
	r := mux.NewRouter()

	// Rotas da API
	r.HandleFunc("/players", getPlayers).Methods("GET")
	r.HandleFunc("/players", AddPlayer).Methods("POST")
	r.HandleFunc("/players/{id}", updateWinRate).Methods("PUT")
	r.HandleFunc("/players/id{id}", deletePlayer).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
