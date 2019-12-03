// Copyright 2019 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package main is a Kubernetes API proxy. It exposes a smaller surface of the
// API and limits operations to specifically selected labels, and deployments

package shades

import (
	"math/rand"
	"testing"
)

func TestRandomInRange(t *testing.T) {
	cases := []struct {
		in   Range
		seed int64
		want float64
	}{
		{Range{0, 1}, 1, 0.6046602879796196},
		{Range{0, 100}, 1, 60.466028797961954},
		{Range{50, 100}, 1, 80.23301439898097},
		{Range{.2, .8}, 1, 0.5627961727877718},
	}

	for _, c := range cases {
		rand.Seed(c.seed)
		got := rando(c.in)
		if got != c.want {
			t.Errorf("randomInRange(%v) got %f, want %f", c.in, got, c.want)
		}
	}

}

func TestBetween(t *testing.T) {
	cases := []struct {
		r    Range
		in   float64
		want bool
	}{
		{Range{-50, 10}, 5, true},
		{Range{0, 100}, 75, true},
		{Range{50, 100}, 5, false},
		{Range{.2, .8}, .5, true},
		{Range{.2, .8}, 5, false},
	}

	for _, c := range cases {
		got := c.r.Between(c.in)
		if got != c.want {
			t.Errorf("Between(%f) got %t, want %t", c.in, got, c.want)
		}
	}

}

func TestNewFamily(t *testing.T) {
	cases := []struct {
		in   string
		want Family
		err  error
	}{
		{"BLUE", Family{"Blue", "0000FF", Range{221, 240}, Range{.1, 1}, Range{.2, 1}}, nil},
		{"BLUEGREEN", Family{}, ErrNotValidFamily},
	}

	for _, c := range cases {
		got, err := NewFamily(c.in)
		if got != c.want {
			t.Errorf("NewFamily(%s) got %v, want %v", c.in, got, c.want)
		}
		if err != c.err {
			t.Errorf("NewFamily(%s) got err '%v', want err '%v'", c.in, err, c.err)
		}

	}

}

func TestRandom(t *testing.T) {

	blue := Family{"Blue", "0000FF", Range{221, 240}, Range{.1, 1}, Range{.2, 1}}
	red := Family{"Red", "FF0000", Range{-10, 20}, Range{.2, 1}, Range{.2, 1}}
	green := Family{"Green", "00FF00", Range{81, 140}, Range{.4, 1}, Range{.3, .8}}

	cases := []struct {
		in   Family
		seed int64
		want string
	}{

		{red, 1, "#fc8b79"},
		{red, 2, "#572428"},
		{green, 1, "#51fc47"},
		{green, 2, "#528225"},
		{blue, 1, "#7a8afb"},
		{blue, 2, "#293452"},
	}

	for _, c := range cases {
		rand.Seed(c.seed)
		got := c.in.Random(c.seed)
		if got != c.want {
			t.Errorf("%s.Random(%d) got %s, want %s", c.in.Name, c.seed, got, c.want)
		}

	}

}

func TestFamilyIn(t *testing.T) {
	cases := []struct {
		family Family
		in     string
		want   bool
	}{
		{list["RED"], "#e58677", true},
		{list["RED"], "#72393d", true},
		{list["RED"], "#00FF00", false},
		{list["GREEN"], "#55e74c", true},
		{list["GREEN"], "#609331", true},
		{list["GREEN"], "#FF0000", false},
		{list["BLUE"], "#7a87e2", true},
		{list["BLUE"], "#404b6a", true},
		{list["BLUE"], "#00FF00", false},
	}

	for _, c := range cases {
		got := c.family.In(c.in)
		if got != c.want {
			t.Errorf("%s.In(%s) got %t, want %t", c.family.Name, c.in, got, c.want)
		}

	}

}

func TestFindFamily(t *testing.T) {
	cases := []struct {
		want string
		in   string
	}{
		{"RED", "#e58677"},
		{"RED", "#72393d"},
		{"GREEN", "#00FF00"},
		{"GREEN", "#55e74c"},
		{"GREEN", "#609331"},
		{"RED", "#FF0000"},
		{"BLUE", "#7a87e2"},
		{"BLUE", "#404b6a"},
		{"GREEN", "#00FF00"},
		{"", ""},
		{"RED", "#F00"},
		{"GREEN", "#0F0"},
		{"BLUE", "#00F"},
	}

	for _, c := range cases {
		got := FindFamily(c.in)
		if got != c.want {
			t.Errorf("FindFamily(%s) got %s, want %s", c.in, got, c.want)
		}

	}

}

func TestInvert(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"#FF0000", "#00FFFF"},
		{"#00FF00", "#FF00FF"},
		{"#FFFFFF", "#000000"},
		{"#000000", "#FFFFFF"},
		{"#DADADA", "#252525"},
		{"#19547A", "#E6AB85"},
	}

	for _, c := range cases {
		got := Invert(c.in)
		if got != c.want {
			t.Errorf("Invert(%s) got %s, want %s", c.in, got, c.want)
		}

	}

}

func TestList(t *testing.T) {
	l := List()

	if len(list) != len(l) {
		t.Errorf("len(List()) got %d, want %d", len(l), len(list))
	}

}
