package http_v1

import (
	"context"
	"fmt"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/controller/http_v1/admin-api/view"
	view_shared "go-sample/internal/modules/catalog/controller/http_v1/shared/view"
	"go-sample/internal/modules/catalog/entity/admin-api"
	usecase "go-sample/internal/modules/catalog/usecase/admin-api"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	categoryURL             = "/v1/catalog/categories"
	categoryItemURL         = "/v1/catalog/categories/:id"
	categoryChangeStatusURL = "/v1/catalog/categories/:id/status"
	categoryItemImageURL    = "/v1/catalog/categories/:id/image"
)

type (
	Category struct {
		section      mrcore.ClientSection
		service      usecase.CategoryService
		serviceImage usecase.CategoryImageService
		listSorter   mrview.ListSorter
	}
)

func NewCategory(
	section mrcore.ClientSection,
	service usecase.CategoryService,
	serviceImage usecase.CategoryImageService,
	listSorter mrview.ListSorter,
) *Category {
	return &Category{
		section:      section,
		service:      service,
		serviceImage: serviceImage,
		listSorter:   listSorter,
	}
}

func (ht *Category) AddHandlers(router mrcore.HttpRouter) {
	moduleAccessFunc := func(next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
		return ht.section.MiddlewareWithPermission(module.PermissionCatalogCategory, next)
	}

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryURL), moduleAccessFunc(ht.GetList()))
	router.HttpHandlerFunc(http.MethodPost, ht.section.Path(categoryURL), moduleAccessFunc(ht.Create()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Get()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Store()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(categoryItemURL), moduleAccessFunc(ht.Remove()))

	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryChangeStatusURL), moduleAccessFunc(ht.ChangeStatus()))

	router.HttpHandlerFunc(http.MethodGet, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.GetImage()))
	router.HttpHandlerFunc(http.MethodPut, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.UploadImage()))
	router.HttpHandlerFunc(http.MethodDelete, ht.section.Path(categoryItemImageURL), moduleAccessFunc(ht.RemoveImage()))
}

func (ht *Category) GetList() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		items, totalItems, err := ht.service.GetList(c.Context(), ht.listParams(c))

		if err != nil {
			return err
		}

		for i := range items {
			items[i].ImageInfo = ht.getImageInfo(c.Context(), items[i].ImagePath)
		}

		return c.SendResponse(
			http.StatusOK,
			view.CategoryListResponse{
				Items: items,
				Total: totalItems,
			},
		)
	}
}

func (ht *Category) listParams(c mrcore.ClientContext) entity.CategoryParams {
	return entity.CategoryParams{
		Filter: entity.CategoryListFilter{
			SearchText: view_shared.ParseFilterString(c, module.ParamNameFilterSearchText),
			Statuses:   view_shared.ParseFilterStatusList(c, module.ParamNameFilterStatuses),
		},
		Sorter: view_shared.ParseSortParams(c, ht.listSorter),
		Pager:  view_shared.ParsePageParams(c),
	}
}

func (ht *Category) Get() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.service.GetItem(c.Context(), ht.getItemID(c))

		if err != nil {
			return err
		}

		item.ImageInfo = ht.getImageInfo(c.Context(), item.ImagePath)

		return c.SendResponse(http.StatusOK, item)
	}
}

func (ht *Category) Create() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.CreateCategoryRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			Caption: request.Caption,
		}

		if err := ht.service.Create(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponse(
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: fmt.Sprintf("%d", item.ID),
				Message: mrctx.Locale(c.Context()).TranslateMessage(
					"msgCategorySuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *Category) Store() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.StoreCategoryRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			ID:         ht.getItemID(c),
			TagVersion: request.Version,
			Caption:    request.Caption,
		}

		if err := ht.service.Store(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) ChangeStatus() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		request := view.ChangeItemStatusRequest{}

		if err := c.Validate(&request); err != nil {
			return err
		}

		item := entity.Category{
			ID:         ht.getItemID(c),
			TagVersion: request.TagVersion,
			Status:     request.Status,
		}

		if err := ht.service.ChangeStatus(c.Context(), &item); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) Remove() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.service.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) GetImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		item, err := ht.serviceImage.Get(c.Context(), ht.getItemID(c))

		if err != nil {
			return err
		}

		defer item.Body.Close()

		return c.SendFile(item.FileInfo, "", item.Body)
	}
}

func (ht *Category) UploadImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		logger := mrctx.Logger(c.Context())

		file, hdr, err := c.Request().FormFile(module.ParamNameFileCatalogCategoryImage)

		if err != nil {
			mrdebug.MultipartForm(logger, c.Request().MultipartForm)
			return mrcore.FactoryErrHttpMultipartFormFile.Caller(-1).Wrap(err, module.ParamNameFileCatalogCategoryImage)
		}

		defer file.Close()

		mrdebug.MultipartFileHeader(logger, hdr)

		item := mrtype.File{
			FileInfo: mrtype.FileInfo{
				ContentType:  hdr.Header.Get("Content-Type"),
				OriginalName: hdr.Filename,
				Size:         hdr.Size,
			},
			Body: file,
		}

		if err = ht.serviceImage.Store(c.Context(), ht.getItemID(c), &item); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) RemoveImage() mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		if err := ht.serviceImage.Remove(c.Context(), ht.getItemID(c)); err != nil {
			return err
		}

		return c.SendResponseNoContent()
	}
}

func (ht *Category) getItemID(c mrcore.ClientContext) mrtype.KeyInt32 {
	return view_shared.ParseIDFromPath(c, "id")
}

func (ht *Category) getImageInfo(ctx context.Context, imagePath string) *mrtype.FileInfo {
	if imagePath == "" {
		return nil
	}

	info, err := ht.serviceImage.GetInfoByPath(ctx, imagePath)

	if err != nil {
		mrctx.Logger(ctx).Err(err)
		return nil
	}

	return info
}
