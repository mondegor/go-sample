package trademark

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/pkg/validate"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UseCaseHelper  mrcore.UseCaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		RequestParsers RequestParsers
		ResponseSender mrserver.ResponseSender

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// RequestParsers - comment struct.
	RequestParsers struct {
		Parser       *validate.Parser
		ExtendParser *validate.ExtendParser
	}
)
