package raster

import (
	"github.com/airbusgeo/godal"
	"image"
)

type Driver interface {
	Render(bands []godal.Band, width, height uint) (image.Image, error)
}
