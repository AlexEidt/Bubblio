package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// Create output directory "Text" to store Bubble Letters
	outputDir := filepath.Join(".", "Text")
	os.MkdirAll(outputDir, os.ModePerm)

	// Parse command line args
	font := flag.String("font", "helvetica", "Font to write letters in.")
	shape := flag.String("shape", "circle", "Component shapes of letters")
	sides := flag.Int("sides", 5, "Number of sides for the Polygon.")
	palette := flag.String("color", "blue", "Color palette for letters")
	scale := flag.Int("scale", 50, "Relative size of components")
	img := flag.Bool("img", false, "Include PNG image if true (for use with -animate)")
	animate := flag.Bool("animate", false, "Create animated GIF with letters")
	frames := flag.Int("frames", 10, "Number of frames in animated GIF")
	orientation := flag.String("o", "L", "Orienation. [L]eft, [R]ight, or [C]entered")

	flag.Parse()

	args := flag.Args()

	input := args[0]
	splitter := args[1]
	filename := args[2]

	symbols := ParseFont(filepath.Join("Fonts", *font+".yaff"))

	unknowns := make([]string, 0, 0)
	check := true
	for _, letter := range input {
		if _, ok := symbols[letter]; !ok {
			unknowns = append(unknowns, string(letter))
			check = false
		}
	}
	if !check {
		fmt.Println("The following characters could not be processed:")
		fmt.Println(strings.Join(unknowns, ", "))
	} else {
		scale := *scale
		width, maxwidth, height, lineheight := 0, 0, 0, 0
		lengths := make(map[string]int)
		lines := strings.Split(input, splitter)
		for _, line := range lines {
			for _, letter := range line {
				s := symbols[letter]
				width += (s.width + 1) * scale
				lineheight = s.height
			}
			lengths[line] = width
			if width > maxwidth {
				maxwidth = width
			}
			width = 0
			height += (lineheight + 2) * scale
		}
		maxwidth += 4 * scale
		outputfile := filename
		iterations := 1
		if *animate {
			iterations = *frames
		}
		splitname := strings.Split(filename, ".")[0]
		for i := 0; i < iterations; i++ {
			if *animate {
				outputfile = splitname + strconv.Itoa(i) + ".png"
			}
			DrawCharacters(
				&symbols,
				&lengths,
				lines,
				*palette,
				scale,
				maxwidth,
				height,
				*orientation,
				*shape,
				*sides,
				outputfile,
			)
		}
		if *animate {
			CreateGIF(splitname, *frames, *img)
		}
	}
}

// Creates a GIF from PNG files of the form: [filename][count].png
// for all count from 0 to count - 1.
// Code inspired from Hirmou Ochiai's GIFFY library on GitHub:
// https://github.com/otiai10/giffy
func CreateGIF(filename string, count int, keep bool) {
	g := &gif.GIF{}

	for i := 0; i < count; i++ {
		fname := filepath.Join("Text", filename+strconv.Itoa(i)+".png")
		f, _ := os.Open(fname)
		if !keep {
			os.Remove(fname)
		}
		defer f.Close()

		img, _, _ := image.Decode(f)

		p := image.NewPaletted(
			img.Bounds(),
			color.Palette{
				image.Transparent,
				color.Black,
				color.RGBA{66, 135, 245, 255},
				color.RGBA{27, 27, 196, 255},
				color.RGBA{27, 159, 196, 255},
				color.RGBA{99, 190, 242, 255},
				color.RGBA{25, 0, 255, 255},
				color.RGBA{119, 255, 0, 255},
				color.RGBA{199, 255, 150, 255},
				color.RGBA{62, 240, 22, 255},
				color.RGBA{105, 255, 71, 255},
				color.RGBA{255, 255, 71, 255},
				color.RGBA{255, 71, 71, 255},
				color.RGBA{255, 0, 0, 255},
				color.RGBA{199, 95, 10, 255},
				color.RGBA{235, 85, 35, 255},
				color.RGBA{158, 41, 2, 255},
				color.RGBA{158, 2, 132, 255},
				color.RGBA{255, 41, 219, 255},
				color.RGBA{247, 109, 224, 255},
				color.RGBA{84, 16, 173, 255},
				color.RGBA{52, 19, 102, 255},
				color.RGBA{0, 36, 38, 255},
				color.RGBA{133, 166, 168, 255},
				color.RGBA{179, 224, 227, 255},
				color.RGBA{0, 0, 0, 255},
				color.RGBA{10, 46, 42, 255},
			},
		)
		draw.Draw(p, p.Bounds(), img, img.Bounds().Min, draw.Over)
		g.Image = append(g.Image, p)
		g.Delay = append(g.Delay, 100)
	}
	f, err := os.Create(filepath.Join("Text", filename+".gif"))
	if err != nil {
		panic(err)
	}
	gif.EncodeAll(f, g)
}
