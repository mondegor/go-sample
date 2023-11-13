package view_shared

import (
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
    "github.com/mondegor/go-webcore/mrtype"
)

func ParseFilterTrademarkList(c mrcore.ClientData, key string) []mrtype.KeyInt32 {
    int64s, err := mrreq.ParseInt64List(c.Request(), key)

    if err != nil {
        mrctx.Logger(c.Context()).Warn(err)
    }

    items := make([]mrtype.KeyInt32, len(int64s))

    for i := range int64s {
        items[i] = mrtype.KeyInt32(int64s[i])
    }

    return items
}
