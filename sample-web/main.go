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
	"net/http"
	"time"

	shades "github.com/tpryan/shades"
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	count := 4

	header := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Document</title>
</head>
<style>
	.square{
		display: inline-block;
		height: 75px; width: 75px;
	}

	.content {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
	}

	.container {
		width: 30%;
		padding: 5px 10px;
	}
</style>
<body>
<div class="content">
`

	footer := `
</div>
</body>
</html>`

	list := shades.List()
	fmt.Fprintln(w, header)
	for _, k := range list {
		shade, err := shades.NewFamily(k)
		if err != nil {
			log.Printf("could not get color family: %v", err)
			continue
		}
		fmt.Fprintln(w, "\t<div class=\"container\">")
		fmt.Fprintf(w, "\t<h1>%s</h1>\n", shade.Name)
		for i := 0; i < count; i++ {
			color := shade.Random(time.Now().UnixNano())
			inverse := shades.Invert(color)
			fmt.Fprintf(w, "\t<div class=\"square\" style=\"background-color: %s; color: %s;\" >%s</div>\n", color, inverse, color)
		}
		fmt.Fprintf(w, "\t</div>\n")
	}
	fmt.Fprintf(w, footer)

}
