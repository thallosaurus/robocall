package api

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

type User struct {
	Id int
}

func SessionInit() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore(authKeyOne, encryptionKeyOne)

	store.Options = &sessions.Options{
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(User{})
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
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
	session, err := store.Get(r, "session")
	session.Options.MaxAge = -1
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(&session.Values)
		if session.Values["user"] != nil {
			log.Println("Logged in: ", r.RequestURI)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "not authorized", http.StatusUnauthorized)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
	})
}
