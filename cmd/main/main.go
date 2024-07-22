package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/thallosaurus/robocall/pkg/api"
	"github.com/thallosaurus/robocall/pkg/conf"
	"github.com/thallosaurus/robocall/pkg/svcctl"
)

func serveLogin(w http.ResponseWriter, _ *http.Request, config *conf.Config) {
	tmpl := template.Must(template.ParseFiles("/opt/robocall/web/tmpl/login.html"))

	err := tmpl.Execute(w, config)

	if err != nil {
		log.Fatal(err)
	}
}

func serveHome(w http.ResponseWriter, _ *http.Request, config *conf.Config) {

	tmpl := template.Must(template.ParseFiles("/opt/robocall/web/tmpl/index.html"))

	err := tmpl.Execute(w, config)

	if err != nil {
		log.Fatal(err)
	}
}

func SigHandler(sigchnl chan os.Signal) {
	for {
		signal := <-sigchnl
		switch signal {
		case syscall.SIGINT:
			fmt.Println("Got CTRL+C signal")
			fmt.Println("Stopping Asterisk...")
			//os.Exit(0)
			svcctl.StopService()
		default:
			fmt.Println("Ignoring signal: ", signal)
		}
	}
}

func main() {
	// Setup SIGINT
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go SigHandler(sigchnl)

	// Setup WaitGroup
	wg := &sync.WaitGroup{}

	api.SessionInit()

	// Initially apply config
	cnf := conf.FromDefaultFile()
	cnf.ApplyConfig()

	c, err := svcctl.RunService(wg)

	if err == nil {
		srv := httpServer(&cnf, wg)

		// Asterisk Child has exited, exit the Parent
		go func() {
			if s := <-c; s == 1 {
				//os.Exit(0)
				fmt.Println("Shutting down HTTP Server")
				srv.Shutdown(context.TODO())

				if err := srv.Shutdown(context.TODO()); err != nil {
					panic(err) // failure/timeout shutting down the server gracefully
				}
			}
		}()

		wg.Wait()
	} else {
		log.Fatal(err)
	}
}

func httpServer(cnf *conf.Config, wg *sync.WaitGroup) *http.Server {
	wg.Add(1)
	root := mux.NewRouter()
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: root,
	}

	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r, cnf)
	}).Methods("GET")

	root.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		api.Login(w, r)
	}).Methods("POST")

	root.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		serveLogin(w, r, cnf)
	}).Methods("GET")

	root.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		api.Logout(w, r)
	}).Methods("GET")

	api.Router(root.PathPrefix("/api").Subrouter(), cnf)
	//root.Handle("/", http.FileServer(http.Dir("/opt/robocall/web/client/")))
	fmt.Println("Listening on port 8080")

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv

}
