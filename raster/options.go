package raster

import (
	"github.com/canghel3/raster2image/parser"
	"path/filepath"
)

type LoadOption func(*GodalDataset)

func WithStyle(style string) LoadOption {
	return func(g *GodalDataset) {
		switch filepath.Ext(filepath.Base(style)) {
		case ".css", "css":
			s, err := parser.NewCSSParser(style).Parse()
			if err == nil {
				g.data.style = s
			}
		}
		return
	}
}
