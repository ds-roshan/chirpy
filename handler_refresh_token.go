package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ds-roshan/chirpy/internal/auth"
	"github.com/ds-roshan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}

	dbRefToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found", err)
		return
	}

	if dbRefToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	if dbRefToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	token, err := auth.MakeJWT(dbRefToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token!", err)
		return
	}

	type response struct {
		Token string `json:"token"`
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

	dbRefToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unauthorized access", err)
		return
	}

	if dbRefToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}

	err = cfg.db.UpdateRefreshTokenRevokeAt(r.Context(), database.UpdateRefreshTokenRevokeAtParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Token: refreshToken,
	})

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Cannot revoke", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
