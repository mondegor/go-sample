package factory_api

import (
	repository_api "go-sample/internal/modules/catalog/infrastructure/repository/api"
	usecase_api "go-sample/internal/modules/catalog/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrtool"
)

func NewCategory(conn *mrpostgres.ConnAdapter, serviceHelper *mrtool.ServiceHelper) *usecase_api.Category {
	return usecase_api.NewCategory(
		repository_api.NewCategoryPostgres(conn),
		serviceHelper,
	)
}
