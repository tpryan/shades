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

// Package shades wraps the go-colorful library and uses it to create random
// shades of a base color.  For example you input red as the color, and can
// return various random shades of red.
package shades

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var seed = time.Now().UnixNano()

func init() {
	rand.Seed(seed)
}

// ErrNotValidFamily indicates that you tried to request a color family that
// does not exist
var ErrNotValidFamily = fmt.Errorf("the input color family is not valid")

var list = map[string]Family{
	"RED":     {"Red", "FF0000", Range{-10, 20}, Range{.2, 1}, Range{.2, 1}},
	"ORANGE":  {"Orange", "FFA500", Range{21, 50}, Range{.3, 1}, Range{.4, 1}},
	"YELLOW":  {"Yellow", "FFFF00", Range{51, 60}, Range{.4, 1}, Range{.63, 1}},
	"GREEN":   {"Green", "00FF00", Range{81, 140}, Range{.4, 1}, Range{.3, .8}},
	"CYAN":    {"Cyan", "00FFFF", Range{170, 200}, Range{.25, 1}, Range{.3, 1}},
	"BLUE":    {"Blue", "0000FF", Range{221, 240}, Range{.1, 1}, Range{.2, 1}},
	"PURPLE":  {"Purple", "800080", Range{241, 280}, Range{.3, 1}, Range{.4, .7}},
	"MAGENTA": {"Magenta", "FF00FF", Range{281, 320}, Range{.35, 1}, Range{.3, .7}},
	"ALL":     {"All", "FF00FF", Range{0, 360}, Range{0, 1}, Range{0, 1}},
}

// Range is a upper and lower bound for a pair of integers for use in the
// go-colorful library
type Range struct {
	Bottom float64
	Top    float64
}

// Between determines if a given number is contained in a range.
func (r *Range) Between(value float64) bool {
	if r.Bottom < 0 {
		diff := 360 + r.Bottom

		if value >= diff {
			value = -(360 - value)
		}

	}
	if value >= r.Bottom && value <= r.Top {
		return true
	}
	return false
}

// Family is a set of definitions of shades of a color.  It has ranges for Hue,
// Saturation, and Luminosity.  These ranges define a set of colors that can be
// considered to be shades of the base color.  This allows us to generate random
// color shades based on that base color.
type Family struct {
	Name string
	Base string
	Hue  Range
	Sat  Range
	Lum  Range
}

// In determines if a given hexidecimal color is withing a given color family.
// If the given hex string is invalid, this function returns false.
func (f *Family) In(hex string) bool {
	color, err := colorful.Hex(hex)
	if err != nil {
		return false
	}

	h, s, l := color.Hsl()

	if f.Hue.Between(h) && f.Sat.Between(s) && f.Lum.Between(l) {
		return true
	}
	return false

}

// Random returns a hexidecimal color representation of a color within the
// shade range of the base color.
func (f *Family) Random(seed int64) string {
	return colorful.Hsl(rando(f.Hue), rando(f.Sat), rando(f.Lum)).Hex()
}

func rando(r Range) float64 {
	return (rand.Float64() * (r.Top - r.Bottom)) + r.Bottom
}

// List returns the whole set of the names of the canonical color families.
func List() []string {
	var r []string
	for k := range list {
		r = append(r, k)
	}
	sort.Strings(r)

	return r
}

// NewFamily returns a new shade family for generating random colors.
func NewFamily(key string) (Family, error) {
	var r Family
	r, ok := list[key]
	if !ok {
		return r, ErrNotValidFamily
	}
	return r, nil
}

// FindFamily returns the name of the family for a given color in a range.
func FindFamily(hex string) string {
	for i, v := range list {
		if i == "ALL" {
			continue
		}
		if v.In(hex) {
			return i
		}
	}
	return ""
}

// Invert returns the color on the opposite side of the hue chart
func Invert(hex string) string {
	hex = strings.ToUpper(hex)
	splitnum := strings.Split(hex, "")
	resultnum := "#"
	simplenum := strings.Split("FEDCBA9876", "")
	complexnum := make(map[string]string)
	complexnum["A"] = "5"
	complexnum["B"] = "4"
	complexnum["C"] = "3"
	complexnum["D"] = "2"
	complexnum["E"] = "1"
	complexnum["F"] = "0"

	for i := 0; i < 7; i++ {
		if isNumeric(splitnum[i]) {
			num, _ := strconv.Atoi(splitnum[i])
			resultnum += simplenum[num]
		} else {
			resultnum += complexnum[splitnum[i]]
		}
	}

	return resultnum
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
