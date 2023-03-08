# Shades
[![GoDoc](https://godoc.org/github.com/tpryan/shades?status.svg)](https://godoc.org/github.com/tpryan/shades)
[![Go](https://github.com/tpryan/shades/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/tpryan/shades/actions/workflows/go.yml)

This is a minimal golang package and app for generating random colors.


## Usage

```go
shade := NewFamily(Red)
color := shade.Random()

fmt.Println(color)
fmt.Printf("color: %s\n", color) // #e58677
```

If you run the sample web app you get a minimal random list of colors.

![Colors Screenshot](sample.png "Screenshot")



"This is not an official Google Project."
