package raster

import (
	"bytes"
	"github.com/airbusgeo/godal"
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
	ds, err := Load(NasaInput)
	if err != nil {
		log.Fatal(err)
	}
	publicGodalDataset = ds

	os.Exit(m.Run())
}

func TestLoad(t *testing.T) {
	ds, err := Load(NasaInput)
	assert.NilError(t, err)
	assert.Check(t, ds != nil)
}

func TestRead(t *testing.T) {
	t.Run("ALREADY LOADED", func(t *testing.T) {
		ds, err := Load(NasaInput)
		assert.NilError(t, err)
		assert.Check(t, ds != nil)

		dsr, err := Read(NasaInput)
		assert.NilError(t, err)
		assert.Check(t, dsr != nil)
	})

	t.Run("NOT ALREADY LOADED", func(t *testing.T) {
		Release(NasaInput)
		ds, err := Read(NasaInput)
		assert.Error(t, err, "no such dataset exists. consider loading it first")
		assert.Assert(t, ds == nil)
	})
}

func TestRelease(t *testing.T) {
	ds, err := Load(NasaInput)
	assert.NilError(t, err)
	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	Release(NasaInput)

	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	ds, err = Read(NasaInput)
	assert.Error(t, err, "no such dataset exists. consider loading it first")
	assert.Check(t, ds == nil)

	n, err := Read(NasaInput)
	assert.Error(t, err, "no such dataset exists. consider loading it first")
	assert.Check(t, n == nil)
}

func TestZoom(t *testing.T) {
	ds, err := Load(NasaInput)
	assert.NilError(t, err)
	assert.Check(t, ds != nil)
	assert.Check(t, ds.data.ds != nil)

	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
	zoomed, err := ds.Zoom(bbox, "EPSG:3857")
	assert.NilError(t, err)

	assert.Check(t, ds.data != zoomed.data)
}

func TestRender(t *testing.T) {
	ds, err := Load(NasaInput)
	assert.NilError(t, err)
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

func TestCopy(t *testing.T) {
	ds, err := Load(NasaInput)
	assert.NilError(t, err)
	assert.Check(t, ds != nil)

	copied, err := ds.Copy()
	assert.NilError(t, err)
	assert.Check(t, copied != nil)

	assert.Check(t, ds != copied)
	//check that they point to different instances in memory
	assert.Check(t, ds.data.ds != copied.data.ds)
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkLoad
BenchmarkLoad-16    	       1	1622663965 ns/op ~ 1622ms/op
*/
func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds, err := Load(NasaInput)
		assert.NilError(b, err)
		assert.Check(b, ds != nil)
		assert.Check(b, ds.data.ds != nil)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkRead
BenchmarkRead-16    	 1645944	       749.9 ns/op
*/
func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds, err := Read(NasaInput)
		assert.NilError(b, err)
		assert.Check(b, ds != nil)
		assert.Check(b, ds.data.ds != nil)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkCopy
BenchmarkCopy-16    	       3	 394223565 ns/op ~ 394ms/op
*/
func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds, err := publicGodalDataset.Copy()
		assert.NilError(b, err)
		assert.Check(b, ds != nil)
		assert.Check(b, ds.data.ds != nil)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkZoom
BenchmarkZoom-16    	    5922	    174486 ns/op ~ 0.17ms/op
*/
func BenchmarkZoom(b *testing.B) {
	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
	for i := 0; i < b.N; i++ {
		ds, err := publicGodalDataset.Zoom(bbox, "EPSG:3857")
		assert.NilError(b, err)
		assert.Check(b, ds != nil)
		assert.Check(b, ds.data.ds != nil)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkZoomInto
BenchmarkZoomInto-16    	    6499	    181467 ns/op ~ 0.18ms/op
*/
func BenchmarkZoomInto(b *testing.B) {
	bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}

	for i := 0; i < b.N; i++ {
		var newDs = new(godal.Dataset)
		err := publicGodalDataset.ZoomInto(&newDs, bbox, "EPSG:3857")
		assert.NilError(b, err)
		assert.Check(b, newDs != nil)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/canghel3/raster2image/raster
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkRender
BenchmarkRender/ZOOM
BenchmarkRender/ZOOM-16         	     423	   2551499 ns/op ~ 2.55ms/op
BenchmarkRender/ZOOM_INTO
BenchmarkRender/ZOOM_INTO-16    	       2	 735253360 ns/op ~ 735ms/op
*/
func BenchmarkRender(b *testing.B) {
	b.Run("ZOOM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
			zoomed, err := publicGodalDataset.Zoom(bbox, "EPSG:3857")
			assert.NilError(b, err)

			render, err := zoomed.Render(256, 256)
			assert.NilError(b, err)

			var buf bytes.Buffer
			err = png.Encode(&buf, render)
			assert.NilError(b, err)

			// Prevent compiler optimizations by using the buffer's length
			_ = buf.Len()
		}
	})

	b.Run("ZOOM INTO", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var newDs = new(godal.Dataset)
			bbox := [4]float64{1364859.5770601074, 5119446.406427965, 1367305.561965233, 5121892.391333092}
			err := publicGodalDataset.ZoomInto(&newDs, bbox, "EPSG:3857")
			assert.NilError(b, err)

			render, err := publicGodalDataset.Render(256, 256)
			assert.NilError(b, err)

			var buf bytes.Buffer
			err = png.Encode(&buf, render)
			assert.NilError(b, err)

			// Prevent compiler optimizations by using the buffer's length
			_ = buf.Len()
		}
	})
}
