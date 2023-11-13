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
	"net/http"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-webcore/mrtool"
)

var (
	configPath string
	logLevel string
)

func init() {
	flag.StringVar(&configPath,"config-path", "./config/config.yaml", "Path to application config file")
	flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
	flag.Parse()

	sharedOptions := &modules.Options{}
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

	// fileProviderAPI, err := factory.NewFileStorage(cfg, logger)
	sharedOptions.MinioAdapter, err = factory.NewS3Minio(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.Locker = mrredislock.NewLockerAdapter(sharedOptions.RedisAdapter.Cli())
	sharedOptions.OrdererComponent = factory.NewOrdererComponent(cfg, sharedOptions.PostgresAdapter, logger, logger)

	// access + sections + router
	modulesAccess, err := factory.NewModulesAccess(cfg, logger)
	appHelper.ExitOnError(err)

	sectionAdminAPI := factory.NewClientSectionAdminAPI(cfg, modulesAccess)
	sectionPublicAPI := factory.NewClientSectionPublicAPI(cfg, modulesAccess)

	router, err := factory.NewHttpRouter(cfg, logger)
	appHelper.ExitOnError(err)

	// section: admin-api
	controllers, err := factory_catalog_adm.NewModule(sharedOptions, sectionAdminAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	// section: public
	controllers, err = factory_catalog_pub.NewModule(sharedOptions, sectionPublicAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	controllers, err = factory_filestation_pub.NewModule(sharedOptions, sectionPublicAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

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
		err = serverAdapter.Close()
		logger.Info("Application stopped")
	case err = <-serverAdapter.Notify():
		logger.Info("Application stopped with error")
	}

	if err != nil && err != http.ErrServerClosed {
		logger.Err(err)
	}
}
