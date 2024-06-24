package factory_api

import (
	repository_api "go-sample/internal/modules/catalog/trademark/infrastructure/repository/api"
	usecase_api "go-sample/internal/modules/catalog/trademark/usecase/api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewTrademark - comment func.
func NewTrademark(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase_api.Trademark {
	return usecase_api.NewTrademark(
		repository_api.NewTrademarkPostgres(client),
		errorWrapper,
	)
}
