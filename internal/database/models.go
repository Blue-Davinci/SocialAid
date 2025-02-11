// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"
)

type Geolocation struct {
	ID          int32
	County      string
	SubCounty   string
	Location    string
	SubLocation string
	CreatedAt   time.Time
}

type Household struct {
	ID            int32
	ProgramID     int32
	GeolocationID int32
	Name          string
	CreatedAt     time.Time
}

type HouseholdHead struct {
	ID          int32
	HouseholdID int32
	Name        string
	NationalID  string
	PhoneNumber string
	Age         int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type HouseholdMember struct {
	ID          int32
	HouseholdID int32
	Name        string
	Age         int32
	Relation    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Program struct {
	ID          int32
	Name        string
	Category    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID        int32
	Email     string
	Name      string
	ApiKey    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
