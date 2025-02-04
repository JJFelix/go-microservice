package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(os.Getenv("SECREY_KEY"))

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "This is a simple Go Microservice")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Header["Token"] != nil{
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token)(interface{}, error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
					return nil, fmt.Errorf("invalid signing method")
				}

				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience{
					return nil, fmt.Errorf("invalid aud")
				}

				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss{
					return nil, fmt.Errorf("invalid iss")
				}

				return mySigningKey, nil
			})
			if err != nil{
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if token.Valid{
				endpoint(w, r)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
		}else{
			// fmt.Fprintf(w, "No authorization token provided")
			http.Error(w, "No authorization token provided", http.StatusUnauthorized)
		}
	})
}

func handleRequests(){
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main(){
	fmt.Println("Server")
	handleRequests()
}