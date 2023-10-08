package view

import (
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterTrademarkList(c mrcore.ClientData, items *[]mrentity.KeyInt32) {
    int64s, err := mrreq.Int64List(c.Request(), "trademarks")

    if err != nil {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    for _, item := range int64s {
        *items = append(*items, mrentity.KeyInt32(item))
    }
}
