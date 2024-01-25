package usecase

import (
	"context"
	"fmt"
	module "go-sample/internal/modules/catalog"
	"go-sample/internal/modules/catalog/entity/admin-api"
	"path"
	"time"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrtool"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryImage struct {
		storage       CategoryImageStorage
		fileAPI       mrstorage.FileProviderAPI
		locker        mrcore.Locker
		eventBox      mrcore.EventBox
		serviceHelper *mrtool.ServiceHelper
	}
)

func NewCategoryImage(
	storage CategoryImageStorage,
	fileAPI mrstorage.FileProviderAPI,
	locker mrcore.Locker,
	eventBox mrcore.EventBox,
	serviceHelper *mrtool.ServiceHelper,
) *CategoryImage {
	return &CategoryImage{
		storage:       storage,
		fileAPI:       fileAPI,
		locker:        locker,
		eventBox:      eventBox,
		serviceHelper: serviceHelper,
	}
}

// GetFile - WARNING you don't forget to call item.File.Body.Close()
func (uc *CategoryImage) GetFile(ctx context.Context, categoryID mrtype.KeyInt32) (mrtype.Image, error) {
	if categoryID < 1 {
		return mrtype.Image{}, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	imageMeta, err := uc.storage.FetchMeta(ctx, categoryID)

	if err != nil {
		return mrtype.Image{}, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if imageMeta.Path == "" {
		return mrtype.Image{}, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	image, err := uc.fileAPI.DownloadFile(ctx, imageMeta.Path)

	if err != nil {
		return mrtype.Image{}, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", imageMeta)
	}

	return mrtype.Image{
		ImageInfo: mrentity.ImageMetaToInfo(imageMeta),
		Body:      image,
	}, nil
}

func (uc *CategoryImage) StoreFile(ctx context.Context, categoryID mrtype.KeyInt32, image mrtype.Image) error {
	if categoryID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	if image.OriginalName == "" || image.Size == 0 {
		return mrcore.FactoryErrServiceInvalidFile.New()
	}

	newImagePath, err := uc.getImagePath(categoryID, image.OriginalName)

	if err != nil {
		return err
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	}

	defer unlock()

	oldImageMeta, err := uc.storage.FetchMeta(ctx, categoryID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	image.Path = newImagePath

	if err = uc.fileAPI.Upload(ctx, image.ToFile()); err != nil {
		return uc.serviceHelper.WrapErrorEntityFailed(err, "FileProviderAPI", image.Path)
	}

	imageMeta := mrentity.ImageMeta{
		Path:         image.Path,
		ContentType:  image.ContentType,
		OriginalName: image.OriginalName,
		Width:        image.Width,
		Height:       image.Height,
		Size:         image.Size,
		UpdatedAt:    mrtype.TimePointer(time.Now().UTC()),
	}

	if err = uc.storage.UpdateMeta(ctx, categoryID, imageMeta); err != nil {
		uc.removeImageFile(ctx, newImagePath, oldImageMeta.Path)
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.eventBoxEmitEntity(ctx, "StoreFile", mrmsg.Data{"categoryId": categoryID, "path": newImagePath, "old-path": oldImageMeta.Path})
	uc.removeImageFile(ctx, oldImageMeta.Path, newImagePath)

	return nil
}

func (uc *CategoryImage) RemoveFile(ctx context.Context, categoryID mrtype.KeyInt32) error {
	if categoryID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	}

	defer unlock()

	imageMeta, err := uc.storage.FetchMeta(ctx, categoryID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if err = uc.storage.DeleteMeta(ctx, categoryID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.eventBoxEmitEntity(ctx, "RemoveFile", mrmsg.Data{"categoryId": categoryID, "meta": imageMeta})
	uc.removeImageFile(ctx, imageMeta.Path, "")

	return nil
}

func (uc *CategoryImage) getLockKey(categoryID mrtype.KeyInt32) string {
	return fmt.Sprintf("%s:%d", entity.ModelNameCategoryImage, categoryID)
}

func (uc *CategoryImage) getImagePath(categoryID mrtype.KeyInt32, filePath string) (string, error) {
	if ext := path.Ext(filePath); ext != "" {
		return fmt.Sprintf(
			"%s/%03x-%x%s",
			module.UnitCategoryImageDir,
			categoryID,
			time.Now().UTC().UnixNano()&0xffff,
			ext,
		), nil
	}

	return "", fmt.Errorf("file %s: ext is empty", filePath)
}

func (uc *CategoryImage) removeImageFile(ctx context.Context, filePath string, prevFilePath string) {
	if filePath == "" || filePath == prevFilePath {
		return
	}

	if err := uc.fileAPI.Remove(ctx, filePath); err != nil {
		mrctx.Logger(ctx).Err(err)
	}
}

func (uc *CategoryImage) eventBoxEmitEntity(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventBox.Emit(
		"%s::%s: %s",
		entity.ModelNameCategoryImage,
		eventName,
		data,
	)
}
