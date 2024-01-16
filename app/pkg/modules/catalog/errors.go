package catalog

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrCategoryNotFound = NewFactory(
		"errCatalogCategoryNotFound", ErrorKindUser, "category with ID={{ .id }} not found")

	FactoryErrTrademarkNotFound = NewFactory(
		"errCatalogTrademarkNotFound", ErrorKindUser, "trademark with ID={{ .id }} not found")
)
