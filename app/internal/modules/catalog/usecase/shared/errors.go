package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrCategoryNotFound = NewFactory(
		"errCategoryNotFound", ErrorKindUser, "category with ID={{ .id }} not found")

	FactoryErrProductArticleAlreadyExists = NewFactory(
		"errProductArticleAlreadyExists", ErrorKindUser, "product article '{{ .name }}' already exist")

	FactoryErrTrademarkNotFound = NewFactory(
		"errTrademarkNotFound", ErrorKindUser, "trademark with ID={{ .id }} not found")
)
