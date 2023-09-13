package http_v1

import (
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func parseFilterTrademarks(c mrcore.ClientData, trademarks *[]mrentity.KeyInt32) {
    items, err := mrreq.Int64List(c.Request(), "trademarks")

    if err == nil {
        for _, item := range items {
            *trademarks = append(*trademarks, mrentity.KeyInt32(item))
        }
    } else {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }
}

func parseFilterStatuses(c mrcore.ClientData, statuses *[]mrcom.ItemStatus) {
    items, err := mrreq.EnumList(c.Request(), "statuses")

    if err == nil {
        var itemStatus mrcom.ItemStatus

        for _, item := range items {
            if item == mrcom.ItemStatusRemoved.String() {
                continue
            }

            err = itemStatus.ParseAndSet(item)

            if err != nil {
                mrctx.Logger(c.Context()).Warn(err.Error())
                continue
            }

            *statuses = append(*statuses, itemStatus)
        }
    } else {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    if len(*statuses) == 0 {
        *statuses = append(*statuses, mrcom.ItemStatusEnabled)
    }
}
