package factory

import (
	"context"
	"go-sample/config"
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

func RegisterSystemHandlers(
	ctx context.Context,
	cfg config.Config,
	router mrserver.HttpRouter,
	section mrperms.AppSection,
) error {
	mrlog.Ctx(ctx).Info().Msgf("Init system handlers in section %s", section.Caption())

	router.HandlerFunc(http.MethodGet, section.Path("/"), mrresp.HandlerGetStatusOKAsJson())
	router.HandlerFunc(http.MethodGet, section.Path("/v1/health"), mrresp.HandlerGetHealth())
	router.HandlerFunc(http.MethodGet, section.Path("/v1/stat-info"), mrresp.HandlerGetStatInfoAsJson())

	systemInfoFunc, err := mrresp.HandlerGetSystemInfoAsJson(
		mrresp.SystemInfoConfig{
			Name:      cfg.AppName,
			Version:   cfg.AppVersion,
			StartedAt: cfg.AppStartedAt,
		},
	)

	if err != nil {
		return err
	}

	router.HandlerFunc(http.MethodGet, section.Path("/v1/system-info"), systemInfoFunc)

	return nil
}
