package factory

import (
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		Locker          mrcore.Locker
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender

		CategoryAPI  catalog.CategoryAPI
		TrademarkAPI catalog.TrademarkAPI
		OrdererAPI   mrorderer.API

		UnitCategory *UnitCategoryOptions
	}

	UnitCategoryOptions struct {
		Dictionary      *mrlang.MultiLangDictionary
		ImageFileAPI    mrstorage.FileProviderAPI
		ImageURLBuilder mrcore.BuilderPath
	}
)
