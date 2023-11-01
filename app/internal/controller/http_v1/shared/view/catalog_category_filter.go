package view_shared

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterCategoryId(c mrcore.ClientData, key string) mrentity.KeyInt32 {
    value, err := mrreq.ParseInt64(c.Request(), key, false)

    if err != nil {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    return mrentity.KeyInt32(value)
}
