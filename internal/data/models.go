package data

import "github.com/Blue-Davinci/SocialAid/internal/database"

type Models struct {
	Program     *ProgramsManagerModel
	GeoLocation *GeoLocationsManagerModel
	HouseHold   *HouseHoldsManagerModel
	Auth        *AuthManagerModel
}

func NewModels(db *database.Queries) Models {
	return Models{
		Program:     &ProgramsManagerModel{DB: db},
		GeoLocation: &GeoLocationsManagerModel{DB: db},
		HouseHold:   &HouseHoldsManagerModel{DB: db},
		Auth:        &AuthManagerModel{DB: db},
	}
}
