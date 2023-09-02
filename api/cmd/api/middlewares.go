package main

import (
	"context"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/auth"
	"github.com/LeonLow97/internal/utils"
)

func setAccessControlHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}

type ContextUserId int
const ContextUserIdKey ContextUserId = 0

func authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.ValidateToken(r)
		if err != nil {
			log.Println(err)

			utils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIdKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
