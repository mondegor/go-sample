package trademark

import (
	"github.com/mondegor/go-sample/internal/catalog/trademark/shared/validate"

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
		RequestParser  *validate.Parser
		ResponseSender mrserver.ResponseSender

		PageSizeMax     uint64
		PageSizeDefault uint64
	}
)
