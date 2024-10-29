package render

import (
	"image"
	"image/color"
)

type RGBRenderer struct {
	width  int
	height int

	data []float64

	min float64
	max float64

	colorMap func(float64) color.RGBA
}

func RGB(data []float64, width, height int, min, max float64, options ...RGBRendererOption) RGBRenderer {
	r := RGBRenderer{
		width:  width,
		height: height,
		data:   data,
		min:    min,
		max:    max,
	}

	for _, option := range options {
		option(&r)
	}

	return r
}

type RGBRendererOption func(*RGBRenderer)

func ColorMapOption(colorMap func(float64) color.RGBA) RGBRendererOption {
	return func(r *RGBRenderer) {
		r.colorMap = colorMap
	}
}

func (rr *RGBRenderer) Draw() (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, rr.width, rr.height))

	rangeVal := rr.max - rr.min
	if rangeVal == 0 {
		rangeVal = 1
	}

	// Normalize and apply the color map
	for y := 0; y < rr.height; y++ {
		for x := 0; x < rr.width; x++ {
			value := rr.data[y*rr.width+x]
			img.Set(x, y, rr.colorMap(value))
		}
	}
	return img, nil
}
