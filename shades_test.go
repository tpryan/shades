// Copyright 2019 Google Inc. All Rights Reserved.
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

package shades

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInRange(t *testing.T) {
	tests := map[string]struct {
		in   Range
		seed int64
		want float64
	}{
		"0,1":    {in: Range{0, 1}, seed: 1, want: 0.6046602879796196},
		"0,100":  {in: Range{0, 100}, seed: 1, want: 60.466028797961954},
		"50,100": {in: Range{50, 100}, seed: 1, want: 80.23301439898098},
		".2,.8":  {in: Range{.2, .8}, seed: 1, want: 0.5627961727877718},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rand.Seed(tc.seed)
			got := rando(tc.in)

			assert.InDelta(t, tc.want, got, 0.0000001)

		})
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
	tests := map[string]struct {
		in   Color
		want Family
	}{
		"Red": {
			in: Red,
			want: Family{
				Name: "Red",
				Base: "FF0000",
				Hue:  Range{-10, 20},
				Sat:  Range{.2, 1},
				Lum:  Range{.2, 1},
			},
		},

		"Orange": {
			in: Orange,
			want: Family{
				Name: "Orange",
				Base: "FFA500",
				Hue:  Range{21, 50},
				Sat:  Range{.3, 1},
				Lum:  Range{.4, 1},
			},
		},

		"Yellow": {
			in: Yellow,
			want: Family{
				Name: "Yellow",
				Base: "FFFF00",
				Hue:  Range{51, 60},
				Sat:  Range{.4, 1},
				Lum:  Range{.63, 1},
			},
		},

		"Green": {
			in: Green,
			want: Family{
				Name: "Green",
				Base: "00FF00",
				Hue:  Range{81, 140},
				Sat:  Range{.4, 1},
				Lum:  Range{.3, .8},
			},
		},

		"Cyan": {
			in: Cyan,
			want: Family{
				Name: "Cyan",
				Base: "00FFFF",
				Hue:  Range{170, 200},
				Sat:  Range{.25, 1},
				Lum:  Range{.3, 1},
			},
		},

		"Blue": {
			in: Blue,
			want: Family{
				Name: "Blue",
				Base: "0000FF",
				Hue:  Range{221, 240},
				Sat:  Range{.1, 1},
				Lum:  Range{.2, 1},
			},
		},
		"Purple": {
			in: Purple,
			want: Family{
				Name: "Purple",
				Base: "800080",
				Hue:  Range{241, 280},
				Sat:  Range{.3, 1},
				Lum:  Range{.4, .7},
			},
		},
		"Magenta": {
			in: Magenta,
			want: Family{
				Name: "Magenta",
				Base: "FF00FF",
				Hue:  Range{281, 320},
				Sat:  Range{.35, 1},
				Lum:  Range{.3, .7},
			},
		},
		"All": {
			in: All,
			want: Family{
				Name: "All",
				Base: "FF00FF",
				Hue:  Range{0, 360},
				Sat:  Range{0, 1},
				Lum:  Range{0, 1},
			},
		},
		"bad value": {
			in:   0,
			want: NewFamily(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewFamily(tc.in)
			assert.Equal(t, tc.want, got)
		})
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
		got := c.in.Random()
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
		{"notacolor", "#53"},
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

func IsGreyScale(hex string) bool {
	in := strings.ReplaceAll(hex, "#", "")

	digits := strings.Split(in, "")

	if len(digits) == 6 {
		tmp := []string{}
		tmp = append(tmp, fmt.Sprintf("%s%s", digits[0], digits[1]))
		tmp = append(tmp, fmt.Sprintf("%s%s", digits[2], digits[3]))
		tmp = append(tmp, fmt.Sprintf("%s%s", digits[4], digits[5]))

		digits = tmp
	}

	if len(digits) == 3 {
		if digits[0] == digits[1] && digits[0] == digits[2] {
			return true
		}
		return false
	}

	return false
}

func TestIsNumeric(t *testing.T) {
	tests := map[string]struct {
		in   string
		want bool
	}{
		"123": {in: "123", want: true},
		"97j": {in: "97j", want: false},
		"1":   {in: "1", want: true},
		"ads": {in: "ads", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := isNumeric(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}

func TestIsHexColor(t *testing.T) {
	tests := map[string]struct {
		in   string
		want bool
	}{
		"#123456": {in: "#123456", want: true},
		"#ccc":    {in: "#ccc", want: true},
		"1AFFa1":  {in: "#1AFFa1", want: true},
		"F00":     {in: "#F00", want: true},
		"123456":  {in: "123456", want: false},
		"123abce": {in: "#123abce", want: false},
		"afafah":  {in: "#afafah", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := isHexColor(tc.in)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
