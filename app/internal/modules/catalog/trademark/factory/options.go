package factory

import (
	view_shared "go-sample/internal/modules/catalog/trademark/controller/httpv1/shared/view"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		RequestParser  *view_shared.Parser
		ResponseSender mrserver.ResponseSender

		PageSizeMax     uint64
		PageSizeDefault uint64
	}
)
