package main

import (
	"fmt"
	"log"
	"net/http"
)

func Server() {
	SetAccessTokenHandler := http.HandlerFunc(SetAccessTokenHandler)
	http.Handle("/get_access", SetAccessTokenHandler)

	RefreshTokens := http.HandlerFunc(RefreshTokens)
	http.Handle("/refresh", RefreshTokens)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func SetAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenStringA, _ := TokenAccess()
	tokenStringR, _ := TokenRefresh()
	serDetails := make(map[string]string)
	serDetails["name1"] = "access"
	serDetails["var1"] = tokenStringA
	serDetails["name2"] = "refresh"
	serDetails["name2"] = tokenStringR
	if r.Method == "POST" {
		cookie := &http.Cookie{
			Name: "Access",

			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(200)

		fmt.Fprintf(w, "POST only")
	}
	return
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {

	return
}
