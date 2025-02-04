package data

import (
	"context"
	"errors"
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

type HouseHoldsManagerModel struct {
	DB *database.Queries
}

const (
	DefaultHouseHoldManDBContextTimeout = 5 * time.Second
)

var (
	ErrGeoLocationDoesNotExist = errors.New("geo location does not exist")
	ErrProgramDoesNotExist     = errors.New("program does not exist")
)

type HouseHold struct {
	ID            int32     `json:"id"`
	ProgramID     int32     `json:"program_id"`
	GeoLocationID int32     `json:"geo_location_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
}

// ValidateHouseHold() validates the house hold struct
func ValidateHouseHold(v *validator.Validator, h *HouseHold) {
	v.Check(h.ProgramID != 0, "program_id", "must be provided")
	v.Check(h.GeoLocationID != 0, "geo_location_id", "must be provided")
	v.Check(h.Name != "", "name", "must be provided")
	v.Check(len(h.Name) <= 255, "name", "must not be more than 255 bytes long")
}

// CreateNewHouseHold() creates a new house hold in the database
// We recieve a pointer to a HouseHold struct and return an error if the house hold already exists or
// if there was an error creating the house hold
func (m HouseHoldsManagerModel) CreateNewHouseHold(houseHold *HouseHold) error {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultHouseHoldManDBContextTimeout)
	defer cancel()
	// create new house hold
	houseHoldInfo, err := m.DB.CreateNewHousehold(ctx, database.CreateNewHouseholdParams{
		ProgramID:     houseHold.ProgramID,
		GeolocationID: houseHold.GeoLocationID,
		Name:          houseHold.Name,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "households" violates foreign key constraint "households_geolocation_id_fkey"`:
			return ErrGeoLocationDoesNotExist
		case err.Error() == `pq: insert or update on table "households" violates foreign key constraint "households_program_id_fkey"`:
			return ErrProgramDoesNotExist
		default:
			return err
		}
	}
	// set the new house hold info
	houseHold.ID = houseHoldInfo.ID
	houseHold.CreatedAt = houseHoldInfo.CreatedAt
	return nil
}
