package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"

	"github.com/mondegor/go-sample/config"
	"github.com/mondegor/go-sample/internal/app"
	"github.com/mondegor/go-sample/pkg/validate"
)

// CreateRequestParsers - создаются и возвращаются парсеры запросов клиента.
func CreateRequestParsers(ctx context.Context, cfg config.Config) (app.RequestParsers, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init base request parsers")

	validator, err := NewValidator(ctx, cfg)
	if err != nil {
		return app.RequestParsers{}, err
	}

	// WARNING: функция использует контекст роутера chi,
	// поэтому её можно менять только при смене самого роутера
	pathFunc := mrchi.URLPathParam

	registeredMimeTypes := mrlib.NewMimeTypeList(logger, cfg.Validation.MimeTypes)

	imageMimeTypes, err := registeredMimeTypes.MimeTypesByExts(cfg.Validation.Images.Image.File.Extensions)
	if err != nil {
		return app.RequestParsers{}, err
	}

	parsers := app.RequestParsers{
		// Bool:      mrparser.NewBool(),
		// DateTime:  mrparser.NewDateTime(),
		Int64:      mrparser.NewInt64(pathFunc),
		ItemStatus: mrparser.NewItemStatus(),
		Uint64:     mrparser.NewUint64(pathFunc),
		ListSorter: mrparser.NewListSorter(mrparser.ListSorterOptions{}),
		ListPager: mrparser.NewListPager(
			mrparser.ListPagerOptions{
				PageSizeMax:     cfg.General.PageSizeMax,
				PageSizeDefault: cfg.General.PageSizeDefault,
			},
		),
		String:    mrparser.NewString(pathFunc),
		UUID:      mrparser.NewUUID(pathFunc),
		Validator: mrparser.NewValidator(mrjson.NewDecoder(), validator),
		Image: mrparser.NewImage(
			logger,
			mrparser.WithImageMaxWidth(cfg.Validation.Images.Image.MaxWidth),
			mrparser.WithImageMaxHeight(cfg.Validation.Images.Image.MaxHeight),
			mrparser.WithImageCheckBody(cfg.Validation.Images.Image.CheckBody),
			mrparser.WithImageFileOptions(
				mrparser.WithFileMinSize(cfg.Validation.Images.Image.File.MinSize),
				mrparser.WithFileMaxSize(cfg.Validation.Images.Image.File.MaxSize),
				mrparser.WithFileMaxFiles(cfg.Validation.Images.Image.File.MaxFiles),
				mrparser.WithFileCheckRequestContentType(cfg.Validation.Images.Image.File.CheckRequestContentType),
				mrparser.WithFileAllowedMimeTypes(imageMimeTypes),
			),
		),
	}

	parsers.Parser = validate.NewParser(
		parsers.Int64,
		parsers.Uint64,
		parsers.String,
		parsers.UUID,
		parsers.Validator,
	)

	parsers.ExtendParser = validate.NewExtendParser(
		parsers.Parser,
		parsers.ItemStatus,
		parsers.ListSorter,
		parsers.ListPager,
	)

	return parsers, nil
}

// NewValidator - создаёт объект mrplayvalidator.ValidatorAdapter.
func NewValidator(ctx context.Context, _ config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init data validator")

	validator := mrplayvalidator.New()

	// registers custom tags for validation (see mrview.validator_tags.go)

	if err := validator.Register("tag_article", mrview.ValidateAnyNotSpaceSymbol); err != nil {
		return nil, err
	}

	return validator, nil
}
