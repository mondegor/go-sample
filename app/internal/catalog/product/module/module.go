package module

const (
	Name       = "Catalog.Product"   // Name - название модуля
	Permission = "modCatalogProduct" // Permission - разрешение модуля

	DBSchema            = "sample_catalog"       // DBSchema - схема БД используемая модулем
	DBTableNameProducts = DBSchema + ".products" // DBTableNameProducts - таблица БД используемая модулем
	DBFieldTagVersion   = "tag_version"          // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt    = "deleted_at"           // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена
)
