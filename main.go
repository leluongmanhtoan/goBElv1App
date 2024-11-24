package main

import (
	"net/http"
)

type UserProfile struct{
	Fullname string
	Age int
	Email string
}

type Account struct{
	Profile UserProfile
	Username string
	Password string
}


func check(e error){
	if e != nil {
		panic(e)
	}
}

func signup(w http.ResponseWriter, req *http.Request) {

}

func signin(w http.ResponseWriter, req *http.Request) {

}

func main() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signin", signin)
	http.ListenAndServe(":8090", nil)
}
