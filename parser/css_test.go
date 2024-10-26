package parser

import (
	"gotest.tools/v3/assert"
	"testing"
)

const SampleCss = "./testdata/styles/sample.css"

func TestCSSParser(t *testing.T) {
	parser := NewCSSParser(SampleCss)
	style, err := parser.Parse()
	assert.NilError(t, err)

	assert.Assert(t, style != nil)
	assert.Assert(t, style.RasterChannels == "auto")
	assert.Assert(t, len(style.ColorMap) == 11)
}
