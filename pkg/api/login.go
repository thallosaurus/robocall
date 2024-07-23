package api

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

type User struct {
	Id int
}

func SessionInit() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	//encryptionKeyOne := securecookie.GenerateRandomKey(32)
	Store = sessions.NewCookieStore(authKeyOne)

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(User{})
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	session.Values["user"] = User{
		Id: 0,
	}

	log.Println(r.RequestURI)
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gather Session
		session, err := Store.Get(r, "session")
		if err != nil {
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(&session.Values)

		// Check if a Session exists, if yes continue with the request
		if session.Values["user"] != nil {
			log.Println("Logged in: ", r.RequestURI)
			next.ServeHTTP(w, r)
		} else {
			//http.Error(w, "not authorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}
