package main

import (
	"context"
	"net/http"

	"github.com/Blue-Davinci/SocialAid/internal/data"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
	"go.uber.org/zap"
)

type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	//app.log.PrintInfo("set user in request context", map[string]string{"name": user.Name, "email": user.Email})
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

// we will have a middleware to check if they provided the correct API key
// authenticate() is a middleware that checks if the user provided the correct API key
// We read the API key from the request, get the user for the API key, if the user exists we continue
// else we return an unauthorized response
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("authenticating user")
		// add vary auth header to the response to indicate that the response will vary based on the authentication
		w.Header().Add("Vary", "X-API-Key")
		// get the API key from the request
		authorizationHeader := r.Header.Get("ApiKey")
		app.logger.Info("authenticating user", zap.String("Authorization Header", authorizationHeader))
		if authorizationHeader == "" {
			app.logger.Info("no authorization header found")
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		// ectract the api key
		apiKey := authorizationHeader
		// validate the api key
		v := validator.New()
		if data.ValidateTokenPlaintext(v, apiKey); !v.Valid() {
			app.logger.Info("invalid api key", zap.Any("error", v.Errors))
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		// get the user for the API key
		user, err := app.models.Auth.GetForApiKey(apiKey)
		if err != nil {
			app.logger.Info("authenticating use failed", zap.Error(err))
			switch {
			case err == data.ErrUserNotFound:
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)

			}
			return
		}
		// log
		app.logger.Info("user", zap.String("email", user.Email))
		// set context user to add user info
		r = app.contextSetUser(r, user)
		// we are good, we can continue
		next.ServeHTTP(w, r)
	})
}

// Create a new requireAuthenticatedUser() middleware to check that a user is not
// anonymous.
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the contextGetUser() helper to retrieve the user
		// information from the request context.
		user := app.contextGetUser(r)
		// If the user is anonymous, then call the authenticationRequiredResponse() to
		// inform the client that they should authenticate before trying again.
		if user.IsAnonymous() {
			app.logger.Info("user", zap.Any("user is anonymous", user))
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) testHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"` //simulate email
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &data.User{
		Email: input.Email,
	}
	app.logger.Info("user", zap.String("Token", user.Email))
	// generate a new api key
	apiKey, err := app.models.Auth.GenerateToken()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.ApiKey = *apiKey
	app.logger.Info("user", zap.Any("Your Api Key", user.ApiKey))
	// output the hashed api key and user
	err = app.writeJSON(w, http.StatusOK, envelope{"Your Api Key": user.ApiKey}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
