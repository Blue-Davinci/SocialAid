package data

import "github.com/Blue-Davinci/SocialAid/internal/database"

type Models struct {
	Program *ProgramsManagerModel
}

func NewModels(db *database.Queries) Models {
	return Models{
		Program: &ProgramsManagerModel{DB: db},
	}
}
