package factory_api

import (
	repository_api "go-sample/internal/modules/catalog/infrastructure/repository/api"
	usecase_api "go-sample/internal/modules/catalog/usecase/api"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewCategory(conn *mrpostgres.ConnAdapter, usecaseHelper *mrcore.UsecaseHelper) *usecase_api.Category {
	return usecase_api.NewCategory(
		repository_api.NewCategoryPostgres(conn),
		usecaseHelper,
	)
}
