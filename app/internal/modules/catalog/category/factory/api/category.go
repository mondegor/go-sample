package factory_api

import (
	repository_api "go-sample/internal/modules/catalog/category/infrastructure/repository/api"
	usecase_api "go-sample/internal/modules/catalog/category/usecase/api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewCategory - comment func.
func NewCategory(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase_api.Category {
	return usecase_api.NewCategory(
		repository_api.NewCategoryPostgres(client),
		errorWrapper,
	)
}
