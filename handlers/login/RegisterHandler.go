package login

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mandoju/BindAPI/models"
	"github.com/mandoju/BindAPI/utils/Database"
	"net/http"
	"time"
)

// RegisterHandlerInput  is the structure oof the JSON input of loginHandler
type RegisterHandlerInput struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterHandlerOutput is the structure oof the JSON output of loginHandler
type RegisterHandlerOutput struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var input RegisterHandlerInput
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err.Error())
	}

	stmt, err := Database.Db.Prepare("SELECT username from users where username = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err.Error())
	}
	stmtt, err := Database.Db.Prepare("Insert INTO users(username,password,email) VALUES(?,?,?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err.Error())
	}

	usernames, err := stmt.Query(input.Username)
	if err != nil {
		fmt.Println(err.Error())

		panic(err.Error())
	}
	if usernames.Next() {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err = stmtt.Exec(input.Username, input.Password, input.Email)
	if err != nil {
		fmt.Println(err.Error())

		panic(err.Error())
	}
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: input.Username,
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
		fmt.Println(err.Error())

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
	output := RegisterHandlerOutput{
		Username: input.Username,
		Token:    tokenString}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
