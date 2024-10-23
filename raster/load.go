package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/render"
	"image"
	"math"
	"path/filepath"
	"sync"
)

var R *Registry

type Registry struct {
	mx       sync.Mutex
	registry map[string]Data
}

type Data struct {
	style string
	min   float64
	max   float64
	ds    *godal.Dataset
}

func init() {
	R = &Registry{
		registry: make(map[string]Data),
		mx:       sync.Mutex{},
	}
}

// Load reads the given raster file and stores it into the registry.
func Load(path string) error {
	ds, err := godal.Open(path)
	if err != nil {
		return err
	}

	min, max, err := minMaxDs(ds)
	if err != nil {
		return err
	}

	R.mx.Lock()
	R.registry[filepath.Base(path)] = Data{
		ds:    ds,
		min:   min,
		max:   max,
		style: "",
	}
	R.mx.Unlock()

	return nil
}

func Release(path string) {
	R.mx.Lock()
	delete(R.registry, filepath.Base(path))
	R.mx.Unlock()
}

func Render(path string, width, height int, switches []string) (image.Image, error) {
	base := filepath.Base(path)
	rd, exists := R.registry[base]
	if !exists {
		return nil, fmt.Errorf("%s not loaded", base)
	}

	// Create a new dataset by warping
	warped, err := rd.ds.Warp("", switches)
	if err != nil {
		return nil, err
	}

	switch len(warped.Bands()) {
	case 1:
		//grayscale (or apply style)
		if len(rd.style) == 0 {
			band := warped.Bands()[0]
			var data = make([]float64, width*height)
			err = band.Read(0, 0, data, width, height)
			if err != nil {
				return nil, err
			}

			grayscale := render.Grayscale(data, width, height, rd.min, rd.max)
			return grayscale.Render()
		}
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 bands", base)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 bands", base)
	}

	return nil, nil
}

func minMaxDs(ds *godal.Dataset) (min, max float64, err error) {
	switch len(ds.Bands()) {
	case 1:
		band := ds.Bands()[0]
		bandStructure := band.Structure()

		var data = make([]float64, bandStructure.SizeX*bandStructure.SizeY)
		err := band.Read(0, 0, data, bandStructure.SizeX, bandStructure.SizeY)
		if err != nil {
			return min, max, err
		}

		min, max = minMax(data)
	}

	return min, max, nil
}

func minMax(data []float64) (min, max float64) {
	min, max = math.MaxFloat64, -math.MaxFloat64
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}
