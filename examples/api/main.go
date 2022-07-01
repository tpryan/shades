package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tpryan/shades"
)

func main() {
	port := "8080"
	srv := &server{
		router: mux.NewRouter().StrictSlash(true),
	}
	srv.routes()

	log.Printf("Starting webserver on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}

type server struct {
	router *mux.Router
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.HandleFunc("/random/{color}", s.handleRandom())
	s.router.HandleFunc("/random", s.handleRandom())
	s.router.HandleFunc("/family", s.handleFamilyList())
	s.router.HandleFunc("/healthz", s.handleHealthz())
	s.router.HandleFunc("/", s.handleHealthz())
}

func (s *server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	}
}

func (s *server) handleFamilyList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		families := shades.List()

		b, err := json.Marshal(families)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func (s *server) handleRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		color := strings.ToUpper(mux.Vars(r)["color"])

		if color == "" {
			color = "ALL"
		}

		shade, err := shades.NewFamily(color)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not get color family: %s", err)
			return
		}
		result := shade.Random()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result)
	}
}
