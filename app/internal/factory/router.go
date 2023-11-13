package factory

import (
	"go-sample/config"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

func NewHttpRouter(cfg *config.Config, logger mrcore.Logger) (mrcore.HttpRouter, error) {
	responseTranslator, err := NewTranslator(cfg, logger)

	if err != nil {
		return nil, err
	}

	requestValidator, err := NewValidator(cfg, logger)

	if err != nil {
		return nil, err
	}

	logger.Info("Create and init http router")

	corsOptions := mrserver.CorsOptions{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		Debug:            cfg.Debug,
	}

	router := mrserver.NewRouter(logger, mrserver.HandlerAdapter(requestValidator))
	router.RegisterMiddleware(
		mrserver.NewCors(corsOptions),
		mrserver.MiddlewareFirst(logger),
		mrserver.MiddlewareAcceptLanguage(responseTranslator),
	)

	router.HandlerFunc(http.MethodGet, "/", mrserver.MainPage)

	return router, nil
}
