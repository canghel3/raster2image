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
	xOff, yOff, xSize, ySize, err := td.getOffsetsAndSize(bbox)
	if err != nil {
		return nil, err
	}

	band := td.dataset.Bands()[0]
	var data = make([]float64, xSize*ySize)
	td.lock.RLock()
	err = band.Read(xOff, yOff, data, xSize, ySize)
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

func (td *TifDriver) getOffsetsAndSize(bbox [4]float64) (xOff, yOff, xSize, ySize int, err error) {
	gt, err := td.dataset.GeoTransform()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	minX, minY, maxX, maxY := bbox[0], bbox[1], bbox[2], bbox[3]

	xOffFloat := (minX - gt[0]) / gt[1]
	yOffFloat := (maxY - gt[3]) / gt[5]
	xEndFloat := (maxX - gt[0]) / gt[1]
	yEndFloat := (minY - gt[3]) / gt[5]

	xOff = int(math.Floor(xOffFloat))
	yOff = int(math.Floor(yOffFloat))
	xEnd := int(math.Ceil(xEndFloat))
	yEnd := int(math.Ceil(yEndFloat))

	xSize = xEnd - xOff
	ySize = yEnd - yOff

	// Clip to dataset boundaries
	dsWidth := td.dataset.Structure().SizeX
	dsHeight := td.dataset.Structure().SizeY

	if xOff < 0 {
		xSize += xOff
		xOff = 0
	}
	if yOff < 0 {
		ySize += yOff
		yOff = 0
	}
	if xOff+xSize > dsWidth {
		xSize = dsWidth - xOff
	}
	if yOff+ySize > dsHeight {
		ySize = dsHeight - yOff
	}

	if xSize <= 0 || ySize <= 0 {
		return 0, 0, 0, 0, fmt.Errorf("requested area is outside the raster extent")
	}

	return xOff, yOff, xSize, ySize, nil
}
