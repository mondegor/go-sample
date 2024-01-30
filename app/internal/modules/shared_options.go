package modules

import (
	"context"
	"go-sample/config"
	factory_catalog "go-sample/internal/modules/catalog/factory"
	factory_filestation "go-sample/internal/modules/file-station/factory"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
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

		// API section
		CatalogCategoryAPI  catalog.CategoryAPI
		CatalogTrademarkAPI catalog.TrademarkAPI
		OrdererAPI          mrorderer.API

		// Modules section
		CatalogModule     factory_catalog.Options
		FileStationModule factory_filestation.Options

		OpenedResources []func(ctx context.Context)
	}

	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		KeyInt32   *mrparser.KeyInt32
		SortPage   *mrparser.SortPage
		String     *mrparser.String
		// UUID       *mrparser.UUID
		Validator *mrparser.Validator
		// File       *mrparser.Image
		Image *mrparser.Image
	}
)
