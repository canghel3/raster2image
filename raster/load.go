package raster

import (
	"errors"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/utils"
	"path/filepath"
	"sync"
)

var R *Registry

type Registry struct {
	mx       sync.RWMutex
	registry map[string]Driver
}

func init() {
	R = &Registry{
		registry: make(map[string]Driver),
		mx:       sync.RWMutex{},
	}
	godal.RegisterAll()
}

// Load opens the given raster file and stores it into the registry.
// Use Load only when opening the file for the first time, because loading is slow.
// For faster access, use Read afterward.
func Load(path string, options ...LoadOption) (Driver, error) {
	ds, err := godal.Open(path)
	if err != nil {
		return nil, err
	}

	//TODO: run gdaladdo to create internal pyramids to improve efficiency

	min, max, err := utils.MinMaxDs(ds)
	if err != nil {
		return nil, err
	}

	var driver Driver
	switch filepath.Ext(path) {
	case ".tif":
		tifDriverData := TifDriverData{
			Name:    path,
			Dataset: ds,
			Min:     min,
			Max:     max,
		}

		driver = NewTifDriver(tifDriverData)
	default:
		return nil, errors.New("file type not supported")
	}

	for _, option := range options {
		option(driver)
	}

	R.mx.Lock()
	R.registry[filepath.Base(path)] = driver
	R.mx.Unlock()

	return driver, nil
}

// Read will retrieve the dataset quickly from the in-memory registry.
func Read(name string) (Driver, error) {
	R.mx.RLock()
	gd, exists := R.registry[filepath.Base(name)]
	R.mx.RUnlock()
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
