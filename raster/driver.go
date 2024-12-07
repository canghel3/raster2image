package raster

import (
	"github.com/canghel3/raster2image/models"
	"image"
)

type Driver interface {
	Render(bbox [4]float64, width, height uint) (image.Image, error)
	Release() error
	setStyle(style *models.RasterStyle)
}
