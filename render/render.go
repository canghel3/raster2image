package render

type Floats interface {
	[]float32 | []float64
}

type RGBRenderer struct {
	width  int
	height int

	data []float64

	min float64
	max float64
}

func RGB(data []float64, width, height int, min, max float64) RGBRenderer {
	return RGBRenderer{
		width:  width,
		height: height,
		data:   data,
		min:    min,
		max:    max,
	}
}
