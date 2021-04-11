package main

import (
	"image/color"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

var colormap = map[string][]color.Color{
	"blue": {
		color.RGBA{66, 135, 245, 255},
		color.RGBA{27, 27, 196, 255},
		color.RGBA{27, 159, 196, 255},
		color.RGBA{99, 190, 242, 255},
		color.RGBA{25, 0, 255, 255},
	},
	"yellow": {
		color.RGBA{119, 255, 0, 255},
		color.RGBA{199, 255, 150, 255},
		color.RGBA{62, 240, 22, 255},
		color.RGBA{105, 255, 71, 255},
		color.RGBA{255, 255, 71, 255},
	},
	"red": {
		color.RGBA{255, 71, 71, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{199, 95, 10, 255},
		color.RGBA{235, 85, 35, 255},
		color.RGBA{158, 41, 2, 255},
	},
	"purple": {
		color.RGBA{158, 2, 132, 255},
		color.RGBA{255, 41, 219, 255},
		color.RGBA{247, 109, 224, 255},
		color.RGBA{84, 16, 173, 255},
		color.RGBA{52, 19, 102, 255},
	},
	"grayscale": {
		color.RGBA{0, 36, 38, 255},
		color.RGBA{133, 166, 168, 255},
		color.RGBA{179, 224, 227, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{10, 46, 42, 255},
	},
}

func DrawCharacters(
	symbols *map[rune]*symbol,
	lengths *map[string]int,
	lines []string,
	palette string,
	scale int,
	width int,
	height int,
	orientation string,
	shape string,
	sides int,
	filename string,
) {
	im := gg.NewContext(width, height)

	// Seeded random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	scaled := scale / 2
	totalheight := scaled
	lineheight := 0
	var totalwidth int

	for _, line := range lines {
		switch orientation {
		case "L": // Left aligned.
			totalwidth = scale
		case "R": // Rigth aligned.
			totalwidth = width - 3*scale - (*lengths)[line]
		default: // Centered by default.
			totalwidth = width/2 - (*lengths)[line]/2
		}
		for _, char := range line {
			s := (*symbols)[char]
			lineheight = s.height
			DrawLetter(scale, scaled, totalwidth, totalheight, shape, sides, palette, s, r, im)
			totalwidth += s.width*scale + scale
		}
		totalheight += lineheight*scale + scale
	}

	im.SavePNG(filepath.Join("Text", filename))
}

func DrawLetter(
	scale int,
	scaled int,
	width int,
	height int,
	shape string,
	sides int,
	palette string,
	s *symbol,
	r *rand.Rand,
	im *gg.Context,
) {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if s.bitmap[y][x] {
				random := r.Intn(scaled) + 1
				for i := 0; i < random; i++ {
					bgcolor := colormap[palette][r.Intn(len(colormap[palette]))]
					radius := float64(r.Intn(scaled)+scaled) / 3
					xc := float64(x*scale + r.Intn(scale) + width)
					yc := float64(y*scale + scaled + r.Intn(scale) + height)
					DrawShape(xc, yc, radius, scaled, shape, sides, random, bgcolor, im)
				}
			}
		}
	}
}

func DrawShape(
	x float64,
	y float64,
	radius float64,
	scaled int,
	shape string,
	sides int,
	random int,
	bg color.Color,
	im *gg.Context,
) {
	if shape == "random" {
		switch random % 4 {
		case 0:
			shape = "square"
		case 1:
			shape = "triangle"
		case 2:
			shape = "polygon"
		case 3:
			shape = "circle"
		}
	}
	scale := float64(scaled * 2)
	switch shape {
	case "square":
		im.DrawRectangle(x, y, scale, scale)
		im.SetColor(bg)
		im.Fill()
		im.DrawRectangle(x, y, scale, scale)
	case "triangle":
		im.DrawRegularPolygon(3, x, y, radius*2, radius*2)
		im.SetColor(bg)
		im.Fill()
		im.DrawRegularPolygon(3, x, y, radius*2, radius*2)
	case "polygon":
		im.DrawRegularPolygon(5, x, y, radius*2, radius*2)
		im.SetColor(bg)
		im.Fill()
		im.DrawRegularPolygon(sides, x, y, radius*2, radius*2)
	default: // circle
		im.DrawCircle(x, y, radius)
		im.SetColor(bg)
		im.Fill()
		im.DrawCircle(x, y, radius)
	}
	// Draw black border
	im.SetColor(color.Black)
	im.SetLineWidth(float64(scaled / 4))
	im.Stroke()
}
