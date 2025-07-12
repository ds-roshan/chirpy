package main

import (
	"net/http"
	"time"

	"github.com/ds-roshan/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldnot get user for refresh token", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token!", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Reresh token not found", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
