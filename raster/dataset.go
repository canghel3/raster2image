package raster

//type GodalDataset struct {
//	//TODO: if lock contention is high find a scalable alternative
//	//each dataset will have its own lock. zooming and rendering overwrites the godal dataset so using a different lock is fine.
//	rw     sync.RWMutex
//	data   Data
//	driver Driver
//	path   string
//}
//
//type Data struct {
//	min   float64
//	max   float64
//	style *models.RasterStyle
//	ds    *godal.Dataset
//}
//
//// TODO: refactor and calculate read dataset with offset instead of warping to the given extent
//func (gd *GodalDataset) Render(width, height int) (image.Image, error) {
//	gd.rw.Lock()
//	if gd.data.ds == nil {
//		return nil, errors.New("godal dataset is nil")
//	}
//	//TODO: fmt.Sprintf is slow, use a different approach
//	warped, err := gd.data.ds.Warp("", []string{
//		"-of", "MEM",
//		"-ts", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height),
//		"-r", "bilinear",
//	})
//	gd.rw.Unlock()
//	if err != nil {
//		return nil, err
//	}
//
//	return gd.driver.Render(warped.Bands(), uint(width), uint(height))
//}
//
//func (gd *GodalDataset) Copy() (*GodalDataset, error) {
//	gd.rw.Lock()
//	if gd.data.ds == nil {
//		return nil, errors.New("godal dataset is nil")
//	}
//	c, err := gd.data.ds.Translate("", []string{
//		"-of", "MEM",
//	})
//	gd.rw.Unlock()
//	if err != nil {
//		return nil, err
//	}
//	gd.data.ds.GeoTransform()
//
//	cpy := gd.shallowCopy()
//	cpy.data.ds = c
//	return cpy, nil
//}
//
//// Zoom warps the dataset to the specified bbox extent (minX,minY,maxX,maxY format).
//// The underlying dataset is not modified.
//// It's recommended to defer Release on the returned dataset to avoid any resource leaks.
//func (gd *GodalDataset) Zoom(bbox [4]float64, srs string) (*GodalDataset, error) {
//	options := []string{
//		"-of", "MEM",
//		"-te", fmt.Sprintf("%f", bbox[0]), fmt.Sprintf("%f", bbox[1]), fmt.Sprintf("%f", bbox[2]), fmt.Sprintf("%f", bbox[3]), // Set bounding box
//		"-s_srs", srs,
//		"-t_srs", srs, // target spatial reference system
//		"-te_srs", "EPSG:3857",
//		"-r", "bilinear",
//	}
//
//	gd.rw.Lock()
//	if gd.data.ds == nil {
//		return nil, errors.New("godal dataset is nil")
//	}
//	warped, err := gd.data.ds.Warp("", options)
//	gd.rw.Unlock()
//	if err != nil {
//		return nil, err
//	}
//
//	newGd := gd.shallowCopy()
//	newGd.data.ds = warped
//
//	return newGd, nil
//}
//
//func (gd *GodalDataset) shallowCopy() *GodalDataset {
//	newGd := &GodalDataset{
//		rw:     sync.RWMutex{},
//		data:   gd.data,
//		driver: gd.driver,
//		path:   gd.path,
//	}
//
//	return newGd
//}
//
//func (gd *GodalDataset) Release() error {
//	gd.rw.Lock()
//	defer gd.rw.Unlock()
//
//	return gd.data.ds.Close()
//}
