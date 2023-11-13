package usecase

import (
    "context"
    "fmt"
    "go-sample/internal/modules/catalog/entity/admin-api"
    "path/filepath"
    "time"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrtool"
    "github.com/mondegor/go-webcore/mrtype"
)

type (
    CategoryImage struct {
        baseImageURL string
        storage CategoryImageStorage
        fileAPI mrstorage.FileProviderAPI
        locker mrcore.Locker
        eventBox mrcore.EventBox
        serviceHelper *mrtool.ServiceHelper
    }
)

func NewCategoryImage(
    baseImageURL string,
    storage CategoryImageStorage,
    fileAPI mrstorage.FileProviderAPI,
    locker mrcore.Locker,
    eventBox mrcore.EventBox,
    serviceHelper *mrtool.ServiceHelper,
) *CategoryImage {
    return &CategoryImage{
        baseImageURL: baseImageURL,
        storage: storage,
        fileAPI: fileAPI,
        locker: locker,
        eventBox: eventBox,
        serviceHelper: serviceHelper,
    }
}

// Get - WARNING you don't forget to call item.File.Body.Close()
func (uc *CategoryImage) Get(ctx context.Context, categoryID mrtype.KeyInt32) (*mrtype.File, error) {
    if categoryID < 1 {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryID})
    }

    imagePath, err := uc.storage.FetchPath(ctx, categoryID)

    if err != nil {
        return nil, uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    if imagePath == "" {
        return nil, mrcore.FactoryErrServiceEntityNotFound.New(categoryID)
    }

    file, err := uc.fileAPI.Download(ctx, imagePath)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, "FileProviderAPI")
    }

    return file, nil
}

func (uc *CategoryImage) GetInfoByPath(ctx context.Context, imagePath string) (*mrtype.FileInfo, error) {
    if imagePath == "" {
        return nil, mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"imagePath": "EMPTY"})
    }

    info, err := uc.fileAPI.Info(ctx, imagePath)

    if err != nil {
        return nil, mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, "FileProviderAPI")
    }

    return &info, nil
}

func (uc *CategoryImage) Store(ctx context.Context, categoryID mrtype.KeyInt32, file *mrtype.File) error {
    if categoryID < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryID})
    }

    newImagePath, err := uc.getImagePath(categoryID, file.OriginalName)

    if err != nil {
        return err
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

    if err != nil {
        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategoryImage)
    }

    defer unlock()

    oldImagePath, err := uc.storage.FetchPath(ctx, categoryID)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    file.Path = newImagePath

    if err = uc.fileAPI.Upload(ctx, file); err != nil {
        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, "FileProviderAPI")
    }

    if err = uc.storage.Update(ctx, categoryID, newImagePath); err != nil {
        uc.removeImageFile(ctx, newImagePath)
        return uc.serviceHelper.WrapErrorForUpdate(err, entity.ModelNameCatalogCategoryImage)
    }

    uc.eventBox.Emit(
        "%s::Upload: cid=%d; path=%s; old-path=%s",
        entity.ModelNameCatalogCategoryImage,
        categoryID,
        newImagePath,
        oldImagePath,
    )

    uc.removeImageFile(ctx, oldImagePath)

    return nil
}

func (uc *CategoryImage) Remove(ctx context.Context, categoryID mrtype.KeyInt32) error {
    if categoryID < 1 {
        return mrcore.FactoryErrServiceIncorrectInputData.New(mrerr.Arg{"categoryId": categoryID})
    }

    unlock, err := uc.locker.Lock(ctx, uc.getLockKey(categoryID))

    if err != nil {
        return mrcore.FactoryErrServiceTemporarilyUnavailable.Wrap(err, entity.ModelNameCatalogCategoryImage)
    }

    defer unlock()

    imagePath, err := uc.storage.FetchPath(ctx, categoryID)

    if err != nil {
        return uc.serviceHelper.WrapErrorForSelect(err, entity.ModelNameCatalogCategoryImage)
    }

    if err = uc.storage.Delete(ctx, categoryID); err != nil {
        return uc.serviceHelper.WrapErrorForRemove(err, entity.ModelNameCatalogCategoryImage)
    }

    uc.eventBox.Emit(
        "%s::Remove: cid=%d; path=%s",
        entity.ModelNameCatalogCategoryImage,
        categoryID,
        imagePath,
    )

    uc.removeImageFile(ctx, imagePath)

    return nil
}

func (uc *CategoryImage) getLockKey(categoryID mrtype.KeyInt32) string {
    return fmt.Sprintf("%s:%d", entity.ModelNameCatalogCategoryImage, categoryID)
}

func (uc *CategoryImage) getImagePath(categoryID mrtype.KeyInt32, path string) (string, error) {
    ext := filepath.Ext(path)

    if ext == "" {
        return "", fmt.Errorf("file %s: ext is empty", path)
    }

    return fmt.Sprintf(
        "%s/%03x-%x%s",
        uc.baseImageURL,
        categoryID,
        time.Now().UnixNano() & 0xffff,
        ext,
    ), nil
}

func (uc *CategoryImage) removeImageFile(ctx context.Context, filePath string) {
    if filePath == "" {
        return
    }

    err := uc.fileAPI.Remove(ctx, filePath)

    if err != nil {
        mrctx.Logger(ctx).Err(err)
    }
}
