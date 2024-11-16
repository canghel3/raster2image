package render

import (
	"github.com/canghel3/raster2image/helpers"
	"github.com/canghel3/raster2image/models"
	"image"
	"image/color"
)

type RGBRenderer struct {
	width  int
	height int

	data []float64

	min float64
	max float64

	styling models.RasterStyle
}

func RGB(data []float64, width, height int, min, max float64, options ...RGBRendererOption) Renderer {
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

	return &r
}

type RGBRendererOption func(*RGBRenderer)

func StyleOption(entry models.RasterStyle) RGBRendererOption {
	return func(r *RGBRenderer) {
		r.styling = entry
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
			img.Set(x, y, rr.getColor(value))
		}
	}
	return img, nil
}

func (rr *RGBRenderer) getColor(value float64) color.RGBA {
	for _, entry := range rr.styling.ColorMap {
		if value <= entry.Quantity {
			return helpers.HexToRGBA(entry.Color)
		}
	}
	// Default to black if no range matches
	return color.RGBA{0, 0, 0, 0}
}
