package factory

import (
    "go-sample/config"

    "github.com/mondegor/go-webcore/mrperms"
)

func NewClientSectionPublic(cfg *config.Config, access *mrperms.ModulesAccess) *mrperms.ClientSection {
    return mrperms.NewClientSection(
        cfg.ClientSections.Public.Name,
        cfg.ClientSections.Public.Caption,
        cfg.ClientSections.Public.Privilege,
        access,
    )
}

func NewClientSectionAdminPanel(cfg *config.Config, access *mrperms.ModulesAccess) *mrperms.ClientSection {
    return mrperms.NewClientSection(
        cfg.ClientSections.AdminPanel.Name,
        cfg.ClientSections.AdminPanel.Caption,
        cfg.ClientSections.AdminPanel.Privilege,
        access,
    )
}
