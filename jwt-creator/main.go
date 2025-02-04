package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(os.Getenv("SECREY_KEY"))

func GetJWT()(string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "jjfelix"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute*2).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil{
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request){
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err != nil{
		fmt.Println("Failed to generate token")
	}
	fmt.Fprintf(w, "%s", string(validToken))
}

func handleRequests(){
	http.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main(){
	handleRequests()
}