package parser

import (
	"errors"
	"fmt"
	"github.com/canghel3/raster2image/models"
	"os"
	"strings"
)

type CSSParser struct {
	path string
}

func NewCSSParser(path string) StyleParser {
	return &CSSParser{
		path: path,
	}
}

func (cp *CSSParser) ParseRasterStyle() (*models.RasterStyle, error) {
	content, err := os.ReadFile(cp.path)
	if err != nil {
		return nil, err
	}

	style := &models.RasterStyle{}
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
			_, err = fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &quantity)
			if err != nil {
				return nil, err
			}

			opacity := 1.0
			_, err = fmt.Sscanf(strings.TrimSpace(parts[2]), "%f", &opacity)
			if err != nil {
				return nil, err
			}

			label := strings.Trim(strings.TrimSpace(parts[3]), `"`)

			// Append new color map entry
			style.ColorMap = append(style.ColorMap, models.ColorMapEntry{
				Color:    color,
				Quantity: quantity,
				Opacity:  opacity,
				Label:    label,
			})
		}
	}

	return style, nil
}

func (cp *CSSParser) ParseVectorStyle() (*models.RasterStyle, error) {
	return nil, errors.New("not implemented")
}
