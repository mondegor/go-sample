package factory

import (
	"context"
	"go-sample/config"
	"go-sample/internal"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrview/mrplayvalidator"

	"github.com/mondegor/go-webcore/mrview"
)

func CreateRequestParsers(ctx context.Context, cfg config.Config) (app.RequestParsers, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init base request parser")

	validator, err := NewValidator(ctx, cfg)

	if err != nil {
		return app.RequestParsers{}, err
	}

	pathFunc := mrjulienrouter.PathParam

	return app.RequestParsers{
		// Bool:      mrparser.NewBool(),
		// DateTime:  mrparser.NewDateTime(),
		Int64:      mrparser.NewInt64(pathFunc),
		ItemStatus: mrparser.NewItemStatus(),
		KeyInt32:   mrparser.NewKeyInt32(pathFunc),
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
		//File: mrparser.NewFile(
		//	mrparser.FileOptions{
		//		AllowedExts:             []string{".pdf", ".zip"},
		//		MinSize:                 512,
		//		MaxSize:                 10 * 1024 * 1024,
		//		CheckRequestContentType: true,
		//	},
		//),
		Image: mrparser.NewImage(
			mrparser.ImageOptions{
				File: mrparser.FileOptions{
					AllowedExts:             []string{".jpeg", ".jpg", ".png"},
					MinSize:                 512,
					MaxSize:                 256 * 1024,
					CheckRequestContentType: true,
				},
				MaxWidth:  1024,
				MaxHeight: 1024,
				CheckBody: true,
			},
		),
	}, nil
}

func NewValidator(ctx context.Context, cfg config.Config) (*mrplayvalidator.ValidatorAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init data validator")

	validator := mrplayvalidator.New()

	// registers custom tags for validation (see mrview.validator_tags.go)

	if err := validator.Register("tag_article", mrview.ValidateAnyNotSpaceSymbol); err != nil {
		return nil, err
	}

	return validator, nil
}
