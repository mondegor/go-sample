package factory

import (
	view_shared "go-sample/internal/modules/catalog/product/controller/httpv1/shared/view"
	"go-sample/pkg/modules/catalog/api"

	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		RequestParser  *view_shared.Parser
		ResponseSender mrserver.ResponseSender

		CategoryAPI  api.CategoryAPI
		TrademarkAPI api.TrademarkAPI
		OrdererAPI   mrsort.Orderer

		PageSizeMax     uint64
		PageSizeDefault uint64
	}
)
