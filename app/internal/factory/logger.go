package factory

import (
	"go-sample/config"

	"github.com/mondegor/go-webcore/mrcore"
)

func NewLogger(cfg *config.Config) (*mrcore.LoggerAdapter, error) {
	prefix := cfg.Log.Prefix

	if prefix != "" {
		prefix = "[" + prefix + "] "
	}

	logger, err := mrcore.NewLogger(prefix, cfg.Log.Level)

	if err != nil {
		return nil, err
	}

	mrcore.SetDefaultLogger(logger)

	logger.Info("%s, version: %s", cfg.AppName, cfg.AppVersion)

	if cfg.AppInfo != "" {
		logger.Info(cfg.AppInfo)
	}

	if mrcore.Debug() {
		logger.Info("DEBUG MODE: ON")
	}

	logger.Info("LOG LEVEL: %s", cfg.Log.Level)
	logger.Info("CONFIG PATH: %s", cfg.ConfigPath)
	logger.Info("APP PATH: %s", cfg.AppPath)

	return logger, nil
}
