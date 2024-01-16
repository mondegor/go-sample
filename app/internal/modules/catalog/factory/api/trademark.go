package factory_api

import (
	repository_api "go-sample/internal/modules/catalog/infrastructure/repository/api"
	usecase_api "go-sample/internal/modules/catalog/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtool"
)

func NewTrademark(conn *mrpostgres.ConnAdapter, serviceHelper *mrtool.ServiceHelper) *usecase_api.Trademark {
	return usecase_api.NewTrademark(
		repository_api.NewTrademarkPostgres(conn),
		serviceHelper,
	)
}
