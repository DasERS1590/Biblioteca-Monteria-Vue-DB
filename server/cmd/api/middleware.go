package main 

import (
    "net/http"
    "strings"
    "context"
    "biblioteca-api/internal/auth"
)


func(app *application) authenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
    
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "No autorizado", http.StatusUnauthorized)
                return
        }
        
        bearerToken := strings.Split(authHeader, "Bearer ")
        if len(bearerToken) != 2 {
            http.Error(w, "Token mal formado", http.StatusUnauthorized)
            return
        }
        
        claims, err := auth.ValidateToken(bearerToken[1])
        if err != nil {
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    } 
}
