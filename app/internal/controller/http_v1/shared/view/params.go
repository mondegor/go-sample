package view_shared

import (
    "go-sample/internal/controller"

    "github.com/mondegor/go-storage/mrentity"
    storage_mrreq "github.com/mondegor/go-storage/mrreq"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterString(c mrcore.ClientData, key string) string {
    str, err := mrreq.ParseStr(c.Request(), key, false)

    if err != nil {
        mrcore.LogErr(err)
    }

    return str
}

func ParseFilterRangeInt64(c mrcore.ClientData, key string) mrentity.RangeInt64 {
    str, err := storage_mrreq.ParseRangeInt64(c.Request(), key)

    if err != nil {
        mrcore.LogErr(err)
    }

    return str
}

func ParseListSorter(c mrcore.ClientData) mrentity.ListSorter {
    ls, err := storage_mrreq.ParseListSorter(
        c.Request(),
        controller.ParamNameSortField,
        controller.ParamNameSortDirection,
    )

    if err != nil {
        mrcore.LogErr(err)
        return mrentity.ListSorter{}
    }

    return ls
}

func ParseListPager(c mrcore.ClientData) mrentity.ListPager {
    lp, err := storage_mrreq.ParseListPager(
        c.Request(),
        controller.ParamNamePageIndex,
        controller.ParamNamePageSize,
    )

    if err != nil || lp.Size > controller.PageSizeMax {
        if err != nil {
            mrcore.LogErr(err)
        }

        return mrentity.ListPager{Size: controller.PageSizeDefault}
    }

    return lp
}
