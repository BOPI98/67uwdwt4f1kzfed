package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Request-Method", "GET,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-requested-with")
			next.ServeHTTP(w, r)
		})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					log.Print(string(debug.Stack()))
				}
			}()
			next.ServeHTTP(w, r)
		})
}
