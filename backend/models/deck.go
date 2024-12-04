package models

import (

)

type Deck struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	DeckName  string  `json:"deck_name"`
}