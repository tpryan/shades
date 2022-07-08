package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tpryan/shades"
)

var (
	errorNoGrayscale = fmt.Errorf("cannot handle grayscale colors")
	errorNoColor     = fmt.Errorf("color cannot be blank")
	errorInValid     = fmt.Errorf("a valid color (#xxxxxx format) must be input to find the family for it")
)

func main() {
	port := os.Getenv("PORT")
	if port != "" {
		port = "8080"
	}

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
	s.router.HandleFunc("/invert", s.handleInvert()).Methods(http.MethodPost)
	s.router.HandleFunc("/family/find", s.handleFamilyFind()).Methods(http.MethodPost)
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

func (s *server) handleFamilyFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		color := strings.ToUpper(r.FormValue("color"))

		if color == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, errorNoColor.Error())
			return
		}

		result := shades.FindFamily(color)
		if result == "" {
			if shades.IsGreyScale(color) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, errorNoGrayscale.Error())
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, errorInValid.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result)
	}
}

func (s *server) handleInvert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		color := strings.ToUpper(r.FormValue("color"))

		if color == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, errorNoColor.Error())
			return
		}

		result := shades.Invert(color)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, result)
	}
}
