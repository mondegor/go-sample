package app

import (
	"context"
	"net/http"

	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsentry"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mondegor/go-sample/config"
	catalogcategory "github.com/mondegor/go-sample/internal/factory/catalog/category"
	catalogproduct "github.com/mondegor/go-sample/internal/factory/catalog/product"
	catalogtrademark "github.com/mondegor/go-sample/internal/factory/catalog/trademark"
	"github.com/mondegor/go-sample/internal/factory/filestation"
	"github.com/mondegor/go-sample/pkg/catalog/api"
	"github.com/mondegor/go-sample/pkg/validate"
)

type (
	// Options - comment struct.
	Options struct {
		Cfg          config.Config
		AppHealth    *mrrun.AppHealth
		ErrorHandler mrcore.ErrorHandler
		EventEmitter mrsender.EventEmitter

		InternalRouter *http.ServeMux
		Sentry         *mrsentry.Adapter
		Prometheus     *prometheus.Registry

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
		CatalogCategoryAvailabilityAPI  api.CategoryAvailability
		CatalogTrademarkAvailabilityAPI api.TrademarkAvailability
		SettingsGetterAPI               mrsettings.DefaultValueGetter
		SettingsSetterAPI               mrsettings.Setter

		// Modules section
		CatalogCategoryModule  catalogcategory.Options
		CatalogProductModule   catalogproduct.Options
		CatalogTrademarkModule catalogtrademark.Options
		FileStationModule      filestation.Options

		SchedulerTasks  []mrworker.Task
		OpenedResources []func(ctx context.Context)
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Bool       *mrparser.Bool
		// DateTime   *mrparser.DateTime
		Int64      *mrparser.Int64
		ItemStatus *mrparser.ItemStatus
		Uint64     *mrparser.Uint64
		ListSorter *mrparser.ListSorter
		ListPager  *mrparser.ListPager
		String     *mrparser.String
		UUID       *mrparser.UUID
		Validator  *mrparser.Validator
		// File       *mrparser.File
		Image *mrparser.Image

		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}

	// ResponseSenders - comment struct.
	ResponseSenders struct {
		Sender     *mrresp.Sender
		FileSender *mrresp.FileSender
	}
)
