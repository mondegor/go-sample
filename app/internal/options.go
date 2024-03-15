package app

import (
	"context"
	"go-sample/config"
	factory_catalog_category "go-sample/internal/modules/catalog/category/factory"
	factory_catalog_product "go-sample/internal/modules/catalog/product/factory"
	factory_catalog_trademark "go-sample/internal/modules/catalog/trademark/factory"
	factory_filestation "go-sample/internal/modules/file-station/factory"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		Cfg              config.Config
		EventEmitter     mrsender.EventEmitter
		UsecaseHelper    *mrcore.UsecaseHelper
		PostgresAdapter  *mrpostgres.ConnAdapter
		RedisAdapter     *mrredis.ConnAdapter
		FileProviderPool *mrstorage.FileProviderPool
		Locker           mrlock.Locker
		Translator       *mrlang.Translator
		RequestParsers   RequestParsers
		ResponseSender   *mrresponse.Sender
		AccessControl    mrperms.AccessControl
		ImageURLBuilder  mrlib.BuilderPath

		// API section
		CatalogCategoryAPI  catalog.CategoryAPI
		CatalogTrademarkAPI catalog.TrademarkAPI
		OrdererAPI          mrorderer.API

		// Modules section
		CatalogCategoryModule  factory_catalog_category.Options
		CatalogProductModule   factory_catalog_product.Options
		CatalogTrademarkModule factory_catalog_trademark.Options
		FileStationModule      factory_filestation.Options

		OpenedResources []func(ctx context.Context)
	}

	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		KeyInt32   *mrparser.KeyInt32
		ListSorter *mrparser.ListSorter
		ListPager  *mrparser.ListPager
		String     *mrparser.String
		// UUID       *mrparser.UUID
		Validator *mrparser.Validator
		// File       *mrparser.File
		Image *mrparser.Image
	}
)
