package raster

import (
	"fmt"
	"github.com/airbusgeo/godal"
	"github.com/canghel3/raster2image/models"
	"github.com/canghel3/raster2image/render"
	"image"
	"math"
	"sync"
)

type TifDriver struct {
	name    string
	lock    sync.RWMutex
	dataset *godal.Dataset
	min     float64
	max     float64
	style   *models.RasterStyle
}

type TifDriverData struct {
	Name    string
	Dataset *godal.Dataset
	Min     float64
	Max     float64
	Style   *models.RasterStyle
}

func NewTifDriver(data TifDriverData) Driver {
	return &TifDriver{
		name:    data.Name,
		dataset: data.Dataset,
		max:     data.Max,
		min:     data.Min,
		style:   data.Style,
	}
}

func (td *TifDriver) Render(bbox [4]float64, width, height uint) (image.Image, error) {
	switch len(td.dataset.Bands()) {
	case 1:
		return td.renderSingleBand(bbox, width, height)
	case 2:
		return nil, fmt.Errorf("cannot render raster %s with 2 Bands", td.name)
	case 3:
		//rgb
	case 4:
		return nil, fmt.Errorf("cannot render raster %s with 4 Bands", td.name)
	}

	return nil, nil
}

func (td *TifDriver) renderSingleBand(bbox [4]float64, width, height uint) (image.Image, error) {
	xOffset, yOffset, err := td.getOffset(bbox)
	if err != nil {
		return nil, err
	}

	band := td.dataset.Bands()[0]
	var data = make([]float64, width*height)
	td.lock.RLock()
	err = band.Read(xOffset, yOffset, data, int(width), int(height))
	td.lock.RUnlock()
	if err != nil {
		return nil, err
	}

	if td.style != nil {
		//setStyle given, so use rgb renderer with the setStyle schema
		rgb := render.NewRGBDrawer(data, int(width), int(height), render.StyleOption(*td.style))
		return rgb.Draw()
	}

	grayscale := render.Grayscale(data, int(width), int(height), td.min, td.max)
	return grayscale.Draw()
}

func (td *TifDriver) Release() error {
	td.lock.Lock()
	defer td.lock.Unlock()
	return td.dataset.Close()
}

func (td *TifDriver) setStyle(style *models.RasterStyle) {
	td.style = style
}

func (td *TifDriver) getOffset(bbox [4]float64) (x, y int, err error) {
	gt, err := td.dataset.GeoTransform()
	if err != nil {
		return 0, 0, err
	}

	//bbox (minX, minY, maxX, maxY)
	minX, _, _, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	//calculate pixel offsets and sizes
	xOffset := int(math.Floor((minX - gt[0]) / gt[1]))
	yOffset := int(math.Floor((maxY - gt[3]) / gt[5]))
	//xEnd := int(math.Ceil((maxX - gt[0]) / gt[1]))
	//yEnd := int(math.Ceil((minY - gt[3]) / gt[5]))

	//width := xEnd - xOffset
	//height := yEnd - yOffset

	return xOffset, yOffset, nil
}
