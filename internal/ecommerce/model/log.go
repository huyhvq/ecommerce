package model

import "time"

type Log struct {
	SearchTerm string    `json:"search_term"`
	Context    []string  `json:"context"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
