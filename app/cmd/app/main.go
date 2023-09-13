package main

import (
    "context"
    "flag"
    "go-sample/config"
    "go-sample/internal/controller/http_v1"
    "go-sample/internal/controller/view"
    "go-sample/internal/factory"
    "go-sample/internal/infrastructure/repository"
    "go-sample/internal/usecase"
    "log"
    "net/http"
    "time"

    sq "github.com/Masterminds/squirrel"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-sysmess/mrlang"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrserver"
    "github.com/mondegor/go-webcore/mrtool"
    "github.com/mondegor/go-webcore/mrview"
)

const appName = "go-sample"
const appVersion = "v0.1.0"

var configPath string
var logLevel string

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

    if logLevel == "" {
        logLevel = cfg.Log.Level
    }

    logger, err := mrcore.NewLogger(appName, logLevel)

    if err != nil {
        log.Fatal(err)
    }

    logger.Info("APP VERSION: %s", appVersion)

    if cfg.Debug {
        logger.Info("DEBUG MODE: ON")
    }

    logger.Info("LOG LEVEL: %s", cfg.Log.Level)
    logger.Info("APP PATH: %s", cfg.AppPath)
    logger.Info("CONFIG PATH: %s", configPath)

    appHelper := mrtool.NewAppHelper(logger)

    responseTranslator, err := mrlang.NewTranslator(
        mrlang.TranslatorOptions{
            DirPath: cfg.Translation.DirPath,
            FileType: cfg.Translation.FileType,
            LangCodes: cfg.Translation.LangCodes,
        },
    )
    appHelper.ExitOnError(err)

    postgresClient, err := factory.NewPostgres(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(postgresClient)

    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    requestValidator := mrview.NewValidator()
    appHelper.ExitOnError(requestValidator.Register("article", view.ValidateArticle))

    serviceHelper := mrtool.NewServiceHelper()

    itemOrdererStorage := mrcom_orderer.NewRepository(postgresClient, queryBuilder)
    itemOrdererComponent := mrcom_orderer.NewComponent(itemOrdererStorage, logger)

    catalogCategoryStorage := repository.NewCatalogCategory(postgresClient, queryBuilder)
    catalogCategoryService := usecase.NewCatalogCategory(catalogCategoryStorage, logger, serviceHelper)
    catalogCategoryHttp := http_v1.NewCatalogCategory(catalogCategoryService)

    catalogTrademarkStorage := repository.NewCatalogTrademark(postgresClient, queryBuilder)
    catalogTrademarkService := usecase.NewCatalogTrademark(catalogTrademarkStorage, logger, serviceHelper)
    catalogTrademarkHttp := http_v1.NewCatalogTrademark(catalogTrademarkService)

    catalogProductStorage := repository.NewCatalogProduct(postgresClient, queryBuilder)
    catalogProductService := usecase.NewCatalogProduct(itemOrdererComponent, catalogProductStorage, catalogTrademarkStorage, logger, serviceHelper)
    catalogProductHttp := http_v1.NewCatalogProduct(catalogProductService, catalogCategoryService, catalogTrademarkService)


    logger.Info("Create router")

    corsOptions := mrserver.CorsOptions{
        AllowedOrigins: cfg.Cors.AllowedOrigins,
        AllowedMethods: cfg.Cors.AllowedMethods,
        AllowedHeaders: cfg.Cors.AllowedHeaders,
        ExposedHeaders: cfg.Cors.ExposedHeaders,
        AllowCredentials: cfg.Cors.AllowCredentials,
        Debug: cfg.Debug,
    }

    router := mrserver.NewRouter(logger, mrserver.HandlerAdapter(requestValidator))
    router.RegisterMiddleware(
        mrserver.NewCors(corsOptions),
        mrserver.MiddlewareFirst(logger),
        mrserver.MiddlewareUserIp(),
        mrserver.MiddlewareAcceptLanguage(responseTranslator),
        mrserver.MiddlewarePlatform(mrcore.PlatformWeb),
        mrserver.MiddlewareAuthenticateUser(),
    )

    router.Register(
        catalogCategoryHttp,
        catalogTrademarkHttp,
        catalogProductHttp,
    )
    router.HandlerFunc(http.MethodGet, "/", MainPage)

    logger.Info("Initialize application")

    server := mrserver.NewServer(logger, mrserver.ServerOptions{
        Handler: router,
        ReadTimeout: time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
        ShutdownTimeout: time.Duration(cfg.Server.ShutdownTimeout) * time.Second,
    })

    logger.Info("Start application")

    err = server.Start(mrserver.ListenOptions{
        AppPath: cfg.AppPath,
        Type: cfg.Listen.Type,
        SockName: cfg.Listen.SockName,
        BindIP: cfg.Listen.BindIP,
        Port: cfg.Listen.Port,
    })
    appHelper.ExitOnError(err)
    defer appHelper.Close(server)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go appHelper.GracefulShutdown(cancel)

    logger.Info("Waiting for requests. To exit press CTRL+C")

    select {
    case <-ctx.Done():
        err = server.Close()
        logger.Info("Application stopped")
    case err = <-server.Notify():
        logger.Info("Application stopped with error")
    }

    if err != nil && err != http.ErrServerClosed {
        logger.Err(err)
    }
}

func MainPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"STATUS\": \"OK\"}"))
}
