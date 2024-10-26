package raster

type LoadOption func(*GodalDataset)

func WithStyle(style string) LoadOption {
	return func(g *GodalDataset) {
		g.data.style = style
	}
}
