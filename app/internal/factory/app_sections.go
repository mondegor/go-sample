package factory

import (
	"context"

	"go-sample/internal/app"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrperms"
)

const (
	sectionAdminAPICaption  = "Admin API"
	sectionAdminAPIBasePath = "/adm/"

	sectionPublicAPICaption  = "Public API"
	sectionPublicAPIBasePath = "/"
)

// NewAppSectionAdminAPI - comment func.
func NewAppSectionAdminAPI(ctx context.Context, opts app.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionAdminAPICaption,
			BasePath:     sectionAdminAPIBasePath,
			Privilege:    opts.Cfg.AppSections.AdminAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.AdminAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.AdminAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}

// NewAppSectionPublicAPI - comment func.
func NewAppSectionPublicAPI(ctx context.Context, opts app.Options) *mrperms.AppSection {
	return mrfactory.NewAppSection(
		ctx,
		mrperms.AppSectionOptions{
			Caption:      sectionPublicAPICaption,
			BasePath:     sectionPublicAPIBasePath,
			Privilege:    opts.Cfg.AppSections.PublicAPI.Privilege,
			AuthSecret:   opts.Cfg.AppSections.PublicAPI.Auth.Secret,
			AuthAudience: opts.Cfg.AppSections.PublicAPI.Auth.Audience,
		},
		opts.AccessControl,
	)
}
