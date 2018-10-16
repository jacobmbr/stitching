package main

import (
	"net/http"

	"github.com/ajstarks/svgo"

	"fmt"
	"image"
	"image/jpeg"

	// "encoding/base64"
	"image/color"
	"os"

	"github.com/anthonynsimon/bild/adjust"
)

var scale = 5
// var file = "./everest.jpg"
var file = "./matterhorn.jpg"
// var file = "./foto.jpg"
var grey = true
var cols = []string{"blue", "red", "wheat", "yellow"}
var brightness = 0.4

func circle(w http.ResponseWriter, req *http.Request) {
	freader, err := os.Open(file)
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer freader.Close()

	img, err := jpeg.Decode(freader)
	bounds := img.Bounds()

	img = adjust.Brightness(img, brightness)

	fmt.Println(img.Bounds())

	if err != nil {
		fmt.Printf("Error: Image could not be decoded: %v\n", err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(bounds.Max.X, bounds.Max.Y)

	newImg := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldPixel := img.At(x, y)
			pixel := color.Gray16Model.Convert(oldPixel)
			newImg.Set(x, y, pixel)
		}
	}

	pixels := make([][]uint32, bounds.Max.Y)
	var r, g, b, a uint32
	for y := 0; y < bounds.Max.Y; y++ {
		slice := make([]uint32, bounds.Max.X)
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, a = newImg.At(x, y).RGBA()
			_ = s.RGBA(int(r>>8), int(g>>8), int(b>>8), float64(a>>8))
			g := (r >> 14)
			slice[x] = g
		}
		pixels[y] = slice
	}

	lines(s, bounds, pixels, scale, grey)
	s.End()
}

func line(s *svg.SVG, typ int, x,y int, bounds image.Rectangle,g bool) {
	s1 := scale
	s2 := scale / 2

	switch typ {
	case 0:
		s.Line(x, y, x+s1, y+s1, fmt.Sprintf("fill-opacity:255.00;stroke-width:1px; stroke:%v", col(0, g)))
	case 1:
		s.Line(x, y+s2, x+s1, y+s2, fmt.Sprintf("fill-opacity:255.00;stroke-width:1px; stroke:%v", col(1, g)))
	case 2:
		s.Line(x, y+s1, x+s1, y, fmt.Sprintf("fill-opacity:255.00;stroke-width:1px; stroke:%v", col(2, g)))
	case 3:
		s.Line(x+s2, y, x+s2, y+s1, fmt.Sprintf("fill-opacity:255.00;stroke-width:1px; stroke:%v", col(3, g)))
	}
}

func col(typ int, grey bool) string {
	if grey {
		return "grey"
	}
	return cols[typ]
}

func lines(s *svg.SVG, b image.Rectangle, pixels [][]uint32, scale int, grey bool) {
	for y := 0; y < b.Max.Y; y += scale {
		for x := 0; x < b.Max.X; x += scale {
			g := pixels[y][x]
			switch g {
			case 0:
				line(s, 3, x, y, b, grey)
				line(s, 2, x, y, b, grey)
				line(s, 1, x, y, b, grey)
				line(s, 0, x, y, b, grey)
			case 1:
				line(s, 3, x, y, b, grey)
				line(s, 2, x, y, b, grey)
				line(s, 1, x, y, b, grey)
			case 2:
				line(s, 3, x, y, b, grey)
				line(s, 2, x, y, b, grey)
			case 3:
			}
		}
	}
}
