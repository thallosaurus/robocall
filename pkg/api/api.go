package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thallosaurus/robocall/pkg/conf"
	"github.com/thallosaurus/robocall/pkg/svcctl"
)

func Router(r *mux.Router, c *conf.Config) {
	r.Use(loginMiddleware)
	r.HandleFunc("/user-config", func(w http.ResponseWriter, r *http.Request) {
		userConfig(w, r, c)
	}).Methods("POST")

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(b)
	}).Methods("GET")

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		svcctl.StopService()
		w.Write([]byte("stopping"))
	}).Methods("GET")

}

func userConfig(w http.ResponseWriter, r *http.Request, c *conf.Config) {
	r.ParseForm()

	c.Sip.Name = "main"
	c.Sip.Host = r.FormValue("gateway-host")
	c.Sip.Username = r.FormValue("username")
	c.Sip.Password = r.FormValue("password")

	fmt.Println(c)
	c.ApplyConfig()
	svcctl.ReloadSIPModule()

	http.Redirect(w, r, "/", http.StatusFound)
}
