package raster

import "image"

type Driver interface {
	Render(width, height uint) (image.Image, error)
}
