package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrperms"

	"github.com/mondegor/go-sample/internal/app"
)

const (
	sectionAdminAPICaption  = "Admin API"
	sectionAdminAPIBasePath = "/adm/"

	sectionPublicAPICaption  = "Public API"
	sectionPublicAPIBasePath = "/"
)

// NewAppSectionAdminAPI - создаёт объект mrperms.AppSection.
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

// NewAppSectionPublicAPI - создаёт объект mrperms.AppSection.
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
