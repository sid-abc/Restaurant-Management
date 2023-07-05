package middlewares

import (
	"context"
	"example/rms/database"
	"example/rms/database/dbHelper"
	"example/rms/models"
	"net/http"
)

func MiddlewareUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*models.Claims)
		userID := claims.UserID
		userRoles, err := dbHelper.GetUserRoles(database.Todo, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var isUser = false
		for _, x := range userRoles {
			if x == models.Role3 {
				isUser = true
				break
			}
		}
		if !isUser {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*models.Claims)
		userID := claims.UserID
		userRoles, err := dbHelper.GetUserRoles(database.Todo, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var isAdmin = false
		for _, x := range userRoles {
			if x == models.Role1 {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareSubAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*models.Claims)
		userID := claims.UserID
		userRoles, err := dbHelper.GetUserRoles(database.Todo, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var isSubAdmin = false
		for _, x := range userRoles {
			if x == models.Role2 {
				isSubAdmin = true
				break
			}
		}
		if !isSubAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
