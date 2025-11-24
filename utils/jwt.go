package utils

import (
    "errors"
    "os"
    "prestasi_backend/app/model"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Membuat Token
func GenerateToken(claim model.JWTClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    return token.SignedString(jwtKey)
}

// Parsing Token
func ParseToken(tokenString string) (*model.JWTClaims, error) {

    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }

    parsed, err := jwt.ParseWithClaims(
        tokenString,
        &model.JWTClaims{},
        func(t *jwt.Token) (interface{}, error) { return jwtKey, nil },
    )

    if err != nil {
        return nil, errors.New("token tidak valid")
    }

    claims, ok := parsed.Claims.(*model.JWTClaims)
    if !ok || !parsed.Valid {
        return nil, errors.New("token tidak valid")
    }

    return claims, nil
}
