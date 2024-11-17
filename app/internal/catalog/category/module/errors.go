package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

// ErrUseCaseCategoryImageNotFound - category image with ID not found.
var ErrUseCaseCategoryImageNotFound = mrerr.NewProto(
	"catalog.errCategoryImageNotFound", mrerr.ErrorKindUser, "category image with ID={{ .id }} not found")
