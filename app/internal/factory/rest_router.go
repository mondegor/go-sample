package factory

import (
	"context"
	"go-sample/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"
)

func NewRestRouter(ctx context.Context, cfg config.Config, translator *mrlang.Translator) (*mrchi.RouterAdapter, error) {
	logger := mrlog.Ctx(ctx)

	corsOptions := mrrscors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		Logger:           logger.With().Str("middleware", "cors").Logger(),
	}

	errorSender, err := NewErrorResponseSender(ctx, cfg)

	if err != nil {
		return nil, err
	}

	router := mrchi.New(
		logger.With().Str("router", "chi").Logger(),
		mrserver.MiddlewareHandlerAdapter(errorSender),
		mrresp.HandlerGetNotFoundAsJson(),
		mrresp.HandlerGetMethodNotAllowedAsJson(),
	)
	router.RegisterMiddleware(
		mrrscors.Middleware(corsOptions),
		mrserver.MiddlewareGeneral(translator, mrresp.ApplyStatRequest),
		mrserver.MiddlewareRecoverHandler(
			mrdebug.IsDebug(),
			mrresp.HandlerGetFatalErrorAsJson(),
		),
	)

	return router, nil
}
