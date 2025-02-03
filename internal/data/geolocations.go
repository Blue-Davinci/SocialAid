package data

import (
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

type GeoLocationsManagerModel struct {
	DB *database.Queries
}

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

func (m GeoLocationsManagerModel) CreateNewGeoLocation(geoLocation *GeoLocation) error {
	return nil
}
