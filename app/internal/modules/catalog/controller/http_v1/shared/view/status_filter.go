package view_shared

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterStatusList(c mrcore.ClientData, key string) []mrenum.ItemStatus {
	items, err := mrreq.ParseItemStatusList(
		c.Request(),
		key,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return items
}
