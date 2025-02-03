package data

import (
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
)

type ProgramsManagerModel struct {
	DB *database.Queries
}

type Program struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
