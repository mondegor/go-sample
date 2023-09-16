package usecase

import (
    "context"
    "fmt"
    "go-sample/internal/entity"
    "path/filepath"
    "time"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrtool"
)

const (
    CatalogCategoryImageDir = "catalog/categories"
)

type (
	CatalogCategoryImage struct {
        storage CatalogCategoryImageStorage
        storageFiles mrstorage.FileProvider
        locker mrcore.Locker
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewCatalogCategoryImage(storage CatalogCategoryImageStorage,
                             storageFiles mrstorage.FileProvider,
                             locker mrcore.Locker,
                             eventBox mrcore.EventBox,
                             serviceHelper *mrtool.ServiceHelper) *CatalogCategoryImage {
    return &CatalogCategoryImage{
        storage: storage,
        storageFiles: storageFiles,
        locker: locker,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
    }
}

// Load - WARNING you don't forget to call item.File.Body.Close()
func (uc *CatalogCategoryImage) Load(ctx context.Context, item *entity.CatalogCategoryImageObject) error {
    if item.CategoryId < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"image.CategoryId": item.CategoryId})
    }

    imagePath, err := uc.storage.Fetch(ctx, item.CategoryId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    if imagePath == "" {
        return mrcore.FactoryErrServiceEntityNotFound.New(item.CategoryId)
    }

    item.File.Name = imagePath
    err = uc.storageFiles.Download(ctx, &item.File)

    if err != nil {
        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, mrstorage.ModelNameFile)
    }

    return nil
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

    oldImagePath, err := uc.storage.Fetch(ctx, item.CategoryId)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    file := item.File
    file.Name = newImagePath

    err = uc.storageFiles.Upload(ctx, &file)

    if err != nil {
        return mrcore.FactoryErrServiceEntityTemporarilyUnavailable.Wrap(err, mrstorage.ModelNameFile)
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
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"image.CategoryId": categoryId})
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryId))

    if err != nil {
        return err
    }

    defer unlock()

    oldImagePath, err := uc.storage.Fetch(ctx, categoryId)

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
        CatalogCategoryImageDir,
        item.CategoryId,
        time.Now().UnixNano() & 0xffffffff,
        ext,
    ), nil
}

func (uc *CatalogCategoryImage) removeFile(ctx context.Context, filePath string) {
    // :TODO: поставить в очередь для удаления (уменьшает время ожидания, а также позволяет не блокировать читателей)
    err := uc.storageFiles.Remove(ctx, filePath)

    if err != nil {
        mrctx.Logger(ctx).Err(err)
    }
}
