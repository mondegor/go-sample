package view_shared

import (
	"go-sample/internal/global"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

func ParseFilterString(c mrcore.ClientData, key string) string {
	value, err := mrreq.ParseStr(c.Request(), key, false)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return value
}

func ParseFilterRangeInt64(c mrcore.ClientData, key string) mrtype.RangeInt64 {
	str, err := mrreq.ParseRangeInt64(c.Request(), key)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return str
}

func ParseListSorter(c mrcore.ClientData, sorter mrview.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		c.Request(),
		global.ParamNameSortField,
		global.ParamNameSortDirection,
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return value
}

func ParseListPager(c mrcore.ClientData) mrtype.PageParams {
	value, err := mrreq.ParsePageParams(
		c.Request(),
		global.ParamNamePageIndex,
		global.ParamNamePageSize,
	)

	if err != nil || value.Size > global.PageSizeMax {
		if err != nil {
			mrctx.Logger(c.Context()).Warn(err)
		}

		return mrtype.PageParams{
			Size: global.PageSizeDefault,
		}
	}

	return value
}
