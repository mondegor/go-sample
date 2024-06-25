package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

// ErrUseCaseCategoryImageNotFound - category image with ID not found.
var ErrUseCaseCategoryImageNotFound = mrerrfactory.NewProtoAppErrorByDefault(
	"errCatalogCategoryImageNotFound", mrerr.ErrorKindUser, "category image with ID={{ .id }} not found")

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrUseCaseCategoryImageNotFound,
	}
}
