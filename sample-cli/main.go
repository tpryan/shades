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

// Package main is a Kubernetes API proxy. It exposes a smaller surface of the
// API and limits operations to specifically selected labels, and deployments
package main

import (
	"fmt"
	"log"

	shades "github.com/tpryan/shades"
)

func main() {
	shade, err := shades.NewFamily("RED")
	if err != nil {
		log.Fatalf("could not get color family: %v", err)
	}
	color := shade.Random(1)

	fmt.Printf("color: %s\n", color)
}
