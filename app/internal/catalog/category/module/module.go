package module

const (
	Name       = "Catalog.Category"   // Name - название модуля
	Permission = "modCatalogCategory" // Permission - разрешение модуля

	DBSchema              = "sample_catalog"         // DBSchema - схема БД используемая модулем
	DBTableNameCategories = DBSchema + ".categories" // DBTableNameCategories - таблица БД используемая модулем
	DBFieldTagVersion     = "tag_version"            // DBFieldTagVersion - поле для хранения версии записи
	DBFieldDeletedAt      = "deleted_at"             // DBFieldDeletedAt - поле содержит дату удаления записи, если NULL, то запись не удалена

	ImageDir = "catalog/categories" // ImageDir - относительный путь для хранения изображений категорий
)
