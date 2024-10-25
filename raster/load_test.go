package raster

import (
	"bytes"
	"gotest.tools/v3/assert"
	"image/png"
	"log"
	"os"
	"testing"
)

const (
	NasaInput = "./testdata/nasa_lights_compr.tif"
	NasaPng   = "./testdata/generated/nasa_lights_compr.png"
)

var (
	publicGodalDataset *GodalDataset
)

func TestMain(m *testing.M) {
	err := Open(NasaInput)
	if err != nil {
		log.Fatal(err)
	}

	ds := Load(NasaInput)
	publicGodalDataset = ds

	os.Exit(m.Run())
}

func TestOpen(t *testing.T) {
	err := Open(NasaInput)
	assert.NilError(t, err)
}

func TestLoad(t *testing.T) {
	err := Open(NasaInput)
	assert.NilError(t, err)

	ds := Load(NasaInput)
	assert.Check(t, ds != nil)
}

func TestRelease(t *testing.T) {
	err := Open(NasaInput)
	assert.NilError(t, err)

	ds := Load(NasaInput)
	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	Release(NasaInput)

	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	ds = Load(NasaInput)
	assert.Check(t, ds == nil)

	n := Load(NasaInput)
	assert.Check(t, n == nil)
}

func TestZoom(t *testing.T) {
	err := Open(NasaInput)
	assert.NilError(t, err)

	ds := Load(NasaInput)
	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
	zoomed, err := ds.Zoom(bbox, "EPSG:3857")
	assert.NilError(t, err)

	assert.Check(t, ds.data != zoomed.data)
}

func TestRender(t *testing.T) {
	err := Open(NasaInput)
	assert.NilError(t, err)

	ds := Load(NasaInput)
	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
	zoomed, err := ds.Zoom(bbox, "EPSG:3857")
	assert.NilError(t, err)

	render, err := zoomed.Render(256, 256)
	assert.NilError(t, err)

	var buf bytes.Buffer
	err = png.Encode(&buf, render)
	assert.NilError(t, err)

	err = os.WriteFile(NasaPng, buf.Bytes(), 0666)
	assert.NilError(t, err)
}

func BenchmarkRender(b *testing.B) {
	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
	zoomed, err := publicGodalDataset.Zoom(bbox, "EPSG:3857")
	assert.NilError(b, err)

	render, err := zoomed.Render(256, 256)
	assert.NilError(b, err)

	var buf bytes.Buffer
	err = png.Encode(&buf, render)
	assert.NilError(b, err)
}
