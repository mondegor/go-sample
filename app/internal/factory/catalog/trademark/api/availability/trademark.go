package availability

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/repository"
	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/usecase"
)

// NewTrademark - создаёт объект usecase.Trademark.
func NewTrademark(client mrstorage.DBConnManager, errorWrapper mrcore.UseCaseErrorWrapper) *usecase.Trademark {
	return usecase.NewTrademark(
		repository.NewTrademarkPostgres(client),
		errorWrapper,
	)
}
