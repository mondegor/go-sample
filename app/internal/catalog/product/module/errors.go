package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrUseCaseProductNotFound - product with ID not found.
	ErrUseCaseProductNotFound = mrerr.NewProto(
		"catalog.errProductNotFound", mrerr.ErrorKindUser, "product with ID={{ .id }} not found")

	// ErrUseCaseProductArticleAlreadyExists - product article already exists.
	ErrUseCaseProductArticleAlreadyExists = mrerr.NewProto(
		"catalog.errProductArticleAlreadyExists", mrerr.ErrorKindUser, "product article '{{ .name }}' already exists")
)
