package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/models"
	"image"
	"path/filepath"
)

type GodalDataset struct {
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
	warped, err := gd.data.ds.Warp("", []string{
		"-of", "MEM",
		"-ts", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height),
	})
	if err != nil {
		return nil, err
	}
	defer warped.Close()

	cpy := gd.shallowCopy()
	cpy.data.ds = warped
	driver, err := cpy.newDriver()
	if err != nil {
		return nil, err
	}

	return driver.Render(width, height)
}

func (gd *GodalDataset) Copy() (*GodalDataset, error) {
	c, err := gd.data.ds.Translate("", []string{
		"-of", "MEM",
	})
	if err != nil {
		return nil, err
	}

	cpy := gd.shallowCopy()
	cpy.data.ds = c
	return cpy, nil
}

func (gd *GodalDataset) newDriver() (Driver, error) {
	ext := filepath.Ext(filepath.Base(gd.path))
	switch ext {
	case ".tif", "tif":
		return NewTifDriver(gd), nil
	}

	return nil, fmt.Errorf("no driver found for %s", ext)
}

// Zoom essentially warps the dataset to the specified bbox extent.
// The underlying dataset is not modified.
// It's recommended to defer Release on the returned dataset to avoid any resource leaks.
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
	return gd.data.ds.Close()
}
