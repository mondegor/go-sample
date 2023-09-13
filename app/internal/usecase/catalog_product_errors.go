package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    FactoryErrCatalogProductArticleAlreadyExists = NewFactory(
        "errCatalogProductArticleAlreadyExists", ErrorKindUser, "product article '{{ .name }}' is already exists")
)
