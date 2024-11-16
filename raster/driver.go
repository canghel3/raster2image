package raster

import "image"

type RenderDriver interface {
	Render(width, height uint) (image.Image, error)
}
