package factory

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-webcore/mrcore/mrcoreerr"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrsender/mrlogadapter"

	"github.com/mondegor/go-sample/cmd/factory/catalog"
	"github.com/mondegor/go-sample/cmd/factory/filestation"
	"github.com/mondegor/go-sample/config"
	"github.com/mondegor/go-sample/internal/app"
)

// InitApp - Настраивает конфигурацию, внешнее окружение приложения, после этого создаёт её модули и компоненты.
func InitApp(ctx context.Context, args []string, stdout io.Writer) (context.Context, app.Options, error) {
	parsedArgs, err := ParseAppArgs(args)
	if err != nil {
		return nil, app.Options{}, err
	}

	cfg, err := config.Create(
		config.Args{
			WorkDir:     parsedArgs.WorkDir,
			ConfigPath:  parsedArgs.ConfigPath,
			DotEnvPath:  parsedArgs.DotEnvPath,
			Environment: parsedArgs.Environment,
			LogLevel:    parsedArgs.LogLevel,
			Stdout:      stdout,
		},
	)
	if err != nil {
		return nil, app.Options{}, fmt.Errorf("factory.InitApp(): %w", err)
	}

	return InitAppEnvironment(ctx, app.Options{Cfg: cfg})
}

// InitAppEnvironment - Настраивает внешнее окружение приложения на основе переданной конфигурации,
// после этого создаёт её модули и компоненты.
// Имеется возможность заранее задать некоторые параметры и компонентов приложения (актуально для использования в тестах):
// К ним относится: opts.PostgresConnManager, opts.RedisAdapter, opts.FileProviderPool.
func InitAppEnvironment(ctx context.Context, opts app.Options) (context.Context, app.Options, error) {
	ctx, err := InitLogger(ctx, opts.Cfg)
	if err != nil {
		return nil, app.Options{}, err
	}

	logger := mrlog.Ctx(ctx)

	// show head info about started app
	logger.Info().Msgf("%s, environment: %s, version: %s", opts.Cfg.App.Name, opts.Cfg.App.Environment, opts.Cfg.App.Version)

	if opts.Cfg.Debugging.Debug {
		logger.Info().Msg("DEBUG MODE: ON")
	}

	logger.Info().Msgf("LOG LEVEL: %s", logger.Level())

	if opts.Cfg.App.WorkDir != "" {
		logger.Debug().Msgf("WORK DIR: %s", opts.Cfg.App.WorkDir)
	}

	logger.Debug().Msgf("CONFIG PATH: %s", opts.Cfg.ConfigPath)

	if opts.Cfg.App.DotEnvPath != "" {
		logger.Debug().Msgf(".ENV PATH: %s", opts.Cfg.App.DotEnvPath)
	}

	opts.AppHealth = mrrun.NewAppHealth()
	opts.ErrorHandler = NewErrorHandler(logger, opts.Cfg)
	opts.EventEmitter = mrlogadapter.NewEventEmitter(logger)

	opts, err = createAppEnvironment(ctx, opts)
	if err != nil {
		return nil, app.Options{}, err
	}

	// Register errors (!!! only after create environment)
	registerAppErrors(opts)

	// Shared APIs init section (!!! only after create environment)
	if opts, err = createAppAPI(ctx, opts); err != nil {
		return nil, app.Options{}, err
	}

	// Shared module's options (!!! only after create APIs)
	if opts, err = createAppModules(ctx, opts); err != nil {
		return nil, app.Options{}, err
	}

	return ctx, opts, nil
}

// createAppEnvironment - создаёт, и настраивает внешнее окружение приложения.
func createAppEnvironment(ctx context.Context, opts app.Options) (enrichedOpts app.Options, err error) {
	opts.ErrorManager = NewErrorManager(opts)
	opts.UseCaseErrorWrapper = mrcoreerr.NewUseCaseErrorWrapper()
	opts.InternalRouter = http.NewServeMux()

	if opts.Cfg.Sentry.DSN != "" {
		sentry, err := NewSentry(ctx, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(sentry))
		opts.Sentry = sentry
	}

	opts.Prometheus = NewPrometheusRegistry(ctx, opts)

	if opts.PostgresConnManager == nil {
		postgresAdapter, err := NewPostgres(ctx, opts)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(postgresAdapter))
		opts.PostgresConnManager = NewPostgresConnManager(ctx, postgresAdapter)

		if opts.Cfg.Storage.MigrationsDir != "" {
			if err = ApplyPostgresMigrations(ctx, opts); err != nil {
				return app.Options{}, err
			}
		}
	}

	if opts.RedisAdapter == nil {
		redisAdapter, err := NewRedis(ctx, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(redisAdapter))
		opts.RedisAdapter = redisAdapter
	}

	if opts.FileProviderPool == nil {
		opts.FileProviderPool, err = NewFileProviderPool(ctx, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(opts.FileProviderPool))
	}

	opts.Locker = mrredislock.NewLockerAdapter(opts.RedisAdapter.Cli())

	if opts.Translator, err = NewTranslator(ctx, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.RequestParsers, err = CreateRequestParsers(ctx, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.ResponseSenders, err = CreateResponseSenders(ctx, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.AccessControl, err = NewAccessControl(ctx, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.ImageURLBuilder, err = NewImageURLBuilder(opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if err = RegisterSystemHandlers(ctx, opts); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

func registerAppErrors(opts app.Options) {
	catalog.RegisterCategoryErrors(opts.ErrorManager)
	catalog.RegisterProductErrors(opts.ErrorManager)
	catalog.RegisterTrademarkErrors(opts.ErrorManager)
}

func createAppAPI(ctx context.Context, opts app.Options) (enrichedOpts app.Options, err error) {
	opts.OrdererAPI = NewOrdererAPI(ctx, opts)

	{
		getter, task := NewSettingsGetterAndTask(ctx, opts)
		opts.SettingsGetterAPI = getter
		opts.SchedulerTasks = append(opts.SchedulerTasks, task)
	}

	opts.SettingsSetterAPI = NewSettingsSetter(ctx, opts)

	if opts.CatalogCategoryAvailabilityAPI, err = catalog.NewCategoryAvailabilityAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogTrademarkAvailabilityAPI, err = catalog.NewTrademarkAvailabilityAPI(ctx, opts); err != nil {
		return opts, err
	}

	return opts, nil
}

func createAppModules(ctx context.Context, opts app.Options) (enrichedOpts app.Options, err error) {
	if opts.CatalogCategoryModule, err = catalog.NewCategoryModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogProductModule, err = catalog.NewProductModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogTrademarkModule, err = catalog.NewTrademarkModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.FileStationModule, err = filestation.NewModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	return opts, nil
}
