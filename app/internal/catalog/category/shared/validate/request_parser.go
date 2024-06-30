package validate

import (
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"

	pkgvalidate "github.com/mondegor/go-sample/pkg/validate"
)

type (
	// RequestCategoryParser - comment interface.
	RequestCategoryParser interface {
		pkgvalidate.RequestExtendParser
		mrserver.RequestParserImage
	}

	// CategoryParser - comment struct.
	CategoryParser struct {
		*pkgvalidate.ExtendParser
		*mrparser.Image
	}
)

// NewCategoryParser - создаёт объект CategoryParser.
func NewCategoryParser(
	p1 *pkgvalidate.ExtendParser,
	p2 *mrparser.Image,
) *CategoryParser {
	return &CategoryParser{
		ExtendParser: p1,
		Image:        p2,
	}
}
