package main

import (
	"net/http"

	// "github.com/ajstarks/svgo"

	"fmt"
	"image"
	"image/jpeg"
	"os"

	// "encoding/base64"
	c "image/color"

	"bytes"
	"image/color"

	"github.com/anthonynsimon/bild/adjust"
	"sync"
)

func colorchange(w http.ResponseWriter, req *http.Request) {
	// saveVariations()
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(bit4(brightness1))
}

var brightness1 = 0.6
var file1 = "./matterhorn.jpg"

// var file = "./everest.jpg"
// var file = "./foto.jpg"

func saveVariations() {
	oldB := brightness1
	bs := [50]float64{}
	var wg sync.WaitGroup
	for x := 1; x < 50; x++ {
		b := -1 + float64(2.0*(float64(x)/50.0))
		go func(x int, b float64) {
			wg.Add(1)
			img := bit4(b)
			f, err := os.Create(fmt.Sprintf("./img%v.jpg", x))
			if err != nil {
				fmt.Println(err)
			}
			f.Write(img)
			f.Close()
			bs[x] = b
			wg.Done()
		}(x, b)
	}
	brightness1 = oldB
	wg.Wait()
	fmt.Println(bs)
}

func bit4(b float64) []byte {
	freader, err := os.Open(file1)
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer freader.Close()

	img, err := jpeg.Decode(freader)
	img = adjust.Brightness(img, b)

	if err != nil {
		fmt.Printf("Error: Image could not be decoded: %v\n", err)
		os.Exit(1)
	}

	// w.Header().Set("Content-Type", "image/svg+xml")
	// s := svg.New(w)
	// s.Start(500, 500)

	newImg := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldPixel := img.At(x, y)
			pixel := color.Gray16Model.Convert(oldPixel)
			newImg.Set(x, y, pixel)
		}
	}

	bounds := img.Bounds()

	pixels := make([][]uint32, bounds.Max.Y)
	// var r, g, b, a uint32
	for y := 0; y < bounds.Max.Y; y++ {
		slice := make([]uint32, bounds.Max.X)
		for x := 0; x < bounds.Max.X; x++ {
			r, _, _, _ := newImg.At(x, y).RGBA()
			// r, g, b, a = newImg.At(x+500, y+500).RGBA()
			// _ = s.RGBA(int(r>>8), int(g>>8), int(b>>8), float64(a>>8))
			g := (r >> 14)
			slice[x] = g
		}
		pixels[y] = slice
	}

	color := false
	var green, red, blue, yellow c.RGBA
	if color {
		green, _ = p.Convert("green")
		red, _ = p.Convert("red")
		blue, _ = p.Convert("blue")
		yellow, _ = p.Convert("yellow")
	} else {
		green, _ = p.Convert("#222")
		red, _ = p.Convert("#444")
		blue, _ = p.Convert("#ccc")
		yellow, _ = p.Convert("#ddd")
	}

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			g := pixels[y][x]
			switch g {
			case 0:
				// s.Rect(x, y, 1, 1, fmt.Sprintf("fill-opacity:255.00; fill:green"))
				newImg.Set(x, y, green)
			case 1:
				// s.Rect(x, y, 1, 1, fmt.Sprintf("fill-opacity:255.00; fill:red"))
				newImg.Set(x, y, red)
			case 2:
				// s.Rect(x, y, 1, 1, fmt.Sprintf("fill-opacity:255.00; fill:blue"))
				newImg.Set(x, y, blue)
			case 3:
				// s.Rect(x, y, 1, 1, fmt.Sprintf("fill-opacity:255.00; fill:yellow"))
				// newImg.Set(x, y, yellow)
			}
		}
	}
	newImg.Set(0, 0, yellow)

	// s.End()
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImg, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		fmt.Println(err)
	}

	return buf.Bytes()
}
