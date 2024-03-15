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
		mrserver.RequestParserListSorter
		mrserver.RequestParserListPager
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
		*mrparser.ListSorter
		*mrparser.ListPager
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
	p5 *mrparser.ListSorter,
	p6 *mrparser.ListPager,
	p7 *mrparser.String,
	// p8 *mrparser.UUID,
	p9 *mrparser.Validator,
	// p10 *mrparser.File,
	p11 *mrparser.Image,
	p12 *mrparser.ItemStatus,
) *Parser {
	return &Parser{
		// Bool:      p1,
		// DateTime:  p2,
		Int64:      p3,
		KeyInt32:   p4,
		ListSorter: p5,
		ListPager:  p6,
		String:     p7,
		// UUID:      p8,
		Validator: p9,
		// File:      p10,
		Image:      p11,
		ItemStatus: p12,
	}
}
