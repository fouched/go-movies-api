package models

import "time"

type Genre struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"-"` // the dash means json ignore
	UpdatedAt time.Time `json:"-"`
}
