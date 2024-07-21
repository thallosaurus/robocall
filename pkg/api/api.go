package api

import (
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
