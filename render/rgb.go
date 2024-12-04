package render

import (
	"github.com/canghel3/raster2image/models"
	"github.com/canghel3/raster2image/utils"
	"image"
	"image/color"
)

type RGBDrawer struct {
	width  int
	height int

	data []float64

	styling models.RasterStyle
}

func NewRGBDrawer(data []float64, width, height int, options ...RGBRendererOption) Drawer {
	r := RGBDrawer{
		width:  width,
		height: height,
		data:   data,
	}

	for _, option := range options {
		option(&r)
	}

	return &r
}

type RGBRendererOption func(*RGBDrawer)

func StyleOption(entry models.RasterStyle) RGBRendererOption {
	return func(r *RGBDrawer) {
		r.styling = entry
	}
}

func (rr *RGBDrawer) Draw() (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, rr.width, rr.height))

	//apply the color map
	for y := 0; y < rr.height; y++ {
		for x := 0; x < rr.width; x++ {
			value := rr.data[y*rr.width+x]
			img.Set(x, y, rr.getColor(value))
		}
	}
	return img, nil
}

// TODO: style is not applied correctly
func (rr *RGBDrawer) getColor(value float64) color.RGBA {
	var previous float64
	for i, entry := range rr.styling.ColorMap {
		if i == 0 {
			previous = entry.Quantity
		}

		if previous < value && value <= entry.Quantity {
			return utils.HexToRGBA(entry.Color)
		}

		previous = entry.Quantity
	}
	// Default to black if no range matches
	return color.RGBA{0, 0, 0, 0}
}
