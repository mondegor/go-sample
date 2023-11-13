package view_shared

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

func ParseFilterCategoryID(c mrcore.ClientData, key string) mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64(c.Request(), key, false)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return mrtype.KeyInt32(value)
}
