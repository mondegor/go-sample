package usecase

import (
    "context"
    "fmt"
    "go-sample/internal/entity/admin-panel"
    "path/filepath"
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrtool"
)

type (
    CatalogCategoryImage struct {
        baseImageUrl string
        storage CatalogCategoryImageStorage
        storageFiles mrstorage.FileProvider
        locker mrcore.Locker
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewCatalogCategoryImage(
    baseImageUrl string,
    storage CatalogCategoryImageStorage,
    storageFiles mrstorage.FileProvider,
    locker mrcore.Locker,
    eventBox mrcore.EventBox,
    serviceHelper *mrtool.ServiceHelper,
) *CatalogCategoryImage {
    return &CatalogCategoryImage{
        baseImageUrl: baseImageUrl,
        storage: storage,
        storageFiles: storageFiles,
        locker: locker,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
    }
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *CatalogCategoryImage) Get(ctx context.Context, categoryId mrentity.KeyInt32) (*entity.CatalogCategoryImageObject, error) {
    if categoryId < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryId})
    }

    imagePath, err := uc.storage.FetchOne(ctx, categoryId)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    if imagePath == "" {
        return nil, mrcore.FactoryErrServiceEntityNotFound.New(categoryId)
    }

    item := entity.CatalogCategoryImageObject{
        CategoryId: categoryId,
        File: mrentity.File{
            Name: imagePath,
        },
    }

    err = uc.storageFiles.Download(ctx, &item.File)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, mrentity.ModelNameFile)
    }

    return &item, nil
}

func (uc *CatalogCategoryImage) Store(ctx context.Context, item *entity.CatalogCategoryImageObject) error {
    if item.CategoryId < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"image.CategoryId": item.CategoryId})
    }

    newImagePath, err := uc.getImagePath(item)

    if err != nil {
        return err
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(item.CategoryId))

    if err != nil {
        return err
    }

    defer unlock()

    oldImagePath, err := uc.storage.FetchOne(ctx, item.CategoryId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    file := item.File
    file.Name = newImagePath

    err = uc.storageFiles.Upload(ctx, &file)

    if err != nil {
        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, mrentity.ModelNameFile)
    }

    err = uc.storage.Update(ctx, item.CategoryId, newImagePath)

    if err != nil {
        uc.removeFile(ctx, newImagePath)
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogCategoryImage)
    }

    uc.eventBox.Emit(
        "%s::Upload: cid=%d; path=%s",
        entity.ModelNameCatalogCategoryImage,
        item.CategoryId,
        newImagePath,
    )

    uc.removeFile(ctx, oldImagePath)

    return nil
}

func (uc *CatalogCategoryImage) Remove(ctx context.Context, categoryId mrentity.KeyInt32) error {
    if categoryId < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryId})
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryId))

    if err != nil {
        return err
    }

    defer unlock()

    oldImagePath, err := uc.storage.FetchOne(ctx, categoryId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    err = uc.storage.Delete(ctx, categoryId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogCategoryImage)
    }

    uc.eventBox.Emit(
        "%s::Remove: cid=%d; path=%s",
        entity.ModelNameCatalogCategoryImage,
        categoryId,
        oldImagePath,
    )

    uc.removeFile(ctx, oldImagePath)

    return nil
}

func (uc *CatalogCategoryImage) getLockKey(categoryId mrentity.KeyInt32) string {
    return fmt.Sprintf("%s:%d", entity.ModelNameCatalogCategoryImage, categoryId)
}

func (uc *CatalogCategoryImage) getImagePath(item *entity.CatalogCategoryImageObject) (string, error) {
    ext := filepath.Ext(item.File.Name)

    if ext == "" {
        return "", fmt.Errorf("ext is empty")
    }

    return fmt.Sprintf(
        "%s/%03x-%x%s",
        uc.baseImageUrl,
        item.CategoryId,
        time.Now().UnixNano() & 0xffff,
        ext,
    ), nil
}

func (uc *CatalogCategoryImage) removeFile(ctx context.Context, filePath string) {
    if filePath == "" {
        return
    }

    err := uc.storageFiles.Remove(ctx, filePath)

    if err != nil {
        mrctx.Logger(ctx).Err(err)
    }
}
