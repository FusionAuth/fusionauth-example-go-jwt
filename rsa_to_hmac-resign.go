package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func main() {

	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := privateKey.PublicKey

	// generate token, presumably in an IdP or elsewhere
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["aud"] = "238d4793-70de-4183-9707-48ed8ecd19d9"
	claims["sub"] = "19016b73-3ffa-4b26-80d8-aa9287738677"
	claims["iss"] = "fusionauth.io"
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	claims["name"] = "Dan Moore"
	var roles [1]string
	roles[0] = "RETRIEVE_TODOS"
	claims["roles"] = roles

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		fmt.Println(fmt.Errorf("Something Went Wrong: %s", err.Error()))
	}

	fmt.Println(string(tokenString))

	// receive token at the API gateway
	decodedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &publicKey, nil
	})

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
	}

	// for real value, pull from secrets
	var mySigningKey = []byte("hello gophers!!!")

	hmacToken := jwt.New(jwt.SigningMethodHS256)
	hmacClaims := hmacToken.Claims.(jwt.MapClaims)
	for key, element := range decodedToken.Claims.(jwt.MapClaims) {
		hmacClaims[key] = element
	}

	// modify claims here if needed
	hmacClaims["exp"] = time.Now().Add(time.Second * 5).Unix()

	// sign JWT
	hmacTokenString, err := hmacToken.SignedString(mySigningKey)

	if err != nil {
		fmt.Println(fmt.Errorf("Something Went Wrong: %s", err.Error()))
	}

	fmt.Println(string(hmacTokenString))

}
