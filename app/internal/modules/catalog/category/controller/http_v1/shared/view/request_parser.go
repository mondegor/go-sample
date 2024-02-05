package view_shared

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
)

type (
	RequestParser interface {
		// mrserver.RequestParserBool
		// mrserver.RequestParserDateTime
		mrserver.RequestParserInt64
		mrserver.RequestParserKeyInt32
		mrserver.RequestParserSortPage
		mrserver.RequestParserString
		// mrserver.RequestParserUUID
		mrserver.RequestParserValidate
		// mrserver.RequestParserFile
		mrserver.RequestParserImage
		mrserver.RequestParserItemStatus
	}

	Parser struct {
		// *mrparser.Bool
		// *mrparser.DateTime
		*mrparser.Int64
		*mrparser.KeyInt32
		*mrparser.SortPage
		*mrparser.String
		// *mrparser.UUID
		*mrparser.Validator
		// *mrparser.File
		*mrparser.Image
		*mrparser.ItemStatus
	}
)

func NewParser(
	// p1 *mrparser.Bool,
	// p2 *mrparser.DateTime,
	p3 *mrparser.Int64,
	p4 *mrparser.KeyInt32,
	p5 *mrparser.SortPage,
	p6 *mrparser.String,
	// p7 *mrparser.UUID,
	p8 *mrparser.Validator,
	// p9 *mrparser.File,
	p10 *mrparser.Image,
	p11 *mrparser.ItemStatus,
) *Parser {
	return &Parser{
		// Bool:      p1,
		// DateTime:  p2,
		Int64:    p3,
		KeyInt32: p4,
		SortPage: p5,
		String:   p6,
		// UUID:      p7,
		Validator: p8,
		// File:      p9,
		Image:      p10,
		ItemStatus: p11,
	}
}
