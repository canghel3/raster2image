package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/models"
	"github.com/canghel3/raster2image/render"
	"image"
)

type TifDriver struct {
	name  string
	bands []godal.Band
	min   float64
	max   float64
	style *models.RasterStyle
}

type TifDriverData struct {
	Name  string
	Bands []godal.Band
	Min   float64
	Max   float64
	Style *models.RasterStyle
}

func NewTifDriver(data TifDriverData) Driver {
	return &TifDriver{
		name:  data.Name,
		bands: data.Bands,
		max:   data.Max,
		min:   data.Min,
		style: data.Style,
	}
}

func (td *TifDriver) Render(width, height uint) (image.Image, error) {
	switch len(td.bands) {
	case 1:
		return td.renderSingleBand(width, height)
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 Bands", td.name)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 Bands", td.name)
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
		rgb := render.NewRGBDrawer(data, int(width), int(height), render.StyleOption(*td.style))
		return rgb.Draw()
	}

	grayscale := render.Grayscale(data, int(width), int(height), td.min, td.max)
	return grayscale.Draw()
}
