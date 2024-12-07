package raster

import (
	"github.com/canghel3/raster2image/parser"
	"path/filepath"
)

type LoadOption func(driver Driver)

func WithStyle(style string) func(driver Driver) {
	return func(driver Driver) {
		switch filepath.Ext(filepath.Base(style)) {
		case ".css", "css":
			s, err := parser.NewCSSParser(style).Parse()
			if err == nil {
				driver.setStyle(s)
			}
		}
	}
}
