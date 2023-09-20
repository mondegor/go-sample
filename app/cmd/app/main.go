package main

import (
    "context"
    "flag"
    "go-sample/config"
    "go-sample/internal/controller/http_v1"
    "go-sample/internal/factory"
    "go-sample/internal/infrastructure/repository"
    "go-sample/internal/usecase"
    "log"
    "net/http"

    sq "github.com/Masterminds/squirrel"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
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

    appHelper := mrtool.NewAppHelper(logger)
    serviceHelper := mrtool.NewServiceHelper()

    postgresAdapter, err := factory.NewPostgres(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(postgresAdapter)

    redisAdapter, err := factory.NewRedis(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(redisAdapter)

    // fileStorage, err := factory.NewFileStorage(cfg, logger)
    fileStorage, err := factory.NewS3Minio(cfg, logger)
    appHelper.ExitOnError(err)

    lockerAdapter := mrredislock.NewLockerAdapter(redisAdapter.Cli())
    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    itemOrdererStorage := mrcom_orderer.NewRepository(postgresAdapter, queryBuilder)
    itemOrdererComponent := mrcom_orderer.NewComponent(itemOrdererStorage, logger)

    catalogCategoryStorage := repository.NewCatalogCategory(postgresAdapter, queryBuilder)
    catalogCategoryImageStorage := repository.NewCatalogCategoryImage(postgresAdapter, queryBuilder)
    catalogCategoryService := usecase.NewCatalogCategory(catalogCategoryStorage, logger, serviceHelper)
    catalogCategoryImageService := usecase.NewCatalogCategoryImage(cfg.FileStorage.CatalogCategoryImageDir, catalogCategoryImageStorage, fileStorage, lockerAdapter, logger, serviceHelper)
    catalogCategoryHttp := http_v1.NewCatalogCategory(catalogCategoryService, catalogCategoryImageService)

    catalogTrademarkStorage := repository.NewCatalogTrademark(postgresAdapter, queryBuilder)
    catalogTrademarkService := usecase.NewCatalogTrademark(catalogTrademarkStorage, logger, serviceHelper)
    catalogTrademarkHttp := http_v1.NewCatalogTrademark(catalogTrademarkService)

    catalogProductStorage := repository.NewCatalogProduct(postgresAdapter, queryBuilder)
    catalogProductService := usecase.NewCatalogProduct(itemOrdererComponent, catalogProductStorage, catalogTrademarkStorage, logger, serviceHelper)
    catalogProductHttp := http_v1.NewCatalogProduct(catalogProductService, catalogCategoryService, catalogTrademarkService)

    router, err := factory.NewHttpRouter(cfg, logger)
    appHelper.ExitOnError(err)

    router.Register(
        catalogCategoryHttp,
        catalogTrademarkHttp,
        catalogProductHttp,
    )

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
