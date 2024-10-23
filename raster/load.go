package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"image"
	"sync"
)

var R *Registry

type RasterData struct {
	bands map[int][]float64
}

type Registry struct {
	mx       sync.Mutex
	registry map[string]RasterData
}

func init() {
	R = &Registry{
		registry: make(map[string]RasterData),
		mx:       sync.Mutex{},
	}
}

// Load reads the given raster file and stores it into the registry.
func Load(filename string) error {
	ds, err := godal.Open(filename)
	if err != nil {
		return err
	}
	defer ds.Close()

	rd := RasterData{
		bands: make(map[int][]float64),
	}

	bands := ds.Bands()
	for i := range len(bands) {
		current := bands[i]
		currentStructure := current.Structure()
		var bandData = make([]float64, currentStructure.SizeX*currentStructure.SizeY)
		err = current.Read(0, 0, bandData, currentStructure.SizeX, currentStructure.SizeY)
		if err != nil {
			return err
		}

		rd.bands[i] = bandData
	}

	R.mx.Lock()
	R.registry[filename] = rd
	R.mx.Unlock()

	return nil
}

func Render(filename string) (*image.NRGBA, error) {
	return nil, fmt.Errorf("not implemented")
}
