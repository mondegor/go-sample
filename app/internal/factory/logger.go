package factory

import (
	"go-sample/config"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog/zerolog"
	"github.com/mondegor/go-webcore/mrlog/zerolog/factory"
)

// NewLogger - comment func.
func NewLogger(cfg config.Config) (*zerolog.LoggerAdapter, error) {
	return factory.NewZeroLogAdapter(
		factory.Options{
			Level:            cfg.Log.Level,
			JsonFormat:       cfg.Log.JsonFormat,
			TimestampFormat:  cfg.Log.TimestampFormat,
			ConsoleColor:     cfg.Log.ConsoleColor,
			PrepareErrorFunc: mrcore.PrepareError,
		},
	)
}
