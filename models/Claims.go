package models

import "github.com/dgrijalva/jwt-go"

// Claims is the structure of given JWT token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
