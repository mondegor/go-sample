package view_shared

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterTrademarkList(c mrcore.ClientData, key string) []mrentity.KeyInt32 {
    int64s, err := mrreq.ParseInt64List(c.Request(), key)

    if err != nil {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    items := make([]mrentity.KeyInt32, len(int64s))

    for i := range int64s {
        items[i] = mrentity.KeyInt32(int64s[i])
    }

    return items
}
