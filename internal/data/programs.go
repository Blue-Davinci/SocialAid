package data

import (
	"context"
	"errors"
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

type ProgramsManagerModel struct {
	DB *database.Queries
}

const (
	DefaultProgramManDBContextTimeout = 5 * time.Second
)

var (
	ErrDuplicateProgram = errors.New("program's name already exists, please choose another one")
)

type Program struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ValidateProgram(v *validator.Validator, p *Program) {
	v.Check(p.Name != "", "name", "must be provided")
	v.Check(p.Category != "", "category", "must be provided")
	v.Check(p.Description != "", "description", "must be provided")
	v.Check(len(p.Name) <= 255, "name", "must not be more than 500 bytes long")
	v.Check(len(p.Category) <= 255, "category", "must not be more than 500 bytes long")
	v.Check(len(p.Description) <= 1000, "description", "must not be more than 1000 bytes long")
}

func (m ProgramsManagerModel) CreateNewProgram(program *Program) error {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultProgramManDBContextTimeout)
	defer cancel()
	programInfo, err := m.DB.CreateNewProgram(ctx, database.CreateNewProgramParams{
		Name:        program.Name,
		Category:    program.Category,
		Description: program.Description,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "programs_name_key"`:
			return ErrDuplicateProgram
		default:
			return err
		}
	}
	// set the new program info
	program.ID = programInfo.ID
	program.CreatedAt = programInfo.CreatedAt
	program.UpdatedAt = programInfo.UpdatedAt
	// no error so return nil
	return nil
}
