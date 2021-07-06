package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func main() {

	var mySigningKey = []byte("hello gophers!!!")

	// user api
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["aud"] = "238d4793-70de-4183-9707-48ed8ecd19d9"
	claims["sub"] = "19016b73-3ffa-4b26-80d8-aa9287738677"
	claims["iss"] = "wrong.io"
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	claims["name"] = "Dan Moore"
	var roles [1]string
	roles[0] = "RETRIEVE_TODOS"
	claims["roles"] = roles

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Println(fmt.Errorf("Something Went Wrong: %s", err.Error()))
	}

	fmt.Println(string(tokenString))

	// todo api
	decodedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     		}
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return
	}

	// verify aud claim
	expectedAud := "238d4793-70de-4183-9707-48ed8ecd19d9"
	checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(expectedAud, false)
	if !checkAudience {
		fmt.Println("invalid aud")
		return
	}
	// verify iss claim
	expectedIss := "fusionauth.io"
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(expectedIss, false)
	if !checkIss {
		fmt.Println("invalid iss")
		return
	}

	fmt.Println("")
	fmt.Printf("%+v\n", decodedToken.Claims)
}
