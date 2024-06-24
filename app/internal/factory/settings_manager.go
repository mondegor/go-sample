package factory

import (
	"context"

	"go-sample/internal/app"

	mrsettingsfactory "github.com/mondegor/go-components/mrsettings/factory"
	"github.com/mondegor/go-components/mrsettings/lightgetter"
	"github.com/mondegor/go-components/mrsettings/setter"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker/mrschedule"
)

const (
	settingsManagerTableName  = "sample_catalog.settings"
	settingsManagerPrimaryKey = "setting_id"
)

// NewSettingsGetter - comment func.
func NewSettingsGetter(ctx context.Context, opts app.Options) (*lightgetter.Component, *mrschedule.TaskShell) {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings getter")

	getter := mrsettingsfactory.NewComponentCacheGetter(
		opts.PostgresConnManager,
		mrsql.NewEntityMeta(settingsManagerTableName, settingsManagerPrimaryKey, nil),
		opts.UsecaseErrorWrapper,
		mrsettingsfactory.ComponentCacheGetterOptions{},
	)

	task := mrschedule.NewTaskShell(
		opts.Cfg.TaskSchedule.SettingsReloader.Caption,
		opts.Cfg.TaskSchedule.SettingsReloader.Startup,
		opts.Cfg.TaskSchedule.SettingsReloader.Period,
		opts.Cfg.TaskSchedule.SettingsReloader.Timeout,
		func(ctx context.Context) error {
			count, err := getter.Reload(ctx)
			if err != nil {
				return err
			}

			if count > 0 {
				mrlog.Ctx(ctx).Info().Msgf("Settings are reloaded: %d", count)
			}

			return nil
		},
	)

	return lightgetter.New(getter), task
}

// NewSettingsSetter - comment func.
func NewSettingsSetter(ctx context.Context, opts app.Options) *setter.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings setter")

	return mrsettingsfactory.NewComponentSetter(
		opts.PostgresConnManager,
		mrsql.NewEntityMeta(settingsManagerTableName, settingsManagerPrimaryKey, nil),
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
		mrsettingsfactory.ComponentSetterOptions{},
	)
}
