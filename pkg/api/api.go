package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/thallosaurus/robocall/pkg/conf"
	"github.com/thallosaurus/robocall/pkg/svcctl"
	"github.com/thallosaurus/robocall/pkg/utils"
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

	r.HandleFunc("/upload-sample", func(w http.ResponseWriter, r *http.Request) {
		uploadSample(w, r, c)
	}).Methods("POST")

	r.HandleFunc("/create-extension", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("POST")

	/*r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		svcctl.StopService()
		w.Write([]byte("stopping"))
	}).Methods("GET")*/

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

func uploadSample(w http.ResponseWriter, r *http.Request, c *conf.Config) {

	fmt.Println("Start parsing")
	err := r.ParseMultipartForm(5e7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	fmt.Println("Convert 2 GSM")

	mpart, mheader, err := r.FormFile("sample")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	var tmp bytes.Buffer
	io.Copy(&tmp, mpart)

	data, err := utils.ConvertToGSM(mheader.Filename, tmp.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		return
	}

	target_name := r.FormValue("name")

	fmt.Println("Uploaded Data converted ", data)

	dest, err := os.Create(fmt.Sprintf("/var/lib/asterisk/sounds/en/%s.gsm", target_name))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		os.Remove(data.Name())
		return
	}

	io.Copy(dest, data)

	c.Samples = append(c.Samples, conf.SampleEntry{
		SoundName: target_name,
		File:      dest,
	})
	c.ToDefaultFile()

	http.Redirect(w, r, "/", http.StatusFound)
}
