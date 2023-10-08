package view

import (
    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

func ParseFilterItemStatusList(c mrcore.ClientData, items *[]mrcom_status.ItemStatus) {
    err := mrcom_status.ParseFilterItemStatusList(
        c.Request(),
        "statuses",
        mrcom_status.ItemStatusEnabled,
        items,
    )

    if err != nil {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }
}
