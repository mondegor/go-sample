package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrCategoryNotFound = NewFactory(
		"errCatalogCategoryNotFound", ErrorKindUser, "category with ID={{ .id }} not found")

	FactoryErrCategoryImageNotFound = NewFactory(
		"errCatalogCategoryImageNotFound", ErrorKindUser, "category image with ID={{ .id }} not found")

	FactoryErrTrademarkNotFound = NewFactory(
		"errCatalogTrademarkNotFound", ErrorKindUser, "trademark with ID={{ .id }} not found")

	FactoryErrProductNotFound = NewFactory(
		"errCatalogProductNotFound", ErrorKindUser, "product with ID={{ .id }} not found")

	FactoryErrProductArticleAlreadyExists = NewFactory(
		"errCatalogProductArticleAlreadyExists", ErrorKindUser, "product article '{{ .name }}' already exist")
)
