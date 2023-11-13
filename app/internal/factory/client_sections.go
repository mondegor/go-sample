package factory

import (
    "go-sample/config"
    "go-sample/internal/global"

    "github.com/mondegor/go-webcore/mrperms"
)

func NewClientSectionAdminAPI(cfg *config.Config, access *mrperms.ModulesAccess) *mrperms.ClientSection {
    return mrperms.NewClientSection(
        cfg.ClientSections.AdminAPI.Caption,
        global.SectionAdminAPIRootPath,
        cfg.ClientSections.AdminAPI.Privilege,
        access,
    )
}

func NewClientSectionPublicAPI(cfg *config.Config, access *mrperms.ModulesAccess) *mrperms.ClientSection {
    return mrperms.NewClientSection(
        cfg.ClientSections.PublicAPI.Caption,
        global.SectionPublicAPIRootPath,
        cfg.ClientSections.PublicAPI.Privilege,
        access,
    )
}
