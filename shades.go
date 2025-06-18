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

// Color is an enum that makes it easy to reference the pre set color values.
type Color int64

const (
	// Red 🟥
	Red Color = iota + 1
	// Orange 🟧
	Orange
	// Yellow 🟨
	Yellow
	// Green 🟩
	Green
	// Cyan 🔵🟢
	Cyan
	// Blue 🟦
	Blue
	// Purple 🔵🔴
	Purple
	// Magenta 🟪
	Magenta
	// All 🌈
	All
)

func (c Color) String() string {
	switch c {

	case Red:
		return "RED"
	case Orange:
		return "ORANGE"
	case Yellow:
		return "YELLOW"
	case Green:
		return "GREEN"
	case Cyan:
		return "CYAN"
	case Blue:
		return "BLUE"
	case Purple:
		return "PURPLE"
	case Magenta:
		return "MAGENTA"
	case All:
		return "ALL"
	}
	return "unknown"
}

var list = map[string]Family{
	"RED": {
		Name: "Red",
		Base: "FF0000",
		Hue:  Range{-10, 20},
		Sat:  Range{.2, 1},
		Lum:  Range{.2, 1},
	},
	"ORANGE": {
		Name: "Orange",
		Base: "FFA500",
		Hue:  Range{21, 50},
		Sat:  Range{.3, 1},
		Lum:  Range{.4, 1},
	},
	"YELLOW": {
		Name: "Yellow",
		Base: "FFFF00",
		Hue:  Range{51, 60},
		Sat:  Range{.4, 1},
		Lum:  Range{.63, 1},
	},
	"GREEN": {
		Name: "Green",
		Base: "00FF00",
		Hue:  Range{81, 140},
		Sat:  Range{.4, 1},
		Lum:  Range{.3, .8},
	},
	"CYAN": {
		Name: "Cyan",
		Base: "00FFFF",
		Hue:  Range{170, 200},
		Sat:  Range{.25, 1},
		Lum:  Range{.3, 1},
	},
	"BLUE": {
		Name: "Blue",
		Base: "0000FF",
		Hue:  Range{221, 240},
		Sat:  Range{.1, 1},
		Lum:  Range{.2, 1},
	},
	"PURPLE": {
		Name: "Purple",
		Base: "800080",
		Hue:  Range{241, 280},
		Sat:  Range{.3, 1},
		Lum:  Range{.4, .7},
	},
	"MAGENTA": {
		Name: "Magenta",
		Base: "FF00FF",
		Hue:  Range{281, 320},
		Sat:  Range{.35, 1},
		Lum:  Range{.3, .7},
	},
	"ALL": {
		Name: "All",
		Base: "FF00FF",
		Hue:  Range{0, 360},
		Sat:  Range{0, 1},
		Lum:  Range{0, 1},
	},
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

// NewFamily returns a new shade family for generating random colors.
func NewFamily(opts ...Option) Family {
	f := &Family{}

	if len(opts) <= 0 {
		opts = append(opts, Only(All))
	}

	for _, opt := range opts {
		opt(f)
	}

	return *f
}

// Option allows us to configure modifications for family
type Option func(f *Family)

func Only(c Color) Option {
	return func(f *Family) {
		tmp := list[c.String()]
		f = &tmp
	}
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
func (f *Family) Random() string {
	return colorful.Hsl(rando(f.Hue), rando(f.Sat), rando(f.Lum)).Hex()
}

func rando(r Range) float64 {
	answer := (rand.Float64() * (r.Top - r.Bottom)) + r.Bottom
	return answer
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

// IsGreyScale will repourt if the colour is greyscale (RGB match)
func IsGreyScale(hex string) bool {
	return IsGrayScale(hex)
}

// IsGrayScale will report if the color is grayscale (RGB match)
func IsGrayScale(hex string) bool {
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
