package factory

import (
	"context"
	"go-sample/internal"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrperms"
)

const (
	sectionAdminAPICaption  = "Admin API"
	sectionAdminAPIRootPath = "/adm/"

	sectionPublicAPICaption  = "Public API"
	sectionPublicAPIRootPath = "/"
)

func NewAppSectionAdminAPI(ctx context.Context, opts app.Options) mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionAdminAPICaption,
			RootPath:     sectionAdminAPIRootPath,
			Privilege:    opts.Cfg.AppSections.AdminAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.AdminAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.AdminAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}

func NewAppSectionPublicAPI(ctx context.Context, opts app.Options) mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionPublicAPICaption,
			RootPath:     sectionPublicAPIRootPath,
			Privilege:    opts.Cfg.AppSections.PublicAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.PublicAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.PublicAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}
