package factory

import (
	view_shared "go-sample/internal/modules/catalog/trademark/controller/http_v1/shared/view"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

type (
	Options struct {
		EventEmitter    mrsender.EventEmitter
		UsecaseHelper   *mrcore.UsecaseHelper
		PostgresAdapter *mrpostgres.ConnAdapter
		RequestParser   *view_shared.Parser
		ResponseSender  *mrresponse.Sender
	}
)
