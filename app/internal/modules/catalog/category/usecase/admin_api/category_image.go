package usecase

import (
	"context"
	"fmt"
	"path"
	"time"

	entity "go-sample/internal/modules/catalog/category/entity/admin_api"
	"go-sample/internal/modules/catalog/category/module"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CategoryImage - comment struct.
	CategoryImage struct {
		storage      CategoryImageStorage
		fileAPI      mrstorage.FileProviderAPI
		locker       mrlock.Locker
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewCategoryImage - comment func.
func NewCategoryImage(
	storage CategoryImageStorage,
	fileAPI mrstorage.FileProviderAPI,
	locker mrlock.Locker,
	eventEmitter mrsender.EventEmitter,
	errorWrapper mrcore.UsecaseErrorWrapper,
) *CategoryImage {
	return &CategoryImage{
		storage:      storage,
		fileAPI:      fileAPI,
		locker:       locker,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// GetFile - comment method.
// WARNING you don't forget to call item.File.Body.Close().
func (uc *CategoryImage) GetFile(ctx context.Context, categoryID uuid.UUID) (mrtype.Image, error) {
	if categoryID == uuid.Nil {
		return mrtype.Image{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	imageMeta, err := uc.storage.FetchMeta(ctx, categoryID)
	if err != nil {
		return mrtype.Image{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if imageMeta.Path == "" {
		return mrtype.Image{}, mrcore.ErrUseCaseEntityNotFound.New()
	}

	image, err := uc.fileAPI.DownloadFile(ctx, imageMeta.Path)
	if err != nil {
		return mrtype.Image{}, uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", imageMeta)
	}

	return mrtype.Image{
		ImageInfo: mrentity.ImageMetaToInfo(imageMeta, nil),
		Body:      image,
	}, nil
}

// StoreFile - comment method.
func (uc *CategoryImage) StoreFile(ctx context.Context, categoryID uuid.UUID, image mrtype.Image) error {
	if categoryID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if image.OriginalName == "" || image.Size == 0 {
		return mrcore.ErrUseCaseInvalidFile.New()
	}

	newImagePath, err := uc.getImagePath(categoryID, image.OriginalName)
	if err != nil {
		return err
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	} else {
		defer unlock()
	}

	oldImageMeta, err := uc.storage.FetchMeta(ctx, categoryID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	image.Path = newImagePath

	if err = uc.fileAPI.Upload(ctx, image.ToFile()); err != nil {
		return uc.errorWrapper.WrapErrorEntityFailed(err, "FileProviderAPI", image.Path)
	}

	imageMeta := mrentity.ImageMeta{
		Path:         image.Path,
		ContentType:  image.ContentType,
		OriginalName: image.OriginalName,
		Width:        image.Width,
		Height:       image.Height,
		Size:         image.Size,
		UpdatedAt:    mrtype.TimeToPointer(time.Now().UTC()),
	}

	if err = uc.storage.UpdateMeta(ctx, categoryID, imageMeta); err != nil {
		uc.removeImageFile(ctx, newImagePath, oldImageMeta.Path)

		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.emitEvent(ctx, "StoreFile", mrmsg.Data{"categoryId": categoryID, "path": newImagePath, "old-path": oldImageMeta.Path})
	uc.removeImageFile(ctx, oldImageMeta.Path, newImagePath)

	return nil
}

// RemoveFile - comment method.
func (uc *CategoryImage) RemoveFile(ctx context.Context, categoryID uuid.UUID) error {
	if categoryID == uuid.Nil {
		return mrcore.ErrUseCaseEntityNotFound.New()
	}

	if unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID)); err != nil {
		return uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	} else {
		defer unlock()
	}

	imageMeta, err := uc.storage.FetchMeta(ctx, categoryID)
	if err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if err = uc.storage.DeleteMeta(ctx, categoryID); err != nil {
		return uc.errorWrapper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.emitEvent(ctx, "RemoveFile", mrmsg.Data{"categoryId": categoryID, "meta": imageMeta})
	uc.removeImageFile(ctx, imageMeta.Path, "")

	return nil
}

func (uc *CategoryImage) getLockKey(categoryID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", entity.ModelNameCategoryImage, categoryID)
}

func (uc *CategoryImage) getImagePath(categoryID uuid.UUID, filePath string) (string, error) {
	if ext := path.Ext(filePath); ext != "" {
		return fmt.Sprintf(
			"%s/%s-%x%s",
			module.ImageDir,
			categoryID,
			time.Now().UTC().UnixNano()&0xffff,
			ext,
		), nil
	}

	return "", fmt.Errorf("file %s: ext is empty", filePath)
}

func (uc *CategoryImage) removeImageFile(ctx context.Context, filePath, prevFilePath string) {
	if filePath == "" || filePath == prevFilePath {
		return
	}

	if err := uc.fileAPI.Remove(ctx, filePath); err != nil {
		mrlog.Ctx(ctx).Error().Err(err)
	}
}

func (uc *CategoryImage) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCategoryImage,
		data,
	)
}
