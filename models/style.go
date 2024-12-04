package models

import (
	"github.com/canghel3/raster2image/utils"
	"image/color"
)

// ColorMapEntry represents each color map entry in the raster-color-map
type ColorMapEntry struct {
	Color    string  // Hex color code
	Quantity float64 // Quantity associated with the color
	Opacity  float64 // Opacity value
	Label    string  // Description label
}

// RasterStyle represents the entire raster style configuration
type RasterStyle struct {
	RasterChannels string          // Channel setting
	ColorMap       []ColorMapEntry // List of color map entries
}

func (rs *RasterStyle) GetColor(value float64) color.RGBA {
	var previous float64
	for i, entry := range rs.ColorMap {
		if i == 0 {
			previous = entry.Quantity
		}

		if previous < value && value <= entry.Quantity {
			return utils.HexToRGBA(entry.Color)
		}

		previous = entry.Quantity
	}
	// Default to black if no range matches
	return color.RGBA{}
}
