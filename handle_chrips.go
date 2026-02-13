package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/deexth/chirpy/internal/auth"
	"github.com/deexth/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Chirp
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})

	// badWords := map[string]struct{}{
	// 	"kerfuffle": {},
	// 	"sharbert":  {},
	// 	"fornax":    {},
	// }
	// cleaned := getCleanedBody(params.Body, badWords)
	//
	// respondWithJSON(w, http.StatusOK, returnVals{
	// 	CleanedBody: cleaned,
	// })

}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context(), 5)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "issue retrieving chirps", err)
		return
	}

	newChirps := make([]Chirp, 0, len(chirps))

	for _, chirp := range chirps {
		newChirps = append(newChirps, Chirp(chirp))
	}

	respondWithJSON(w, http.StatusOK, newChirps)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusForbidden, "unauthorized", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "issue retrieving chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp(chirp))
}

func (cfg *apiConfig) handleChirpDeletion(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "unauthorized", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusForbidden, "unauthorized", err)
		return
	}

	numAffectedRows, err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userID,
		ID:     chirpID,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "chirp not found", err)
		return
	}

	if numAffectedRows == 0 {
		respondWithError(w, http.StatusForbidden, "chirp not found", errors.New("no rows affected"))
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
