package modules

import (
	"go-sample/config"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-storage/mrstorage"
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
	S3Pool           *mrstorage.FileProviderPool
	Locker           mrcore.Locker
	OrdererComponent mrorderer.Component
}
