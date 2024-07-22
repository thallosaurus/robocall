package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/thallosaurus/robocall/pkg/api"
	"github.com/thallosaurus/robocall/pkg/conf"
	"github.com/thallosaurus/robocall/pkg/svcctl"
)

func serveHome(w http.ResponseWriter, _ *http.Request, config *conf.Config) {

	tmpl := template.Must(template.ParseFiles("/opt/robocall/web/tmpl/index.html"))

	err := tmpl.Execute(w, config)

	if err != nil {
		log.Fatal(err)
	}
}

var root = mux.NewRouter()

func main() {

	/*cnf := conf.Config{
		Name:     "test",
		Username: "test",
		Password: "test",
		Host:     "test",
	}*/
	api.SessionInit()

	cnf := conf.FromDefaultFile()

	// Initially apply config
	cnf.ApplyConfig()

	c, err := svcctl.RunService()

	if err == nil {

		// Asterisk Child has exited, exit the Parent
		go func() {
			if s := <-c; s == 1 {
				os.Exit(0)
			}
		}()

		root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			serveHome(w, r, &cnf)
		}).Methods("GET")

		root.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
			api.Login(w, r)
		})

		root.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
			api.Logout(w, r)
		}).Methods("GET")

		api.Router(root.PathPrefix("/api").Subrouter(), &cnf)
		//root.Handle("/", http.FileServer(http.Dir("/opt/robocall/web/client/")))
		fmt.Println("Listening on port 8080")
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", root))
	} else {
		log.Fatal(err)
	}
}
