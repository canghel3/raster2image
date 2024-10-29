package raster

import (
	"github.com/airbusgeo/godal"
	"image/color"
	"math"
	"strconv"
)

func minMaxDs(ds *godal.Dataset) (min, max float64, err error) {
	switch len(ds.Bands()) {
	case 1:
		band := ds.Bands()[0]
		bandStructure := band.Structure()

		var data = make([]float64, bandStructure.SizeX*bandStructure.SizeY)
		err := band.Read(0, 0, data, bandStructure.SizeX, bandStructure.SizeY)
		if err != nil {
			return min, max, err
		}

		min, max = minMax(data)
	}

	return min, max, nil
}

func minMax(data []float64) (min, max float64) {
	min, max = math.MaxFloat64, -math.MaxFloat64
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

func hexToRGBA(hex string) color.RGBA {
	// Handle empty string or non-hex format
	if len(hex) == 0 {
		return color.RGBA{}
	}

	// Remove the leading '#' if it exists
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b, a uint8
	a = 255 // default alpha is 100% opaque

	switch len(hex) {
	case 3:
		// 3-digit hex (RGB)
		r = parseHexToByte(hex[0:1] + hex[0:1])
		g = parseHexToByte(hex[1:2] + hex[1:2])
		b = parseHexToByte(hex[2:3] + hex[2:3])
	case 6:
		// 6-digit hex (RRGGBB)
		r = parseHexToByte(hex[0:2])
		g = parseHexToByte(hex[2:4])
		b = parseHexToByte(hex[4:6])
	case 8:
		// 8-digit hex (RRGGBBAA)
		r = parseHexToByte(hex[0:2])
		g = parseHexToByte(hex[2:4])
		b = parseHexToByte(hex[4:6])
		a = parseHexToByte(hex[6:8])
	default:
		// Return default color if the format is invalid
		return color.RGBA{}
	}

	return color.RGBA{R: r, G: g, B: b, A: a}
}

func parseHexToByte(hexStr string) uint8 {
	value, err := strconv.ParseUint(hexStr, 16, 8)
	if err != nil {
		return 0
	}
	return uint8(value)
}
