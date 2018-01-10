package main

import (
	"net/http"
	"encoding/json"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"time"
	"os"
)

var jwtSecret = []byte("iiGImCNggvt7zg7hXAaAO8cRL1rQDI2D")
var tokenType = "jwt"

func main() {
	// instantiate the gorilla/mux router
	r := mux.NewRouter()

	// the default route handler is just the hello handler
	r.Handle("/", HelloHandler).Methods("GET")

	// the route for obtaining the token
	r.Handle("/token", CreateTokenHandler).Methods("GET")

	// additionally, there is a dynamic "name" route, protected by the JWT middleware
	r.Handle("/{name}", jwtMiddleware.Handler(HelloHandler)).Methods("GET")

	http.ListenAndServe(":3300", handlers.LoggingHandler(os.Stdout, r))
}

var HelloHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	var messageText string

	if len(name) > 0 {
		messageText = "Hello, " + name
	} else {
		messageText = "Hello stranger!"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{} {
		"message": messageText,
	})
})

var CreateTokenHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// create map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["name"] = "Michael"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// sign the token with the secret
	tokenString, _ := token.SignedString(jwtSecret)

	// prepare response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{} {
		"token": tokenString,
		"type": tokenType,
	})
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})