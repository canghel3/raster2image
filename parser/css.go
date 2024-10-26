package parser

import (
	"fmt"
	"os"
	"strings"
)

// ColorMapEntry represents each color map entry in the raster-color-map
type ColorMapEntry struct {
	Color    string  // Hex color code
	Quantity int     // Quantity associated with the color
	Opacity  float64 // Opacity value
	Label    string  // Description label
}

// RasterStyle represents the entire raster style configuration
type RasterStyle struct {
	RasterChannels string          // Channel setting
	ColorMap       []ColorMapEntry // List of color map entries
}

type CSSParser struct {
	path string
}

func NewCSSParser(path string) *CSSParser {
	return &CSSParser{
		path: path,
	}
}

func (cp *CSSParser) Parse() (*RasterStyle, error) {
	content, err := os.ReadFile(cp.path)
	if err != nil {
		return nil, err
	}

	style := &RasterStyle{}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Set raster channels
		if strings.HasPrefix(line, "raster-channels") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				style.RasterChannels = strings.TrimSuffix(strings.TrimSpace(parts[1]), ";")
			}
		}

		// Parse color map entries
		if strings.HasPrefix(line, "color-map-entry") {
			line = strings.TrimPrefix(line, "color-map-entry(")
			line = strings.TrimSuffix(line, ")")
			parts := strings.Split(line, ",")

			// Check for valid number of parts in the color-map-entry
			if len(parts) != 4 {
				return nil, fmt.Errorf("invalid color-map-entry format")
			}

			// Parse each part
			color := strings.TrimSpace(parts[0])
			quantity := 0
			fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &quantity)
			opacity := 1.0
			fmt.Sscanf(strings.TrimSpace(parts[2]), "%f", &opacity)
			label := strings.Trim(strings.TrimSpace(parts[3]), `"`)

			// Append new color map entry
			style.ColorMap = append(style.ColorMap, ColorMapEntry{
				Color:    color,
				Quantity: quantity,
				Opacity:  opacity,
				Label:    label,
			})
		}
	}

	return style, nil
}
