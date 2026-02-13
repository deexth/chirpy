package main

import (
	"encoding/json"
	"net/http"

	"github.com/deexth/chirpy/internal/auth"
	"github.com/deexth/chirpy/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unautorized", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unautorized", err)
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong", err)
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "someting went wrong", err)
		return
	}

	err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:       userID,
		Email:    params.Email,
		Password: hashedPwd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Email string `json:"email"`
	}{
		Email: params.Email,
	})
}
