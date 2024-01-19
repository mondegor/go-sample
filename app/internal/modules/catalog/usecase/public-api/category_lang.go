package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/public-api"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryLangDecorator struct {
		service CategoryService
		dict    *mrlang.MultiLangDictionary
	}
)

func NewCategoryLangDecorator(
	service CategoryService,
	dict *mrlang.MultiLangDictionary,
) *CategoryLangDecorator {
	return &CategoryLangDecorator{
		service: service,
		dict:    dict,
	}
}

func (uc *CategoryLangDecorator) GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error) {
	items, total, err := uc.service.GetList(ctx, params)

	if err != nil {
		return nil, 0, err
	}

	if dict := uc.getDict(ctx, params.LanguageID); dict != nil {
		for i := range items {
			items[i].Caption = dict.ItemByID(int(items[i].ID)).Attr("caption", items[i].Caption)
		}
	}

	return items, total, nil
}

func (uc *CategoryLangDecorator) GetItem(ctx context.Context, id mrtype.KeyInt32, languageID uint16) (*entity.Category, error) {
	item, err := uc.service.GetItem(ctx, id, languageID)

	if err != nil {
		return nil, err
	}

	item.Caption = uc.getDict(ctx, languageID).ItemByID(int(item.ID)).Attr("caption", item.Caption)

	return item, nil
}

func (uc *CategoryLangDecorator) getDict(ctx context.Context, languageID uint16) *mrlang.Dictionary {
	dict, err := uc.dict.ByLangID(languageID)

	if err != nil {
		mrctx.Logger(ctx).Warn(err)
		return nil
	}

	return dict
}
