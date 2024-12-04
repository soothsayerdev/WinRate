package models

import (
	"time"
)

type Match struct {
	ID               int       	  `json:"match_id"`
	UserDeckName     string       `json:"user_deck_id"`
	OpponentDeckName string       `json:"opponent_deck_id"`
	Victories        int          `json:"victories"`
	Defeats          int          `json:"defeats"`
	CreatedAt        time.Time    `json:"created_at"`
}