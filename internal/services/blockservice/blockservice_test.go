package blockservice

import (
	"testing"

	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	var rsl bool
	lblock := []schemas.Block{
		{
			Type:  schemas.BlockTypeScheme,
			Value: "ftp",
		},
		{
			Type:  schemas.BlockTypeScheme,
			Value: "rtmp",
		},
		{
			Type:  schemas.BlockTypeDomain,
			Value: "blockme",
		},
		{
			Type:  schemas.BlockTypeDomain,
			Value: "blockmetoo",
		},
		{
			Type:  schemas.BlockTypeRegex,
			Value: "(?m)regex - contain",
		},
	}

	s := &S{}
	//test invalid URL
	rsl = s.IsBlock("is not url", lblock)
	assert.True(t, rsl)

	//test block scheme
	rsl = s.IsBlock("ftp://test/block/ftp", lblock)
	assert.True(t, rsl)
	rsl = s.IsBlock("http://test/block/ftp", lblock)
	assert.False(t, rsl)
	rsl = s.IsBlock("rtmp://test/block/ftp", lblock)
	assert.True(t, rsl)

	//test block domain
	rsl = s.IsBlock("http://sholdnotblockme/test", lblock)
	assert.False(t, rsl)
	rsl = s.IsBlock("http://sholdnotblockmetoo/test", lblock)
	assert.False(t, rsl)
	rsl = s.IsBlock("http://blockme/test", lblock)
	assert.True(t, rsl)
	rsl = s.IsBlock("http://blockmetoo/test", lblock)
	assert.True(t, rsl)

	//tet block regex
	rsl = s.IsBlock("http://regex/port/true?search=regex - contain", lblock)
	assert.True(t, rsl)
	rsl = s.IsBlock("http://regex/port/true?search=regex%20-%20contain", lblock)
	assert.True(t, rsl)
	rsl = s.IsBlock("http://regex/port/true", lblock)
	assert.False(t, rsl)
}
