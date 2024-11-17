package module

const (
	Name       = "Catalog.Trademark"   // Name - название модуля
	Permission = "modCatalogTrademark" // Permission - разрешение модуля

	DBSchema              = "sample_catalog"         // DBSchema - схема БД используемая модулем
	DBTableNameTrademarks = DBSchema + ".trademarks" // DBTableNameTrademarks - таблица БД используемая модулем
	DBFieldTagVersion     = "tag_version"            // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt      = "deleted_at"             // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
