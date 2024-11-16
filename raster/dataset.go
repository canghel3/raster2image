package raster

import (
	"errors"
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/models"
	"image"
	"path/filepath"
	"sync"
)

type GodalDataset struct {
	//TODO: if lock contention is high, use a renderer worker pool or find a scalable alternative
	rw   sync.RWMutex
	data Data
	path string
}

type Data struct {
	min   float64
	max   float64
	style *models.RasterStyle
	ds    *godal.Dataset
}

func (gd *GodalDataset) Render(width, height uint) (image.Image, error) {
	gd.rw.Lock()
	if gd.data.ds == nil {
		return nil, errors.New("godal dataset is nil")
	}
	warped, err := gd.data.ds.Warp("", []string{
		"-of", "MEM",
		"-ts", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height),
	})
	gd.rw.Unlock()
	if err != nil {
		return nil, err
	}

	cpy := gd.shallowCopy()
	cpy.data.ds = warped
	driver, err := cpy.newRenderDriver()
	if err != nil {
		return nil, err
	}

	return driver.Render(width, height)
}

func (gd *GodalDataset) Copy() (*GodalDataset, error) {
	gd.rw.Lock()
	if gd.data.ds == nil {
		return nil, errors.New("godal dataset is nil")
	}
	c, err := gd.data.ds.Translate("", []string{
		"-of", "MEM",
	})
	gd.rw.Unlock()
	if err != nil {
		return nil, err
	}

	cpy := gd.shallowCopy()
	cpy.data.ds = c
	return cpy, nil
}

func (gd *GodalDataset) newRenderDriver() (RenderDriver, error) {
	ext := filepath.Ext(filepath.Base(gd.path))
	switch ext {
	case ".tif", "tif":
		return NewTifDriver(gd.data.ds.Bands(), gd.path, gd.data.min, gd.data.max, gd.data.style), nil
	}

	return nil, fmt.Errorf("no driver found for %s", ext)
}

// Zoom essentially warps the dataset to the specified bbox extent (minX,minY,maxX,maxY format).
// The underlying dataset is not modified.
// It's recommended to defer Release on the returned dataset to avoid any resource leaks.
func (gd *GodalDataset) Zoom(bbox [4]float64, srs string) (*GodalDataset, error) {
	options := []string{
		"-of", "MEM",
		"-te", fmt.Sprintf("%f", bbox[0]), fmt.Sprintf("%f", bbox[1]), fmt.Sprintf("%f", bbox[2]), fmt.Sprintf("%f", bbox[3]), // Set bounding box
		"-t_srs", srs, // target spatial reference system
		"-te_srs", "EPSG:3857",
	}

	gd.rw.Lock()
	if gd.data.ds == nil {
		return nil, errors.New("godal dataset is nil")
	}
	warped, err := gd.data.ds.Warp("", options)
	gd.rw.Unlock()
	if err != nil {
		return nil, err
	}

	newGd := gd.shallowCopy()
	newGd.data.ds = warped

	return newGd, nil
}

func (gd *GodalDataset) shallowCopy() *GodalDataset {
	newGd := &GodalDataset{
		data: gd.data,
		path: gd.path,
	}

	return newGd
}

func (gd *GodalDataset) Release() error {
	gd.rw.Lock()
	defer gd.rw.Unlock()

	return gd.data.ds.Close()
}
