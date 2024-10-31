### This golang package is meant to help with converting raster files into images

- only supports single band .tif file with byte data (raster values are between 0-255) and .css styles
- 
```go
go get github.com/canghel3/raster2image
```

Usage:
```go
func main() {
	ds, err := raster.Load("/path/to/file", options.WithStyle("/path/to/cssStyle"))
	//once a dataset is loaded it can be read with Read much faster
	ds, err = raster.Read("file") // or use /path/to/file instead of file
	ds, err = ds.Zoom([]float64{0, 0, 0, 0}) // use a valid bbox extent
	image, err := ds.Render(256, 256)
	var buf bytes.Buffer
	png.Encode(&buf, image)
}
```