-- name: CreateNewGeoLocation :one
INSERT INTO geolocations (county, sub_county, location, sub_location)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

