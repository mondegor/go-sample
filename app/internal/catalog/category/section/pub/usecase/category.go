package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrpath"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
)

type (
	// Category - comment struct.
	Category struct {
		storage      pub.CategoryStorage
		errorWrapper mrcore.UseCaseErrorWrapper
		imgBaseURL   mrpath.PathBuilder
		dict         *mrlang.MultiLangDictionary
	}
)

// NewCategory - создаёт объект Category.
func NewCategory(
	storage pub.CategoryStorage,
	errorWrapper mrcore.UseCaseErrorWrapper,
	imgBaseURL mrpath.PathBuilder,
	dict *mrlang.MultiLangDictionary,
) *Category {
	return &Category{
		storage:      storage,
		errorWrapper: errorWrapper,
		imgBaseURL:   imgBaseURL,
		dict:         dict,
	}
}

// GetList - comment method.
func (uc *Category) GetList(ctx context.Context, params entity.CategoryParams) (items []entity.Category, countItems uint64, err error) {
	items, countItems, err = uc.storage.FetchWithTotal(ctx, params)
	if err != nil {
		return nil, 0, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategory)
	}

	if countItems == 0 {
		return make([]entity.Category, 0), 0, nil
	}

	dict := uc.getDict(ctx, params.LanguageID)

	for i := range items {
		uc.prepareItem(&items[i], dict)
	}

	return items, countItems, nil
}

// GetItem - comment method.
func (uc *Category) GetItem(ctx context.Context, itemID uuid.UUID, languageID uint16) (entity.Category, error) {
	if itemID == uuid.Nil {
		return entity.Category{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	item, err := uc.storage.FetchOne(ctx, itemID)
	if err != nil {
		return entity.Category{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategory, itemID)
	}

	dict := uc.getDict(ctx, languageID)
	uc.prepareItem(&item, dict)

	return item, nil
}

func (uc *Category) getDict(ctx context.Context, languageID uint16) *mrlang.Dictionary {
	dict, err := uc.dict.ByLangID(languageID)
	if err != nil {
		mrlog.Ctx(ctx).Warn().Err(err)

		return nil
	}

	return dict
}

func (uc *Category) prepareItem(item *entity.Category, dict *mrlang.Dictionary) {
	if dict != nil {
		item.Caption = dict.ItemByKey(item.ID.String()).Attr("caption", item.Caption)
	}

	item.ImageURL = uc.imgBaseURL.BuildPath(item.ImageURL)
}
