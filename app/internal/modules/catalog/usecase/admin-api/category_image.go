package usecase

import (
	"context"
	"fmt"
	"go-sample/internal/modules/catalog/entity/admin-api"
	"path/filepath"
	"time"

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

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *CategoryImage) Get(ctx context.Context, categoryID mrtype.KeyInt32) (*mrtype.File, error) {
	if categoryID < 1 {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	imagePath, err := uc.storage.FetchPath(ctx, categoryID)

	if err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if imagePath == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	file, err := uc.fileAPI.Download(ctx, imagePath)

	if err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", imagePath)
	}

	return file, nil
}

func (uc *CategoryImage) GetInfoByPath(ctx context.Context, imagePath string) (*mrtype.FileInfo, error) {
	if imagePath == "" {
		return nil, mrcore.FactoryErrServiceEntityNotFound.New()
	}

	info, err := uc.fileAPI.Info(ctx, imagePath)

	if err != nil {
		return nil, uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, "FileProviderAPI", imagePath)
	}

	return &info, nil
}

func (uc *CategoryImage) Store(ctx context.Context, categoryID mrtype.KeyInt32, file *mrtype.File) error {
	if categoryID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	newImagePath, err := uc.getImagePath(categoryID, file.OriginalName)

	if err != nil {
		return err
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	}

	defer unlock()

	oldImagePath, err := uc.storage.FetchPath(ctx, categoryID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	file.Path = newImagePath

	if err = uc.fileAPI.Upload(ctx, file); err != nil {
		return uc.serviceHelper.WrapErrorEntityFailed(err, "FileProviderAPI", file.Path)
	}

	if err = uc.storage.Update(ctx, categoryID, newImagePath); err != nil {
		uc.removeImageFile(ctx, newImagePath)
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.eventBoxEmitEntity(ctx, "Store", mrmsg.Data{"categoryId": categoryID, "path": newImagePath, "old-path": oldImagePath})

	uc.removeImageFile(ctx, oldImagePath)

	return nil
}

func (uc *CategoryImage) Remove(ctx context.Context, categoryID mrtype.KeyInt32) error {
	if categoryID < 1 {
		return mrcore.FactoryErrServiceEntityNotFound.New()
	}

	unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

	if err != nil {
		return uc.serviceHelper.WrapErrorFailed(err, entity.ModelNameCategoryImage)
	}

	defer unlock()

	oldImagePath, err := uc.storage.FetchPath(ctx, categoryID)

	if err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	if err = uc.storage.Delete(ctx, categoryID); err != nil {
		return uc.serviceHelper.WrapErrorEntityNotFoundOrFailed(err, entity.ModelNameCategoryImage, categoryID)
	}

	uc.eventBoxEmitEntity(ctx, "Remove", mrmsg.Data{"categoryId": categoryID, "old-path": oldImagePath})

	uc.removeImageFile(ctx, oldImagePath)

	return nil
}

func (uc *CategoryImage) getLockKey(categoryID mrtype.KeyInt32) string {
	return fmt.Sprintf("%s:%d", entity.ModelNameCategoryImage, categoryID)
}

func (uc *CategoryImage) getImagePath(categoryID mrtype.KeyInt32, path string) (string, error) {
	ext := filepath.Ext(path)

	if ext == "" {
		return "", fmt.Errorf("file %s: ext is empty", path)
	}

	return fmt.Sprintf(
		"%03x-%x%s",
		categoryID,
		time.Now().UnixNano()&0xffff,
		ext,
	), nil
}

func (uc *CategoryImage) removeImageFile(ctx context.Context, filePath string) {
	if filePath == "" {
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
