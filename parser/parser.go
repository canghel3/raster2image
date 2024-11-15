package parser

import "github.com/canghel3/raster2image/models"

type StyleParser interface {
	RasterStyleParser
	VectorStyleParser
}

type RasterStyleParser interface {
	ParseRasterStyle() (*models.RasterStyle, error)
}

type VectorStyleParser interface {
	ParseVectorStyle() (*models.RasterStyle, error)
}
