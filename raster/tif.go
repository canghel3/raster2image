package raster

import (
	"fmt"
	"github.com/canghel3/raster2image/render"
	"image"
	"image/color"
)

type TifDriver struct {
	gd *GodalDataset
}

func NewTifDriver(gd *GodalDataset) Driver {
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
		band := td.gd.data.ds.Bands()[0]
		var data = make([]float64, width*height)
		err := band.Read(0, 0, data, int(width), int(height))
		if err != nil {
			return nil, err
		}

		if td.gd.data.style == nil {
			grayscale := render.Grayscale(data, int(width), int(height), td.gd.data.min, td.gd.data.max)
			return grayscale.Draw()
		} else {
			//style given, so use rgb renderer with the style schema
			colorMap := td.generateColorMap()
			rgb := render.RGB(data, int(width), int(height), td.gd.data.min, td.gd.data.max, render.ColorMapOption(colorMap))
			return rgb.Draw()
		}
	} else {
		//check if style is set
		//if style is set, check that the dataset values are within the style ranges
		//if no style set, normalize the values in uint8 range
		return nil, fmt.Errorf("cannot render raster %s with values larger than 255", td.gd.path)
	}
}

func (td *TifDriver) generateColorMap() func(float64) color.RGBA {
	return func(f float64) color.RGBA {
		var previous = td.gd.data.style.ColorMap[0]
		for i, colorEntry := range td.gd.data.style.ColorMap {
			if i == 0 {
				if f <= float64(colorEntry.Quantity) {
					return hexToRGBA(colorEntry.Color)
				}
			} else if i == len(td.gd.data.style.ColorMap)-1 {
				if f >= float64(colorEntry.Quantity) {
					return hexToRGBA(colorEntry.Color)
				}
			} else {
				if f >= float64(previous.Quantity) && f <= float64(colorEntry.Quantity) {
					return hexToRGBA(colorEntry.Color)
				}
			}

			previous = td.gd.data.style.ColorMap[i]
		}

		return color.RGBA{}
	}
}
