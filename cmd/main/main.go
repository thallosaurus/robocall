package main

import (
	"fmt"
	"log"
	"net/http"
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

	cnf := conf.FromDefaultFile()

	// Initially apply config
	cnf.ApplyConfig()

	svcctl.RunService()
	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r, &cnf)
	}).Methods("GET")

	api.Router(root.PathPrefix("/api").Subrouter(), &cnf)
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", root))
}
