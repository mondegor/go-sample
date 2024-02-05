package usecase_shared

import (
	. "github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrCategoryImageNotFound = NewFactory(
		"errCatalogCategoryImageNotFound", ErrorKindUser, "category image with ID={{ .id }} not found")
)
