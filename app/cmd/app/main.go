package main

import (
	"context"
	"flag"
	"go-sample/config"
	"go-sample/internal/factory"
	"go-sample/internal/modules"
	factory_catalog_adm "go-sample/internal/modules/catalog/factory/admin-api"
	factory_catalog_pub "go-sample/internal/modules/catalog/factory/public-api"
	factory_filestation_pub "go-sample/internal/modules/file-station/factory/public-api"
	"log"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrtool"
)

var (
	configPath string
	logLevel   string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "Path to application config file")
	flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
	flag.Parse()

	cfg, err := config.New(configPath)

	if err != nil {
		log.Fatal(err)
	}

	if logLevel != "" {
		cfg.Log.Level = logLevel
	}

	logger, err := factory.NewLogger(cfg)

	if err != nil {
		log.Fatal(err)
	}

	sharedOptions := &modules.Options{}
	sharedOptions.Cfg = cfg
	sharedOptions.Logger = logger
	sharedOptions.EventBox = logger

	appHelper := mrtool.NewAppHelper(logger)

	sharedOptions.ServiceHelper = mrtool.NewServiceHelper()

	sharedOptions.PostgresAdapter, err = factory.NewPostgres(cfg, logger)
	appHelper.ExitOnError(err)
	defer appHelper.Close(sharedOptions.PostgresAdapter)

	sharedOptions.RedisAdapter, err = factory.NewRedis(cfg, logger)
	appHelper.ExitOnError(err)
	defer appHelper.Close(sharedOptions.RedisAdapter)

	sharedOptions.FileProviderPool, err = factory.NewFileProviderPool(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.Locker = mrredislock.NewLockerAdapter(sharedOptions.RedisAdapter.Cli())

	sharedOptions.Translator, err = factory.NewTranslator(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.RequestParsers, err = factory.NewRequestParsers(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.ResponseSender, err = factory.NewResponseSender(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.AccessControl, err = factory.NewAccessControl(cfg, logger)
	appHelper.ExitOnError(err)

	// API
	sharedOptions.CatalogCategoryAPI, err = factory.NewCatalogCategoryAPI(sharedOptions)
	appHelper.ExitOnError(err)

	sharedOptions.CatalogTrademarkAPI, err = factory.NewCatalogTrademarkAPI(sharedOptions)
	appHelper.ExitOnError(err)

	sharedOptions.OrdererAPI = factory.NewOrdererAPI(sharedOptions)

	// module's options
	catalogOptions, err := factory.NewCatalogOptions(sharedOptions)
	appHelper.ExitOnError(err)

	fileStationOptions, err := factory.NewFileStationOptions(sharedOptions)
	appHelper.ExitOnError(err)

	// http router
	router, err := factory.NewHttpRouter(cfg, logger, sharedOptions.Translator)
	appHelper.ExitOnError(err)

	// section: admin-api
	sectionAdminAPI := factory.NewAppSectionAdminAPI(sharedOptions)

	appHelper.ExitOnError(
		factory.RegisterSystemHandlers(cfg, logger, router, sectionAdminAPI),
	)

	controllers, err := factory_catalog_adm.CreateModule(catalogOptions)
	appHelper.ExitOnError(err)
	router.Register(
		mrfactory.WithMiddlewareCheckAccess(controllers, sectionAdminAPI, sharedOptions.AccessControl)...,
	)

	// section: public
	sectionPublicAPI := factory.NewAppSectionPublicAPI(sharedOptions)

	appHelper.ExitOnError(
		factory.RegisterSystemHandlers(cfg, logger, router, sectionPublicAPI),
	)

	controllers, err = factory_catalog_pub.CreateModule(catalogOptions)
	appHelper.ExitOnError(err)
	router.Register(
		mrfactory.WithMiddlewareCheckAccess(controllers, sectionPublicAPI, sharedOptions.AccessControl)...,
	)

	controllers, err = factory_filestation_pub.CreateModule(fileStationOptions)
	appHelper.ExitOnError(err)
	router.Register(
		mrfactory.WithMiddlewareCheckAccess(controllers, sectionPublicAPI, sharedOptions.AccessControl)...,
	)

	// http server
	serverAdapter, err := factory.NewHttpServer(cfg, logger, router)
	appHelper.ExitOnError(err)
	defer appHelper.Close(serverAdapter)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go appHelper.GracefulShutdown(cancel)

	logger.Info("Waiting for requests. To exit press CTRL+C")

	select {
	case <-ctx.Done():
		logger.Info("Application stopped")
	case err = <-serverAdapter.Notify():
		logger.Info("Application stopped with error")
	}

	logger.Err(err)
}
