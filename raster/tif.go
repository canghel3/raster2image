package raster

import (
	"fmt"
	"github.com/canghel3/raster2image/render"
	"image"
)

type TifDriver struct {
	gd *GodalDataset
}

func NewTifDriver(gd *GodalDataset) *TifDriver {
	return &TifDriver{
		gd: gd,
	}
}

// TODO: check if gdal_translate to PNG and styling is possible
func (td *TifDriver) Render(width, height uint) (image.Image, error) {
	switch len(td.gd.data.ds.Bands()) {
	case 1:
		return td.renderSingleBand(width, height)
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 bands", td.gd.path)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 bands", td.gd.path)
	}

	return nil, nil
}

func (td *TifDriver) renderSingleBand(width, height uint) (image.Image, error) {
	if td.gd.data.min == 0 && td.gd.data.max == 255 {
		//grayscale (or apply style)
		if td.gd.data.style == nil {
			band := td.gd.data.ds.Bands()[0]
			var data = make([]float64, width*height)
			err := band.Read(0, 0, data, int(width), int(height))
			if err != nil {
				return nil, err
			}

			grayscale := render.Grayscale(data, int(width), int(height), td.gd.data.min, td.gd.data.max)
			return grayscale.Draw()
		} else {
			//style given, so use rgb renderer with the style schema
			return nil, fmt.Errorf("not implemented")
		}
	} else {
		//check if style is set
		//if style is set, check that the dataset values are within the style ranges
		//if no style set, normalize the values in uint8 range
		return nil, fmt.Errorf("cannot render raster %s with values larger than 255", td.gd.path)
	}
}
