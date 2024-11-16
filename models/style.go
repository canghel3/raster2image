package models

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
