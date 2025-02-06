package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

type AuthManagerModel struct {
	DB *database.Queries
}

type Apikey struct {
	Plaintext string
	Hash      []byte
}

const (
	DefaultAuthManDBContextTimeout = 5 * time.Second
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        int32     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ApiKey    Apikey    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Declare a new AnonymousUser variable.
var AnonymousUser = &User{}

// Check if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Check that the plaintext token has been provided and is exactly 26 bytes long.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	//v.Check(len(tokenPlaintext) == 36, "token", "must be valid")
}

// GetForApiKey() is a method that returns a user for a given api key
// We recieve a plaintext API key as a string and return a pointer to a User struct and an error
func (m AuthManagerModel) GetForApiKey(apiKey string) (*User, error) {
	// create context
	ctx, cancel := contextGenerator(context.Background(), DefaultAuthManDBContextTimeout)
	defer cancel()
	// calculate sha256 hash of the api key
	hash := sha256.Sum256([]byte(apiKey))

	// get user
	user, err := m.DB.GetForApiKey(ctx, hash[:])
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	// make a user
	authenticatedUser := &User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		ApiKey:    Apikey{Plaintext: apiKey, Hash: hash[:]},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return authenticatedUser, nil
}

func (m AuthManagerModel) GenerateToken() (*Apikey, error) {
	apiKey := &Apikey{}
	// Initialize a zero-valued byte slice with a length of 16 bytes.
	randomBytes := make([]byte, 16)
	// Use the Read() function from the crypto/rand package to fill the byte slice random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field.
	apiKey.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table.
	hash := sha256.Sum256([]byte(apiKey.Plaintext))
	apiKey.Hash = hash[:]
	return apiKey, nil
}
