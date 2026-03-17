package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deexth/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "something went wrong", err)
		return
	}

	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "something went wrong", errors.New("apikey mismatch"))
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	numAffectedRows, err := cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "something went wrong", err)
		return
	}

	if numAffectedRows == 0 {
		respondWithError(w, http.StatusNotFound, "something went wrong", errors.New("no rows affected"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
