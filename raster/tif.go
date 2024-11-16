package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/models"
	"github.com/canghel3/raster2image/render"
	"image"
)

type TifDriver struct {
	bands []godal.Band
	path  string
	min   float64
	max   float64
	style *models.RasterStyle
}

func NewTifDriver(bands []godal.Band, path string, min, max float64, style *models.RasterStyle) RenderDriver {
	return &TifDriver{
		bands: bands,
		path:  path,
		min:   min,
		max:   max,
		style: style,
	}
}

// TODO: check if gdal_translate to PNG and styling is possible
func (td *TifDriver) Render(width, height uint) (image.Image, error) {
	switch len(td.bands) {
	case 1:
		return td.renderSingleBand(width, height)
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 bands", td.path)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 bands", td.path)
	}

	return nil, nil
}

func (td *TifDriver) renderSingleBand(width, height uint) (image.Image, error) {
	band := td.bands[0]
	var data = make([]float64, width*height)
	err := band.Read(0, 0, data, int(width), int(height))
	if err != nil {
		return nil, err
	}

	if td.style != nil {
		//style given, so use rgb renderer with the style schema
		rgb := render.RGB(data, int(width), int(height), td.min, td.max, render.StyleOption(*td.style))
		return rgb.Draw()
	}

	if td.min == 0 && td.max == 255 {
		//grayscale (or apply style)
		grayscale := render.Grayscale(data, int(width), int(height), td.min, td.max)
		return grayscale.Draw()
	} else {

		return nil, fmt.Errorf("cannot render raster %s with values larger than 255", td.path)
	}
}
