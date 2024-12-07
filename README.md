### This golang package is meant to help with converting raster files into images

- for now only supports single band .tif file with byte data (raster values are between 0-255) and .css styles

```go
go get github.com/canghel3/raster2image
```

Usage:
```go
func main() {
	driver, err := raster.Load("/path/to/file.tif", options.WithStyle("/path/to/style.css"))
	//once a dataset is loaded it can be read with Read much faster
	driver, err = raster.Read("file.tif") // or use /path/to/file.tif instead of file.tif, works either way
	bbox := [4]float64{0, 0, 0, 0} //use a valid bbox
	image, err := driver.Render(bbox, 256, 256)
	var buf bytes.Buffer
	png.Encode(&buf, image)
}
```