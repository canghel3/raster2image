package raster

import (
	"errors"
	"github.com/airbusgeo/godal"
	"path/filepath"
	"sync"
)

var R *Registry

type Registry struct {
	mx       sync.RWMutex
	registry map[string]*GodalDataset
}

func init() {
	R = &Registry{
		registry: make(map[string]*GodalDataset),
		mx:       sync.RWMutex{},
	}
	godal.RegisterAll()
}

// Load opens the given raster file and stores it into the registry.
// Use Load only when opening the file for the first time, because loading is slow.
// For faster access, use Read afterward.
func Load(path string, options ...LoadOption) (*GodalDataset, error) {
	ds, err := godal.Open(path)
	if err != nil {
		return nil, err
	}

	min, max, err := minMaxDs(ds)
	if err != nil {
		return nil, err
	}

	gd := GodalDataset{
		data: Data{
			ds:  ds,
			min: min,
			max: max,
		},
		path: path,
	}

	for _, option := range options {
		option(&gd)
	}

	R.mx.Lock()
	R.registry[filepath.Base(path)] = &gd
	R.mx.Unlock()

	return &gd, nil
}

// Read will retrieve the dataset quickly from the in-memory registry.
func Read(name string) (*GodalDataset, error) {
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
