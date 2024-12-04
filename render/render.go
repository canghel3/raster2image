package render

import "image"

type Drawer interface {
	Draw() (image.Image, error)
}
