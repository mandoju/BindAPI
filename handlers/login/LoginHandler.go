package login

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/mandoju/BindAPI/models"
	"github.com/mandoju/BindAPI/utils"
	"net/http"
	"time"
)

var jwtKey, _ = utils.GetJwtKey()

// Users HardCoded
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// LoginHandlerInput  is the structure oof the JSON input of loginHandler
type LoginHandlerInput struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginHandlerOutput is the structure oof the JSON output of loginHandler
type LoginHandlerOutput struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// LoginHandler is the handle that offers the user a jwt token given the right login/password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds LoginHandlerInput
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	output := LoginHandlerOutput{
		Username: creds.Username,
		Token:    tokenString}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
