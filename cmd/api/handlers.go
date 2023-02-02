package main

import (
	"errors"
	"github.com/aidanb22/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := app.DB.AllPlayers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, players)
}
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.AllUsers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) PlayerDatabase(w http.ResponseWriter, r *http.Request) {
	players, err := app.DB.AllPlayers()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, players)
}

func (app *application) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //gets id for the player
	playerID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	player, err := app.DB.OnePlayer(playerID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, player)

}

func (app *application) PlayerForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //gets id for the player
	playerID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	player, err := app.DB.OnePlayerForEdit(playerID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var payload = struct {
		Player *models.Player `json:"player"`
	}{
		player, //populating the struct here
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) InsertPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := app.readJSON(w, r, &player)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//try to get an image

	player.CreatedAt = time.Now()
	player.UpdatedAt = time.Now()
	_, err = app.DB.InsertPlayer(player)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := JSONResponse{
		Error:   false,
		Message: "Player updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)

}

func (app *application) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	var payload models.Player

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
	}

	player, err := app.DB.OnePlayer(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
	}

	player.Plyname = payload.Plyname
	player.College = payload.College
	player.Age = payload.Age
	player.Height = payload.Height
	player.Description = payload.Description
	player.UpdatedAt = time.Now()

	err = app.DB.UpdatePlayer(*player)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "player updated",
	}
	app.writeJSON(w, http.StatusAccepted, resp)

}

func (app *application) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
	}

	err = app.DB.DeletePlayer(id)
	if err != nil {
		app.errorJSON(w, err)
	}
	resp := JSONResponse{
		Error:   false,
		Message: "player deleted",
	}
	app.writeJSON(w, http.StatusAccepted, resp)

}
