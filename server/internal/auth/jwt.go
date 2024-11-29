package auth

import (
    "time"
    "errors"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("biblioteca_secret_key_2024")

type Claims struct {
    UserID int    `json:"user_id"`
    Username string `json:"username"`
    Role    string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("token inválido")
    }

    return claims, nil
}
