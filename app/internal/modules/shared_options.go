package modules

import (
	"go-sample/config"
	"go-sample/pkg/modules/catalog"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

type Options struct {
	Cfg              *config.Config
	Logger           mrcore.Logger
	EventBox         mrcore.EventBox
	ServiceHelper    *mrtool.ServiceHelper
	PostgresAdapter  *mrpostgres.ConnAdapter
	RedisAdapter     *mrredis.ConnAdapter
	FileProviderPool *mrstorage.FileProviderPool
	Locker           mrcore.Locker
	OrdererAPI       mrorderer.API
	Translator       *mrlang.Translator

	CatalogCategoryAPI  catalog.CategoryAPI
	CatalogTrademarkAPI catalog.TrademarkAPI
}
