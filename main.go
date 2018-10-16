package main

import (
	"net/http"
	"fmt"
	"image"
	"log"
	parser "github.com/wayneashleyberry/css-color"
	"image/jpeg"
)

import(
	"os"
)

var p *parser.Parser

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	p = parser.New()
}

func main() {
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	http.Handle("/circle", http.HandlerFunc(circle))
	http.Handle("/color", http.HandlerFunc(colorchange))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
