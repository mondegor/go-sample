package app

import (
	"context"

	"go-sample/config"
	factory_catalog_category "go-sample/internal/modules/catalog/category/factory"
	factory_catalog_product "go-sample/internal/modules/catalog/product/factory"
	factory_catalog_trademark "go-sample/internal/modules/catalog/trademark/factory"
	factory_filestation "go-sample/internal/modules/filestation/factory"
	"go-sample/pkg/modules/catalog/api"

	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsentry"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Options - comment struct.
	Options struct {
		Cfg          config.Config
		ErrorHandler mrcore.ErrorHandler
		EventEmitter mrsender.EventEmitter

		Sentry              *mrsentry.Adapter
		Prometheus          *prometheus.Registry
		ErrorManager        *mrinit.ErrorManager
		UsecaseErrorWrapper mrcore.UsecaseErrorWrapper

		PostgresConnManager *mrpostgres.ConnManager
		RedisAdapter        *mrredis.ConnAdapter
		FileProviderPool    *mrstorage.FileProviderPool
		Locker              mrlock.Locker
		Translator          *mrlang.Translator
		RequestParsers      RequestParsers
		ResponseSenders     ResponseSenders
		AccessControl       *mrperms.RoleAccessControl
		ImageURLBuilder     mrpath.PathBuilder

		// API section
		CatalogCategoryAPI  api.CategoryAPI
		CatalogTrademarkAPI api.TrademarkAPI
		OrdererAPI          mrsort.Orderer
		SettingsGetterAPI   mrsettings.DefaultValueGetter
		SettingsSetterAPI   mrsettings.Setter

		// Modules section
		CatalogCategoryModule  factory_catalog_category.Options
		CatalogProductModule   factory_catalog_product.Options
		CatalogTrademarkModule factory_catalog_trademark.Options
		FileStationModule      factory_filestation.Options

		SchedulerTasks  []mrworker.Task
		OpenedResources []func(ctx context.Context)
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		KeyInt32   *mrparser.KeyInt32
		ListSorter *mrparser.ListSorter
		ListPager  *mrparser.ListPager
		String     *mrparser.String
		UUID       *mrparser.UUID
		Validator  *mrparser.Validator
		// File       *mrparser.File
		Image *mrparser.Image
	}

	// ResponseSenders - comment struct.
	ResponseSenders struct {
		Sender     *mrresp.Sender
		FileSender *mrresp.FileSender
	}
)
