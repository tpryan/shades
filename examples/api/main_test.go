// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tpryan/shades"
)

func TestHealthzHandler(t *testing.T) {
	srv := &server{
		router: mux.NewRouter().StrictSlash(true),
	}
	srv.routes()

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handleHealthz())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `ok`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestMeta(t *testing.T) {
	tests := map[string]struct {
		input bool
		want  bool
	}{
		"1": {input: true, want: true},
		"2": {input: false, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestRandomHandler(t *testing.T) {
	tests := map[string]struct {
		color  string
		want   string
		status int
	}{
		"all":  {color: "", want: "#5995fa", status: http.StatusOK},
		"blue": {color: "blue", want: "#7a8afb", status: http.StatusOK},
		"yuck": {color: "yuck", want: fmt.Sprintf("could not get color family: %s", shades.ErrNotValidFamily), status: http.StatusInternalServerError},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rand.Seed(1)
			srv := &server{
				router: mux.NewRouter().StrictSlash(true),
			}
			srv.routes()
			req, err := http.NewRequest("GET", "/random/", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.color != "" {
				vars := map[string]string{
					"color": tc.color,
				}
				req = mux.SetURLVars(req, vars)

			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.handleRandom())
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			if rr.Body.String() != tc.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.want)
			}
		})
	}
}

func TestFamilyListHandler(t *testing.T) {
	srv := &server{
		router: mux.NewRouter().StrictSlash(true),
	}
	srv.routes()
	req, err := http.NewRequest("GET", "/family", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handleFamilyList())
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `["ALL","BLUE","CYAN","GREEN","MAGENTA","ORANGE","PURPLE","RED","YELLOW"]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFamilyFindHandler(t *testing.T) {
	tests := map[string]struct {
		want   string
		color  string
		status int
	}{
		"Red 1":         {"RED", "#e58677", http.StatusOK},
		"Red 2":         {"RED", "#72393d", http.StatusOK},
		"Green 1":       {"GREEN", "#00FF00", http.StatusOK},
		"Green 2":       {"GREEN", "#55e74c", http.StatusOK},
		"Green 3":       {"GREEN", "#609331", http.StatusOK},
		"Red 4":         {"RED", "#FF0000", http.StatusOK},
		"Blue 1":        {"BLUE", "#7a87e2", http.StatusOK},
		"Blue 2":        {"BLUE", "#404b6a", http.StatusOK},
		"Green 4":       {"GREEN", "#00FF00", http.StatusOK},
		"Blank":         {errorNoColor.Error(), "", http.StatusInternalServerError},
		"Red 5":         {"RED", "#F00", http.StatusOK},
		"Green 5":       {"GREEN", "#0F0", http.StatusOK},
		"Blue 3":        {"BLUE", "#00F", http.StatusOK},
		"Silver":        {errorNoGrayscale.Error(), "#ccc", http.StatusInternalServerError},
		"Black":         {errorNoGrayscale.Error(), "#000000", http.StatusInternalServerError},
		"InAppropriate": {errorInValid.Error(), "InAppropriate", http.StatusInternalServerError},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rand.Seed(1)
			srv := &server{
				router: mux.NewRouter().StrictSlash(true),
			}
			srv.routes()

			reader := strings.NewReader(fmt.Sprintf("color=%s", tc.color))

			req, err := http.NewRequest("POST", "/family/find", reader)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tc.color != "" {
				vars := map[string]string{
					"color": tc.color,
				}
				req = mux.SetURLVars(req, vars)

			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.handleFamilyFind())
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			if rr.Body.String() != tc.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.want)
			}
		})
	}
}

func TestInvertHandler(t *testing.T) {
	tests := map[string]struct {
		color  string
		want   string
		status int
	}{
		"FF0000": {"#FF0000", "#00FFFF", http.StatusOK},
		"00FF00": {"#00FF00", "#FF00FF", http.StatusOK},
		"FFFFFF": {"#FFFFFF", "#000000", http.StatusOK},
		"000000": {"#000000", "#FFFFFF", http.StatusOK},
		"DADADA": {"#DADADA", "#252525", http.StatusOK},
		"19547A": {"#19547A", "#E6AB85", http.StatusOK},
		"":       {"", errorNoColor.Error(), http.StatusInternalServerError},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv := &server{
				router: mux.NewRouter().StrictSlash(true),
			}
			srv.routes()

			reader := strings.NewReader(fmt.Sprintf("color=%s", tc.color))

			req, err := http.NewRequest("POST", "/invert", reader)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tc.color != "" {
				vars := map[string]string{
					"color": tc.color,
				}
				req = mux.SetURLVars(req, vars)

			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.handleInvert())
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			if rr.Body.String() != tc.want {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.want)
			}
		})
	}
}
