package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deexth/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string        `json:"password"`
		Email     string        `json:"email"`
		ExpiresIn time.Duration `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON body", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if params.ExpiresIn == 0 {
		params.ExpiresIn = 3600
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, params.ExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue creating token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: accessToken,
	})
}
