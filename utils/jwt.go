package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mandoju/BindAPI/models"
	"net/http"
)

var jwtKey, _ = GetJwtKey()

func checkJwtToken(r *http.Request) (string,int) {
	// Check if cookie has token value
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			return "", http.StatusUnauthorized
		}
		// For any other type of error, return a bad request status
		return "", http.StatusBadRequest
	}

	// Getting  token Value
	tokenValue := c.Value

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", http.StatusUnauthorized
		}
		return "", http.StatusBadRequest
	}
	if !tkn.Valid {
		return "", http.StatusUnauthorized
	}

	return claims.Username, http.StatusOK
}