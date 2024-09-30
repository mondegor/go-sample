package category

import (
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrpath"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/go-sample/internal/catalog/category/shared/validate"
)

type (
	// Options - comment struct.
	Options struct {
		EventEmitter   mrsender.EventEmitter
		UseCaseHelper  mrcore.UseCaseErrorWrapper
		DBConnManager  mrstorage.DBConnManager
		Locker         mrlock.Locker
		RequestParsers RequestParsers
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

	// RequestParsers - comment struct.
	RequestParsers struct {
		// Parser       *pkgvalidate.Parser
		// ExtendParser *pkgvalidate.ExtendParser
		ModuleParser *validate.CategoryParser
	}
)
