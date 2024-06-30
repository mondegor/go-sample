package product

import (
	"github.com/mondegor/go-components/mrsort"

	"github.com/mondegor/go-sample/pkg/validate"

	"github.com/mondegor/go-sample/pkg/catalog/api"

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
		RequestParsers RequestParsers
		ResponseSender mrserver.ResponseSender

		CategoryAPI  api.CategoryAvailability
		TrademarkAPI api.TrademarkAvailability
		OrdererAPI   mrsort.Orderer

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}
)
