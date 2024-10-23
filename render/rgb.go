package render

import "image"

type RGBRenderer struct {
	width  int
	height int

	data [][][]float64

	min float64
	max float64
}

func RGB(data [][][]float64, width, height int, min, max float64) RGBRenderer {
	return RGBRenderer{
		width:  width,
		height: height,
		data:   data,
		min:    min,
		max:    max,
	}
}

func (rr *RGBRenderer) Render() (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, rr.width, rr.height))

	rangeVal := rr.max - rr.min
	if rangeVal == 0 {
		rangeVal = 1
	}

	// Normalize and apply the color map
	for y := 0; y < rr.height; y++ {
		for x := 0; x < rr.width; x++ {
			//value := rr.data[y*rr.width+x]
			// Normalize the value between 0 and 255
			//normalized := uint8((value - r.min) / (r.max - r.min) * 255)
			//normalized := uint8((value - r.min) / rangeVal * 255)
			//col := applyColorMap(uint8(value))
			//img.Set(x, y, col)
		}
	}
	return img, nil
}
