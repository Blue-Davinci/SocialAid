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
	ErrHouseHoldDoesNotExist   = errors.New("house hold does not exist")
	ErrHouseHoldAlreadyExists  = errors.New("house hold already exists and has a head")
	ErrHouseHoldMemberExists   = errors.New("house hold member already exists")
)

type HouseHold struct {
	ID            int32     `json:"id"`
	ProgramID     int32     `json:"program_id"`
	GeoLocationID int32     `json:"geo_location_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
}

type HouseHoldHead struct {
	ID          int32     `json:"id"`
	HouseHoldID int32     `json:"house_hold_id"`
	Name        string    `json:"name"`
	NationalID  string    `json:"national_id"`
	PhoneNumber string    `json:"phone_number"`
	Age         int32     `json:"age"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HouseHoldMember struct {
	ID          int32     `json:"id"`
	HouseHoldID int32     `json:"house_hold_id"`
	Name        string    `json:"name"`
	Age         int32     `json:"age"`
	Relation    string    `json:"relation"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ValidateHouseHold() validates the house hold struct
func ValidateHouseHold(v *validator.Validator, h *HouseHold) {
	v.Check(h.ProgramID != 0, "program_id", "must be provided")
	v.Check(h.GeoLocationID != 0, "geo_location_id", "must be provided")
	v.Check(h.Name != "", "name", "must be provided")
	v.Check(len(h.Name) <= 255, "name", "must not be more than 255 bytes long")
}

func ValidateHouseHoldHead(v *validator.Validator, h *HouseHoldHead) {
	v.Check(h.HouseHoldID != 0, "house_hold_id", "must be provided")
	v.Check(h.Name != "", "name", "must be provided")
	v.Check(len(h.Name) <= 255, "name", "must not be more than 255 bytes long")
	v.Check(h.NationalID != "", "national_id", "must be provided")
	v.Check(len(h.NationalID) <= 15, "national_id", "must not be more than 255 bytes long")
	v.Check(h.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(h.Age != 0, "age", "must be provided")
}

func ValidateHouseHoldMember(v *validator.Validator, h *HouseHoldMember) {
	v.Check(h.HouseHoldID != 0, "house_hold_id", "must be provided")
	v.Check(h.Name != "", "name", "must be provided")
	v.Check(len(h.Name) <= 255, "name", "must not be more than 255 bytes long")
	v.Check(h.Age != 0, "age", "must be provided")
	v.Check(h.Relation != "", "relation", "must be provided")
	v.Check(len(h.Relation) <= 255, "relation", "must not be more than 255 bytes long")
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

// CreateNewHouseholdHead() creates a new house hold head in the database
// We recieve a pointer to a HouseHold struct and return an error if the house hold already has a head or
// if there was an error creating the house hold head.
// We encrypt the phone number before storing it in the database
func (m HouseHoldsManagerModel) CreateNewHouseholdHead(houseHoldHead *HouseHoldHead, encryption_key string) error {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultHouseHoldManDBContextTimeout)
	defer cancel()
	// handle phone number encryption
	// decrypt our hex encoded key
	decodedKey, err := DecodeEncryptionKey(encryption_key)
	if err != nil {
		return err
	}
	// encrypt and set the password
	encryptedPhoneNumber, err := EncryptData(houseHoldHead.PhoneNumber, decodedKey)
	if err != nil {
		return err
	}
	// create new house hold head
	houseHoldHeadInfo, err := m.DB.CreateNewHouseholdHead(ctx, database.CreateNewHouseholdHeadParams{
		HouseholdID: houseHoldHead.HouseHoldID,
		Name:        houseHoldHead.Name,
		NationalID:  houseHoldHead.NationalID,
		PhoneNumber: encryptedPhoneNumber,
		Age:         houseHoldHead.Age,
	})
	if err != nil {
		switch {
		// check if the house hold already has a head, if so return an error
		case err.Error() == `pq: insert or update on table "household_heads" violates foreign key constraint "household_heads_household_id_fkey"`:
			return ErrHouseHoldDoesNotExist
		case err.Error() == `pq: duplicate key value violates unique constraint "household_heads_household_id_key"`:
			return ErrHouseHoldAlreadyExists
		default:
			return err
		}
	}
	// set the new house hold head info
	houseHoldHead.ID = houseHoldHeadInfo.ID
	houseHoldHead.CreatedAt = houseHoldHeadInfo.CreatedAt
	houseHoldHead.UpdatedAt = houseHoldHeadInfo.UpdatedAt
	// return nil if everything is successful
	return nil
}

// CreateNewHouseholdMember() creates a new house hold member in the database
// We recieve a pointer to a HouseHold struct and return an error if the house hold already has a head or
// if there was an error creating the house hold head.
func (m HouseHoldsManagerModel) CreateNewHouseholdMember(houseHoldMember *HouseHoldMember) error {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultHouseHoldManDBContextTimeout)
	defer cancel()
	// create new house hold member
	houseHoldMemberInfo, err := m.DB.CreateNewHouseholdMember(ctx, database.CreateNewHouseholdMemberParams{
		HouseholdID: houseHoldMember.HouseHoldID,
		Name:        houseHoldMember.Name,
		Age:         houseHoldMember.Age,
		Relation:    houseHoldMember.Relation,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "household_members" violates foreign key constraint "household_members_household_id_fkey"`:
			return ErrHouseHoldDoesNotExist
		case err.Error() == `pq: duplicate key value violates unique constraint "unique_household_member"`:
			return ErrHouseHoldMemberExists
		default:
			return err
		}
	}
	// set the new house hold member info
	houseHoldMember.ID = houseHoldMemberInfo.ID
	houseHoldMember.CreatedAt = houseHoldMemberInfo.CreatedAt
	houseHoldMember.UpdatedAt = houseHoldMemberInfo.UpdatedAt
	// return nil if everything is successful
	return nil
}
