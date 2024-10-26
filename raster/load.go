package raster

import (
	"errors"
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
	registry map[string]*GodalDataset
}

type Data struct {
	style string
	min   float64
	max   float64
	ds    *godal.Dataset
}

type GodalDataset struct {
	data Data
	path string
}

func init() {
	R = &Registry{
		registry: make(map[string]*GodalDataset),
		mx:       sync.Mutex{},
	}
	godal.RegisterAll()
}

// Load opens the given raster file and stores it into the registry.
// Use Load only when opening the file for the first time, because loading is slow.
// For faster access, use Read afterward.
func Load(path string) (*GodalDataset, error) {
	ds, err := godal.Open(path)
	if err != nil {
		return nil, err
	}

	min, max, err := minMaxDs(ds)
	if err != nil {
		return nil, err
	}

	gd := &GodalDataset{
		data: Data{
			ds:    ds,
			min:   min,
			max:   max,
			style: "",
		},
		path: path,
	}

	R.mx.Lock()
	R.registry[filepath.Base(path)] = gd
	R.mx.Unlock()

	return gd, nil
}

// Read will retrieve the dataset quickly from the in-memory registry.
func Read(name string) (*GodalDataset, error) {
	gd, exists := R.registry[filepath.Base(name)]
	if exists {
		return gd, nil
	}

	return nil, errors.New("no such dataset exists. consider loading it first")
}

func Release(path string) {
	R.mx.Lock()
	registry, exists := R.registry[filepath.Base(path)]
	if exists {
		registry.Release()
	}
	delete(R.registry, filepath.Base(path))
	R.mx.Unlock()
}

func (gd *GodalDataset) Render(width, height int) (image.Image, error) {
	warped, err := gd.data.ds.Warp("", []string{
		"-of", "MEM",
		"-ts", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height),
	})
	if err != nil {
		return nil, err
	}
	defer warped.Close()

	switch len(warped.Bands()) {
	case 1:
		return gd.renderSingleBand(warped, width, height)
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 bands", gd.path)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 bands", gd.path)
	}

	return nil, nil
}

func (gd *GodalDataset) Copy() (*GodalDataset, error) {
	c, err := gd.data.ds.Translate("", []string{
		"-of", "MEM",
	})
	if err != nil {
		return nil, err
	}

	cpy := &GodalDataset{
		data: gd.data,
		path: gd.path,
	}

	cpy.data.ds = c
	return cpy, nil
}

func (gd *GodalDataset) renderSingleBand(warped *godal.Dataset, width, height int) (image.Image, error) {
	if gd.data.min == 0 && gd.data.max == 255 {
		//grayscale (or apply style)
		if len(gd.data.style) == 0 {
			band := warped.Bands()[0]
			var data = make([]float64, width*height)
			err := band.Read(0, 0, data, width, height)
			if err != nil {
				return nil, err
			}

			grayscale := render.Grayscale(data, width, height, gd.data.min, gd.data.max)
			return grayscale.Render()
		} else {
			//style given, so use rgb renderer with the style schema
			return nil, fmt.Errorf("not implemented")
		}
	} else {
		//requires normalization
		return nil, fmt.Errorf("cannot render raster %s with values larger than 255", gd.path)
	}
}

// Zoom essentially warps the dataset to the specified bbox extent.
// The underlying dataset is not modified.
func (gd *GodalDataset) Zoom(bbox [4]float64, srs string) (*GodalDataset, error) {
	options := []string{
		"-of", "MEM",
		"-te", fmt.Sprintf("%f", bbox[0]), fmt.Sprintf("%f", bbox[1]), fmt.Sprintf("%f", bbox[2]), fmt.Sprintf("%f", bbox[3]), // Set bounding box
		"-t_srs", srs, // target spatial reference system
		"-te_srs", "EPSG:3857",
	}

	warped, err := gd.data.ds.Warp("", options)
	if err != nil {
		return nil, err
	}

	newGd := &GodalDataset{
		data: gd.data,
		path: gd.path,
	}

	newGd.data.ds = warped

	return newGd, nil
}

func (gd *GodalDataset) Release() error {
	return gd.data.ds.Close()
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
