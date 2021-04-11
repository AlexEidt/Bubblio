package main

import (
	"image/color"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

// Mapping of colors to RGB colors in that "range".
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

// Draws all characters (letters) on the PNG.
func DrawCharacters(
	symbols *map[rune]*symbol, // Mapping of letters to bitmaps
	lengths *map[string]int, // Map of lines and their lengths
	lines []string, // Map of lines
	palette string, // The color palette for this text
	scale int, // The density of shapes in letters
	width int, // Width of PNG
	height int, // Height of PNG
	orientation string, // Text alignment
	shape string, // Type of shapes making up letters
	sides int, // Number of sides for polygon (for shape == "polygon" only)
	filename string, // Filename to store PNG in
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

// Draws one bubblified letter on the PNG.
func DrawLetter(
	scale int, // Density of shapes in letters
	scaled int, // Density divided by 2
	width int, // Width of letter
	height int, // Height of letter
	shape string, // Shapes making up letter
	sides int, // Sides of polygon (for shape == "polygon" only)
	palette string, // Color palette for this letter
	s *symbol, // Bitmap of the letter
	r *rand.Rand, // Random number generator object
	im *gg.Context, // PNG Image
) {
	isrand := palette == "random"
	var bgcolor color.Color
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if s.bitmap[y][x] {
				random := r.Intn(scaled) + 1
				for i := 0; i < random; i++ {
					if isrand {
						// Choose random key from colormap
						for k := range colormap {
							bgcolor = colormap[k][r.Intn(len(colormap["blue"]))]
							break
						}
					} else {
						bgcolor = colormap[palette][r.Intn(len(colormap[palette]))]
					}
					radius := float64(r.Intn(scaled)+scaled) / 3
					xc := float64(x*scale + r.Intn(scale) + width)
					yc := float64(y*scale + scaled + r.Intn(scale) + height)
					DrawShape(xc, yc, radius, scaled, shape, sides, random, bgcolor, im)
				}
			}
		}
	}
}

// Draws the shape given by the "shape" parameter onto the PNG.
func DrawShape(
	x float64, // X coordinate of center
	y float64, // Y coordinate of center
	radius float64, // radius (for circles and polygons)
	scaled int, // Relative size of shape
	shape string, // Type of shape
	sides int, // Number of sides for shape == "polygon"
	random int, // Random number of shape == "random"
	bg color.Color, // Background color of shape
	im *gg.Context, // PNG Image
) {
	if shape == "random" {
		switch random % 4 {
		case 0:
			shape = "square"
		case 1:
			shape = "triangle"
		case 2:
			shape = "polygon"
			sides = random%7 + 5
		case 3:
			shape = "circle"
		}
	}
	radius2 := radius * 2
	switch shape {
	case "square":
		radius3 := radius * 3
		im.DrawRectangle(x, y, radius3, radius3)
		im.SetColor(bg)
		im.Fill()
		im.DrawRectangle(x, y, radius3, radius3)
	case "triangle":
		im.DrawRegularPolygon(3, x, y, radius2, radius2)
		im.SetColor(bg)
		im.Fill()
		im.DrawRegularPolygon(3, x, y, radius2, radius2)
	case "polygon":
		im.DrawRegularPolygon(sides, x, y, radius2, radius2)
		im.SetColor(bg)
		im.Fill()
		im.DrawRegularPolygon(sides, x, y, radius2, radius2)
	default: // circle
		im.DrawCircle(x, y, radius*3/2)
		im.SetColor(bg)
		im.Fill()
		im.DrawCircle(x, y, radius*3/2)
	}
	// Draw black border
	im.SetColor(color.Black)
	im.SetLineWidth(float64(scaled / 3))
	im.Stroke()
}
