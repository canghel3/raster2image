package render

import (
	"image"
	"image/color"
)

type GrayscaleRenderer struct {
	width  int
	height int

	data []float64

	min float64
	max float64
}

func Grayscale(data []float64, width, height int, min, max float64) Renderer {
	return &GrayscaleRenderer{
		width:  width,
		height: height,
		data:   data,
		min:    min,
		max:    max,
	}
}

func (gr *GrayscaleRenderer) Draw() (image.Image, error) {
	img := image.NewGray(image.Rect(0, 0, gr.width, gr.height))

	// Normalize and apply the color map
	for y := 0; y < gr.height; y++ {
		for x := 0; x < gr.width; x++ {
			img.SetGray(x, y, color.Gray{Y: uint8(gr.data[y*gr.width+x])})
		}
	}
	return img, nil
}
