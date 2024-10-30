package render

import "image"

type Renderer interface {
	Draw() (image.Image, error)
}
