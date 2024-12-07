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

	// If the requested output size differs from xSize, ySize, we must resample
	// If width == xSize && height == ySize, no resampling needed
	//finalWidth := int(width)
	//finalHeight := int(height)
	//var dataToDraw []float64
	//if finalWidth != xSize || finalHeight != ySize {
	//	log.Println("resampling")
	//	dataToDraw = nearestResample(data, xSize, ySize, finalWidth, finalHeight)
	//} else {
	//	log.Println("NOT RESAMPLING")
	//	dataToDraw = data
	//}

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

	// Convert from georeferenced to pixel space
	xOffFloat := (minX - gt[0]) / gt[1]
	yOffFloat := (maxY - gt[3]) / gt[5] // Note: gt[5] is usually negative
	xEndFloat := (maxX - gt[0]) / gt[1]
	yEndFloat := (minY - gt[3]) / gt[5]

	// Use floor/ceil to get pixel boundaries
	xOff = int(math.Floor(xOffFloat))
	yOff = int(math.Floor(yOffFloat))
	xEnd := int(math.Ceil(xEndFloat))
	yEnd := int(math.Ceil(yEndFloat))

	xSize = xEnd - xOff
	ySize = yEnd - yOff

	dsWidth := td.dataset.Structure().SizeX
	dsHeight := td.dataset.Structure().SizeY

	// Clamp values to dataset boundaries
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

func nearestResample(src []float64, srcWidth, srcHeight, dstWidth, dstHeight int) []float64 {
	if dstWidth == srcWidth && dstHeight == srcHeight {
		// No resampling needed, just return a copy
		out := make([]float64, len(src))
		copy(out, src)
		return out
	}

	out := make([]float64, dstWidth*dstHeight)
	// Compute ratios
	xRatio := float64(srcWidth) / float64(dstWidth)
	yRatio := float64(srcHeight) / float64(dstHeight)

	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// Map the output pixel back to source coordinates
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)

			// Clamp to avoid any floating rounding issues (shouldn't normally happen)
			if srcX >= srcWidth {
				srcX = srcWidth - 1
			}
			if srcY >= srcHeight {
				srcY = srcHeight - 1
			}

			out[y*dstWidth+x] = src[srcY*srcWidth+srcX]
		}
	}

	return out
}

func bilinearResample(src []float64, srcWidth, srcHeight, dstWidth, dstHeight int) []float64 {
	out := make([]float64, dstWidth*dstHeight)
	xRatio := float64(srcWidth-1) / float64(dstWidth-1)
	yRatio := float64(srcHeight-1) / float64(dstHeight-1)

	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// Map dst pixel to fractional src coordinates
			srcX := float64(x) * xRatio
			srcY := float64(y) * yRatio

			xFloor := int(math.Floor(srcX))
			yFloor := int(math.Floor(srcY))
			xCeil := xFloor + 1
			yCeil := yFloor + 1

			if xCeil >= srcWidth {
				xCeil = srcWidth - 1
			}
			if yCeil >= srcHeight {
				yCeil = srcHeight - 1
			}

			topLeft := src[yFloor*srcWidth+xFloor]
			topRight := src[yFloor*srcWidth+xCeil]
			bottomLeft := src[yCeil*srcWidth+xFloor]
			bottomRight := src[yCeil*srcWidth+xCeil]

			xFrac := srcX - float64(xFloor)
			yFrac := srcY - float64(yFloor)

			// Interpolate in X
			topVal := topLeft + (topRight-topLeft)*xFrac
			bottomVal := bottomLeft + (bottomRight-bottomLeft)*xFrac

			// Interpolate in Y
			out[y*dstWidth+x] = topVal + (bottomVal-topVal)*yFrac
		}
	}

	return out
}
