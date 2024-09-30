package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-sample/internal/catalog/category/api/availability/repository"
	"github.com/mondegor/go-sample/internal/catalog/category/api/availability/usecase"
)

// NewCategory - создаёт объект usecase.Category.
func NewCategory(client mrstorage.DBConnManager, errorWrapper mrcore.UseCaseErrorWrapper) *usecase.Category {
	return usecase.NewCategory(
		repository.NewCategoryPostgres(client),
		errorWrapper,
	)
}
