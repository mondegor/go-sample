package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mrerrfactory"
)

var (
	// ErrUseCaseProductNotFound - product with ID not found.
	ErrUseCaseProductNotFound = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogProductNotFound", mrerr.ErrorKindUser, "product with ID={{ .id }} not found")

	// ErrUseCaseProductArticleAlreadyExists - product article already exists.
	ErrUseCaseProductArticleAlreadyExists = mrerrfactory.NewProtoAppErrorByDefault(
		"errCatalogProductArticleAlreadyExists", mrerr.ErrorKindUser, "product article '{{ .name }}' already exists")
)

// Errors - comment func.
func Errors() []*mrerr.ProtoAppError {
	return []*mrerr.ProtoAppError{
		ErrUseCaseProductNotFound,
		ErrUseCaseProductArticleAlreadyExists,
	}
}
