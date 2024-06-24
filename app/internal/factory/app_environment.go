package factory

import (
	"context"

	"go-sample/config"
	"go-sample/internal/app"
	factory_catalog "go-sample/internal/factory/modules/catalog"
	factory_filestation "go-sample/internal/factory/modules/filestation"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-webcore/mrcore/mrcoreerr"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender/mrlogadapter"
)

// CreateAppEnvironment - comment func.
func CreateAppEnvironment(configPath, logLevel string) (context.Context, app.Options, error) {
	cfg, err := config.Create(configPath)
	if err != nil {
		return nil, app.Options{}, err
	}

	// create and init logger
	if logLevel != "" {
		cfg.Log.Level = logLevel
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		return nil, app.Options{}, err
	}

	if err = mrlog.SetDefault(logger); err != nil {
		return nil, app.Options{}, err
	}

	// show head info about started app
	logger.Info().Msgf("%s, version: %s, environment: %s", cfg.AppName, cfg.AppVersion, cfg.AppEnvironment)

	if cfg.Debugging.Debug {
		logger.Info().Msg("DEBUG MODE: ON")
	}

	logger.Info().Msgf("LOG LEVEL: %s", logger.Level())
	logger.Debug().Msgf("CONFIG PATH: %s", cfg.ConfigPath)

	ctx := mrlog.WithContext(context.Background(), logger)

	opts := app.Options{
		Cfg:          cfg,
		ErrorHandler: NewErrorHandler(logger, cfg),
		EventEmitter: mrlogadapter.NewEventEmitter(logger),
	}

	return ctx, opts, nil
}

// InitAppEnvironment - comment func.
func InitAppEnvironment(ctx context.Context, opts app.Options) (app.Options, error) {
	// init shared options
	if opts.Cfg.Sentry.Enable {
		sentry, err := NewSentry(ctx, opts.Cfg)
		if err != nil {
			return opts, err
		}

		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(sentry))

		opts.Sentry = sentry
	}

	opts.Prometheus = NewPrometheusRegistry(ctx, opts)

	opts.ErrorManager = NewErrorManager(opts)
	opts.UsecaseErrorWrapper = mrcoreerr.NewUsecaseErrorWrapper()

	postgresAdapter, err := NewPostgres(ctx, opts)
	if err != nil {
		return opts, err
	}

	opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(postgresAdapter))

	opts.PostgresConnManager = NewPostgresConnManager(ctx, postgresAdapter)

	opts.RedisAdapter, err = NewRedis(ctx, opts.Cfg)
	if err != nil {
		return opts, err
	}

	opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(opts.RedisAdapter))

	if opts.FileProviderPool, err = NewFileProviderPool(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	opts.Locker = mrredislock.NewLockerAdapter(opts.RedisAdapter.Cli())

	if opts.Translator, err = NewTranslator(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.RequestParsers, err = CreateRequestParsers(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.ResponseSenders, err = CreateResponseSenders(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.AccessControl, err = NewAccessControl(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.ImageURLBuilder, err = NewImageURLBuilder(opts.Cfg); err != nil {
		return opts, err
	}

	// Register errors (!!! only after init opts)
	factory_catalog.RegisterCategoryErrors(opts.ErrorManager)
	factory_catalog.RegisterProductErrors(opts.ErrorManager)
	factory_catalog.RegisterTrademarkErrors(opts.ErrorManager)

	// Shared APIs init section (!!! only after init opts)
	if opts.CatalogCategoryAPI, err = factory_catalog.NewCategoryAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogTrademarkAPI, err = factory_catalog.NewTrademarkAPI(ctx, opts); err != nil {
		return opts, err
	}

	opts.OrdererAPI = NewOrdererAPI(ctx, opts)

	{
		getter, task := NewSettingsGetter(ctx, opts)
		opts.SettingsGetterAPI = getter
		opts.SchedulerTasks = append(opts.SchedulerTasks, task)
	}

	opts.SettingsSetterAPI = NewSettingsSetter(ctx, opts)

	// Shared module's options (!!! only after init APIs)
	if opts.CatalogCategoryModule, err = factory_catalog.NewCategoryModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogProductModule, err = factory_catalog.NewProductModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogTrademarkModule, err = factory_catalog.NewTrademarkModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.FileStationModule, err = factory_filestation.NewModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	return opts, nil
}
