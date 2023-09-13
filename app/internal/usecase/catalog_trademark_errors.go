package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrCatalogTrademarkNotFound = NewFactory(
        "errCatalogTrademarkNotFound", ErrorKindUser, "trademark with ID={{ .id }} not found")
)
