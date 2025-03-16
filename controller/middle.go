package controller

import (
	"go-web/model"
	"log"
	"net/http"
)

func middleAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getSessionUser(r)
		log.Println("middle:", username)
		_, err2 := model.GetUserByUsername(username)
		if username != "" && err == nil && err2 == nil {
			log.Println("Last seen:", username)
			model.UpdateLastSeen(username)
		}
		if err != nil || err2 != nil {
			log.Println("middle get session err and redirect to login")
			// log.Println("err:", err, "err2:", err2)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	}
}
