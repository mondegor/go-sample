package factory

import (
	"go-sample/internal/app"

	"github.com/mondegor/go-sysmess/mrcaller"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/features"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/mrparser/mrparserinit"
	"github.com/mondegor/go-webcore/mrserver/mrreq/mrreqinit"
)

// NewErrorManager - comment func.
func NewErrorManager(opts app.Options) *mrinit.ErrorManager {
	extra := mrerr.ProtoExtra{}

	// create Caller for Errors
	if opts.Cfg.Debugging.ErrorCaller.Enable {
		caller := mrcaller.New(
			mrcaller.DepthOption(opts.Cfg.Debugging.ErrorCaller.Depth),
			mrcaller.ShowFuncNameOption(opts.Cfg.Debugging.ErrorCaller.ShowFuncName),
			mrcaller.FilterStackTraceOption(
				mrcaller.FilterStackTraceTrimUpper(opts.Cfg.Debugging.ErrorCaller.UpperBounds),
			),
		)

		extra = mrerr.ProtoExtra{
			Caller: func() mrerr.StackTracer {
				return caller.StackTrace()
			},
			OnCreated: func(_ *mrerr.AppError) (instanceID string) {
				return features.GenerateInstanceID()
			},
		}
	}

	if opts.Sentry != nil {
		extra.OnCreated = func(err *mrerr.AppError) (instanceID string) {
			if instanceID = opts.Sentry.CaptureAppError(err); instanceID != "" {
				return instanceID
			}

			return features.GenerateInstanceID()
		}
	}

	manager := mrinit.NewErrorManager(extra)

	manager.Register(
		mrinit.ManagedError{
			Err:           mrerr.ErrErrorIsNilPointer,
			WithCaller:    true,
			WithOnCreated: true,
		},
	)

	manager.RegisterList(mrinit.ManagedInternalErrors())
	manager.RegisterList(mrinit.ManagedStorageErrors())
	manager.RegisterList(mrinit.ManagedUseCaseErrors())
	manager.RegisterList(mrinit.ManagedHttpErrors())
	manager.RegisterList(mrreqinit.ManagedHttpErrors())
	manager.RegisterList(mrparserinit.ManagedHttpErrors())

	return manager
}
