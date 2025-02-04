package data

import (
	"context"
	"errors"
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

type GeoLocationsManagerModel struct {
	DB *database.Queries
}

const (
	DefaultGeoLocManDBContextTimeout = 5 * time.Second
)

var (
	ErrDuplicateGeoLocation = errors.New("geo location already exists, please choose another one")
)

type GeoLocation struct {
	ID          int32     `json:"id"`
	County      string    `json:"county"`
	SubCounty   string    `json:"sub_county"`
	Location    string    `json:"location"`
	SubLocation string    `json:"sub_location"`
	CreatedAt   time.Time `json:"created_at"`
}

func ValidateGeoLocation(v *validator.Validator, g *GeoLocation) {
	v.Check(g.County != "", "county", "must be provided")
	v.Check(g.SubCounty != "", "sub_county", "must be provided")
	v.Check(g.Location != "", "location", "must be provided")
	v.Check(g.SubLocation != "", "sub_location", "must be provided")
	v.Check(len(g.County) <= 255, "county", "must not be more than 255 bytes long")
	v.Check(len(g.SubCounty) <= 255, "sub_county", "must not be more than 255 bytes long")
	v.Check(len(g.Location) <= 255, "location", "must not be more than 255 bytes long")
	v.Check(len(g.SubLocation) <= 255, "sub_location", "must not be more than 255 bytes long")
}

// CreateNewGeoLocation() creates a new geo location in the database
// We recieve a pointer to a GeoLocation struct and return an error if the geo location already exists or
// if there was an error creating the geo location
func (m GeoLocationsManagerModel) CreateNewGeoLocation(geoLocation *GeoLocation) error {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultGeoLocManDBContextTimeout)
	defer cancel()
	// create new geo location
	geoLocationInfo, err := m.DB.CreateNewGeoLocation(ctx, database.CreateNewGeoLocationParams{
		County:      geoLocation.County,
		SubCounty:   geoLocation.SubCounty,
		Location:    geoLocation.Location,
		SubLocation: geoLocation.SubLocation,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "geolocations_sub_location_key"`:
			return ErrDuplicateGeoLocation
		default:
			return err
		}
	}
	// set the new geo location info
	geoLocation.ID = geoLocationInfo.ID
	geoLocation.CreatedAt = geoLocationInfo.CreatedAt

	return nil
}
