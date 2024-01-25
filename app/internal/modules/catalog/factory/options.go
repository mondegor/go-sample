package factory

import (
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
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
		RequestParsers  *RequestParsers
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

	RequestParsers struct {
		String *mrparser.String
		Image  *view_shared.ParserImage
		Parser *view_shared.Parser
	}
)
