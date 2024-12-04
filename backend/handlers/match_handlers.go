package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/soothsayerdev/WinRate/backend/models"
	"github.com/soothsayerdev/WinRate/backend/config"

)

func CalculateWinRate(victories int, defeats int) float64 {
	totalGames := victories + defeats
	if totalGames == 0 {
		return 0
	}
	return (float64(victories) / float64(totalGames)) * 100
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	var match models.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, "Dados inválidos no corpo da requisição", http.StatusBadRequest)
		return
	}

	if match.UserDeckName == "" || match.OpponentDeckName == "" {
		http.Error(w, "Os nomes dos decks são obrigatórios", http.StatusBadRequest)
		return
	}

	winRate := CalculateWinRate(match.Victories, match.Defeats)

 _, err = config.DB.Exec(
		"INSERT INTO matches (user_deck_name, opponent_deck_name, victories, defeats, created_at) VALUES (?, ?, ?, ?, ?)",
		match.UserDeckName, match.OpponentDeckName, match.Victories, match.Defeats, time.Now(),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inserir partida no banco: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Match criada com sucesso!",
		"winRate": winRate,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	var match models.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// vars := mux.Vars(r)
	// matchID := vars["id"]
	print(match.ID)
	winRate := CalculateWinRate(match.Victories, match.Defeats)

	_, err = config.DB.Exec("UPDATE matches SET victories = ?, defeats = ? WHERE id = ?", match.Victories, match.Defeats, match.ID)
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

func GetMatches(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Consulta para pegar os nomes dos decks
	query := `
		SELECT m.id, m.user_deck_name, m.opponent_deck_name, m.victories, m.defeats, m.created_at,
			ud.deck_name AS user_deck_name, od.deck_name AS opponent_deck_name
		FROM matches m
		JOIN decks ud ON m.user_deck_name = ud.deck_name
		JOIN decks od ON m.opponent_deck_name = od.deck_name
		WHERE m.user_id = ?`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Erro ao buscar partidas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		if err := rows.Scan(&match.ID, &match.UserDeckName, &match.OpponentDeckName, &match.Victories, &match.Defeats, &match.CreatedAt); err != nil {
			http.Error(w, "Erro ao processar partidas: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Não formate a data aqui
		matches = append(matches, match)
	}

	w.Header().Set("Content-Type", "application/json")
	// O Go vai automaticamente serializar CreatedAt como string no formato padrão
	json.NewEncoder(w).Encode(matches)
}
