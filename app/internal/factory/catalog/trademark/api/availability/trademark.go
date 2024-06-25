package availability

import (
	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/repository"
	"github.com/mondegor/go-sample/internal/catalog/trademark/api/availability/usecase"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

// NewTrademark - comment func.
func NewTrademark(client mrstorage.DBConnManager, errorWrapper mrcore.UsecaseErrorWrapper) *usecase.Trademark {
	return usecase.NewTrademark(
		repository.NewTrademarkPostgres(client),
		errorWrapper,
	)
}
