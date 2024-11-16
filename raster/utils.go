package raster

import (
	"github.com/airbusgeo/godal"
	"math/rand"
)

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
	min, max = 0, 0
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

const (
	minExtentX = 736241.4564
	minExtentY = 5826336.044
	maxExtentX = 2372605.358
	maxExtentY = 7049328.4966
	tileSize   = 256 // Tile size in pixels
	zoomLevel  = 7
)

func generateRandomBBoxWithinExtent() [4]float64 {
	// Calculate the tile size in map units (meters) at zoom level 7
	totalMapWidth := maxExtentX - minExtentX
	totalMapHeight := maxExtentY - minExtentY
	tileWidth := totalMapWidth / float64(1<<zoomLevel)   // Width of one tile in meters
	tileHeight := totalMapHeight / float64(1<<zoomLevel) // Height of one tile in meters

	// Determine the maximum number of tiles in the x and y ranges within the extent
	maxTileX := int(totalMapWidth / tileWidth)
	maxTileY := int(totalMapHeight / tileHeight)

	// Randomly select a tile index within the calculated range
	tileX := rand.Intn(maxTileX)
	tileY := rand.Intn(maxTileY)

	// Calculate the bounding box coordinates for the selected tile
	minX := minExtentX + float64(tileX)*tileWidth
	maxX := minX + tileWidth
	minY := maxExtentY - float64(tileY+1)*tileHeight
	maxY := minY + tileHeight

	return [4]float64{minX, minY, maxX, maxY}
}
