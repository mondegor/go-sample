package view_shared

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	RequestParserImage interface {
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserString
		mrserver.RequestParserImage
	}

	ParserImage struct {
		*mrparser.KeyInt32
		*mrparser.String
		*mrparser.Image
	}
)

func NewParserImage(
	p1 *mrparser.KeyInt32,
	p2 *mrparser.String,
	p3 *mrparser.Image,
) *ParserImage {
	return &ParserImage{
		KeyInt32: p1,
		String:   p2,
		Image:    p3,
	}
}
