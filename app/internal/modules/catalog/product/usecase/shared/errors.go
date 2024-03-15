package usecase_shared

import (
	. "github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrProductNotFound = NewFactory(
		"errCatalogProductNotFound", ErrorKindUser, "product with ID={{ .id }} not found")

	FactoryErrProductArticleAlreadyExists = NewFactory(
		"errCatalogProductArticleAlreadyExists", ErrorKindUser, "product article '{{ .name }}' already exists")
)
