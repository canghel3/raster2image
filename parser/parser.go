package parser

import "github.com/canghel3/raster2image/models"

type StyleParser interface {
	Parse() (*models.RasterStyle, error)
}
