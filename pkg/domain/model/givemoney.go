package model

import "time"

// GiveMoney represents the give_money entity
type GiveMoney struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Nickname  string     `json:"nickname"`
	Figure    int        `json:"figure"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}