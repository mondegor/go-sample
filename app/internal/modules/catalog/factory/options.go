package factory

import (
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type (
	Options struct {
		Logger          mrcore.Logger
		EventBox        mrcore.EventBox
		ServiceHelper   *mrtool.ServiceHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		Locker          mrcore.Locker
		OrdererAPI      mrorderer.API

		CategoryAPI  catalog.CategoryAPI
		TrademarkAPI catalog.TrademarkAPI

		UnitCategory *UnitCategoryOptions
	}

	UnitCategoryOptions struct {
		Dictionary      *mrlang.MultiLangDictionary
		ImageFileAPI    mrstorage.FileProviderAPI
		ImageURLBuilder mrcore.BuilderPath
	}
)
