package factory

import (
	view_shared "go-sample/internal/modules/catalog/category/controller/httpv1/shared/view"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UsecaseHelper  mrcore.UsecaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		Locker         mrlock.Locker
		RequestParser  *view_shared.Parser
		ResponseSender mrserver.FileResponseSender

		UnitCategory UnitCategoryOptions

		PageSizeMax     uint64
		PageSizeDefault uint64
	}

	// UnitCategoryOptions - comment struct.
	UnitCategoryOptions struct {
		Dictionary      *mrlang.MultiLangDictionary
		ImageFileAPI    mrstorage.FileProviderAPI
		ImageURLBuilder mrpath.PathBuilder
	}
)
