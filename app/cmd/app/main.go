package main

import (
    "context"
    "flag"
    "go-sample/config"
    http_v1_adm "go-sample/internal/controller/http_v1/admin-panel"
    http_v1_public "go-sample/internal/controller/http_v1/public"
    entity_adm "go-sample/internal/entity/admin-panel"
    "go-sample/internal/factory"
    repository_adm "go-sample/internal/infrastructure/repository/admin-panel"
    repository_public "go-sample/internal/infrastructure/repository/public"
    usecase_adm "go-sample/internal/usecase/admin-panel"
    usecase_public "go-sample/internal/usecase/public"
    "log"
    "net/http"

    "github.com/mondegor/go-components/mrorderer"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-storage/mrredislock"
    "github.com/mondegor/go-storage/mrsql"
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

    mrpostgresWhere := mrpostgres.NewSqlBuilderWhere()
    mrpostgresPager := mrpostgres.NewSqlBuilderPager(1000)

    itemOrdererStorage := mrorderer.NewRepository(postgresAdapter)
    itemOrdererComponent := mrorderer.NewComponent(itemOrdererStorage, logger)

    modulesAccess, err := factory.NewModulesAccess(cfg, logger)
    appHelper.ExitOnError(err)

    sectionPublic := factory.NewClientSectionPublic(cfg, modulesAccess)
    sectionAdminPanel := factory.NewClientSectionAdminPanel(cfg, modulesAccess)

    // section: public
    publicCatalogCategoryStorage := repository_public.NewCatalogCategory(
        postgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgresWhere,
            mrpostgres.NewSqlBuilderOrderBy("category_caption", mrentity.SortDirectionASC),
            mrpostgresPager,
        ),
    )
    publicCatalogCategoryService := usecase_public.NewCatalogCategory(publicCatalogCategoryStorage, serviceHelper)
    publicCatalogCategoryHttp := http_v1_public.NewCatalogCategory(sectionAdminPanel, publicCatalogCategoryService)

    fileService := usecase_public.NewFileItem(fileStorage)
    publicImageHttp := http_v1_public.NewImageItem(sectionPublic, fileService)

    // section: admin-panel
    catalogCategoryEntityMetaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity_adm.CatalogCategory{})
    appHelper.ExitOnError(err)

    catalogCategoryStorage := repository_adm.NewCatalogCategory(
        postgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgresWhere,
            mrpostgres.NewSqlBuilderOrderByWithFieldMap(
                catalogCategoryEntityMetaOrderBy.FieldMap,
                catalogCategoryEntityMetaOrderBy.DefaultDbField,
                mrentity.SortDirectionASC,
            ),
            mrpostgresPager,
        ),
    )
    catalogCategoryImageStorage := repository_adm.NewCatalogCategoryImage(postgresAdapter)
    catalogCategoryService := usecase_adm.NewCatalogCategory(catalogCategoryStorage, logger, serviceHelper)
    catalogCategoryImageService := usecase_adm.NewCatalogCategoryImage(cfg.FileStorage.CatalogCategoryImageDir, catalogCategoryImageStorage, fileStorage, lockerAdapter, logger, serviceHelper)
    catalogCategoryHttp := http_v1_adm.NewCatalogCategory(sectionAdminPanel, catalogCategoryService, catalogCategoryImageService)

    catalogTrademarkEntityMetaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity_adm.CatalogTrademark{})
    appHelper.ExitOnError(err)

    catalogTrademarkStorage := repository_adm.NewCatalogTrademark(
        postgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgresWhere,
            mrpostgres.NewSqlBuilderOrderByWithFieldMap(
                catalogTrademarkEntityMetaOrderBy.FieldMap,
                catalogTrademarkEntityMetaOrderBy.DefaultDbField,
                mrentity.SortDirectionASC,
            ),
            mrpostgresPager,
        ),
    )
    catalogTrademarkService := usecase_adm.NewCatalogTrademark(catalogTrademarkStorage, logger, serviceHelper)
    catalogTrademarkHttp := http_v1_adm.NewCatalogTrademark(sectionAdminPanel, catalogTrademarkService)

    catalogProductEntityMetaOrderBy, err := mrsql.NewEntityMetaOrderBy(entity_adm.CatalogProduct{})
    appHelper.ExitOnError(err)

    catalogProductEntityMetaUpdate, err := mrsql.NewEntityMetaUpdate(entity_adm.CatalogProduct{})
    appHelper.ExitOnError(err)

    catalogProductStorage := repository_adm.NewCatalogProduct(
        postgresAdapter,
        mrsql.NewBuilderSelect(
            mrpostgresWhere,
            mrpostgres.NewSqlBuilderOrderByWithFieldMap(
                catalogProductEntityMetaOrderBy.FieldMap,
                catalogProductEntityMetaOrderBy.DefaultDbField,
                mrentity.SortDirectionASC,
            ),
            mrpostgresPager,
        ),
        mrsql.NewBuilderUpdateWithMeta(
            catalogProductEntityMetaUpdate,
            mrpostgres.NewSqlBuilderSet(),
        ),
    )
    catalogProductService := usecase_adm.NewCatalogProduct(itemOrdererComponent, catalogProductStorage, catalogTrademarkStorage, logger, serviceHelper)
    catalogProductHttp := http_v1_adm.NewCatalogProduct(sectionAdminPanel, catalogProductService, catalogCategoryService, catalogTrademarkService)

    router, err := factory.NewHttpRouter(cfg, logger)
    appHelper.ExitOnError(err)

    router.Register(
        // section: public
        publicCatalogCategoryHttp,
        publicImageHttp,

        // section: admin-panel
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
