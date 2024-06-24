# Go Sample Changelog
Все изменения сервиса Go Sample будут документироваться на этой странице.

## 2024-06-25
### Added
- Подключён планировщик задач `mrworker`.`mrschedule`;
- Подключён компонент `mrsettings` для доступа к произвольным настройкам;
- Подключены `prometheus` метрики для сбора статистики http запросов и работы БД;
- Подключён `mrsentry.Adapter` для отправки ошибок в `sentry`;
- Подключены линтеры с их настройками (`.golangci.yaml`);
- Добавлены комментарии для публичных объектов и методов;
- Добавлена конфигурационная переменная `Environment` для задания рабочего окружения;

### Changed
- Обновлена система формирования ошибок на основе новой версии библиотеки `go-sysmess`:
  - изменён формат создания новых ошибок;
  - объект `AppErrorFactory` заменён на `ProtoAppError` который теперь сам является ошибкой;
- `MimeTypeList` теперь задаётся из `config.yaml`;

### Removed
- Удалена поддержка соединения http сервера по сокету,
  также удалены `ListenTypeSock`, `ListenTypePort`;

## 2024-04-10
### Added
- Добавлена поддержка роутера `chi`;
- Добавлен `MiddlewareRecoverHandler` и обработчик `mrresp.HandlerGetFatalErrorAsJson`;
- Добавлен обработчик `mrresp.HandlerGetStatInfoAsJson` (`/v1/stat-info`) для фиксации и
  отображения статистики по запросам;

### Changed
- Изменение формата параметров в URL для роутера: `/v1/sample/:id -> /v1/sample/{id}`,
  такого вида параметры совместимы с `chi` роутером;
- Переименован пакет: `mrresponse` -> `mrresp`;
- Заменена строка `*path` в URL получения файлов на константу `mrserver.VarRestOfURL`;
- Переработан механизм `soft delete`, добавлено поле `deleted_at`;

### Removed
- Удалён статус `ItemStatusRemoved`, теперь удаление контролируется через
  отдельное поле `deleted_at`;

## 2024-03-23
### Added
- Добавлены следующие типы ошибок:
    - `FactoryErrCategoryNotAvailable`;
    - `FactoryErrTrademarkNotAvailable`;

### Changed
- В местах использования метода `mrfactory.WithPermission` добавлен `mrfactory.PrepareEachController`;
- `mrserver.NewMiddlewareHttpHandlerAdapter -> mrserver.MiddlewareHandlerAdapter`;
- Доработаны функции типа `factory.registerAdminAPIControllers`, заменены на `createAdminAPIControllers`
  с использованием новой функции `factory.registerControllers`;

### Removed
- Удален метод `IsExist` вместо него теперь используется `FetchStatus`;

## 2024-03-20
### Changed
- Обновлена структура БД (в том числе поля created_at и updated_at размещены внизу таблицы);
- Поле `Product.Price` теперь является указателем для возможности записи в
  БД нулевых значений;
- В `factory.NewRestServer` создание модулей вынесено в методы подобные этому:
  `registerAdminAPIControllers`;

## 2024-03-19
### Added
- Добавлены следующие ошибки:
    - `catalog.FactoryErrCategoryRequired`;
    - `catalog.FactoryErrTrademarkRequired`;
- Добавлен API компонент `App.Field.RewriteName`;

### Changed
- Переименованы методы:
    - `NewFetchParams -> NewSelectParams`;
    - `GetMetaData -> NewOrderMeta`;

## 2024-03-18
### Changed
- Внедрена новая версия библиотеки `go-sysmess`, в связи с этим:
    - в функции `IsAutoCallerOnFunc` изменено условие с использованием `HasCallStack()`;
- В некоторых API методах тип `PUT` преобразован в `PATCH` для более строгого соответствия API спецификации;

### Fixed
- В `Catalog.Categories.AdminAPI.Request.Path.CategoryID` исправлен тип с `intеger` на `UUID`;

## 2024-03-17
### Added
- Добавлен `App.Response.Model.SuccessCreatedItemInt32` в API и в `pkg`;

### Changed
- Идентификатор категории был заменён с int на uuid;
- Перенесены в `pkg` часто используемые сервисом модели:
    - `SuccessCreatedItemResponse`;
    - `ChangeItemStatusRequest`;
    - `MoveItemRequest`;

### Removed
- Удалён `App.Response.Model.Success`;
- Удалён `App.Response.Model.SuccessModifyItem`;

## 2024-03-16
### Changed
- Все поля БД типа `timestamp` теперь с `with time zone`;
- Заменено `version -> tagVersion`;

## 2024-03-15
### Changed
- Рефакторинг кода:
    - переименование `FactoryErrService* -> FactoryErrUseCase*`, `errService* -> errUseCase*`;
    - переименование интерфейсов `*Service -> *UseCase`;
    - замена методов `LoadOne` на `FetchOne`;
    - методы `Create`, `Insert` теперь возвращают ID записи;
    - схема БД переименована `gscatalog -> sample_catalog`;
- Вся мета информация об изображениях стала необязательной (`imageUrl`, и т.д.);
- Настройки `PageSizeMax` и `PageSizeDefault` вынесены в общие настройки модулей `ModulesSettings.General`;
- Парсер `SortPage` разделён на два: `ListSorter`, `ListPager`;
- Удалены неиспользуемые параметры запросов в каждом из модулей, отсортированы по алфавиту оставшиеся;
- В логгер добавлена поддержка `IsAutoCallerOnFunc`;


- Рефакторинг API документации:
    - Добавлены компоненты:
        - `App.Response.Model.BinaryAnyFile`;
        - `App.Response.Model.BinaryImage`;
        - `App.Response.Model.BinaryMedia`;
        - `App.Response.Model.JsonFile`;
        - `App.Response.Model.SuccessModifyItem`;
        - `App.Response.Model.TextFile.yaml`;
    - Доработка описания фильтрации, сортировки при получении списков записей;
    - Доработка описания ограничений при добавлении/обновлении записей;
    - Для всех модулей поля-идентификаторы описаны как отдельные сущности;

## 2024-02-05
### Changed
- Переименованы:
    - `datetime_created -> created_at`;
    - `datetime_updated -> updated_at`;
    - `modules.Options -> app.Options`;
- Создание модулей переехало в `factory/modules/*`;
- Большинство юнитов было преобразовано в модули, которые объединены доменами;

## 2024-01-30
### Changed
- Внедрён новый интерфейс логгера, добавлен режим трассировки запросов;
- Для многих методов добавлен параметр `ctx context.Context`;
- Заменён устаревший интерфейс `mrcore.EventBox` на `mrsender.EventEmitter`;
- Переименован `ServiceHelper -> UsecaseHelper`;
- Внедрены `mrlib.CallEachFunc`, `CloseFunc` для группового закрытия ресурсов;
- Переименован `CorrelationID` на `X-Correlation-ID`;
- Объекты конфигураций/опций теперь передаются по значению (`*Config -> Config`, `*Options -> Options`);
- Внедрён `oklog/run` для управления одновременным запуском нескольких серверов (http, grpc)
- Добавлены методы для создания и инициализации всех глобальных настроек приложения
  (`CreateAppEnvironment`, `InitAppEnvironment`);
- Теперь модули собираются в рамках отдельных серверов (см. `factory.NewRestServer`);
- Изменены некоторые переменные окружения:
    - удалён `APPX_LOG_PREFIX`;
    - добавлен `APPX_LOG_TIMESTAMP=RFC3339|RFC3339Nano|DateTime|TimeOnly` (формат даты в логах);
    - добавлен `APPX_LOG_JSON=true|false` (вывод логов в json формате);
    - добавлен `APPX_LOG_COLOR=true|false` (использование цветного вывода логов в консоле);
    - переименованы:
        - `APPX_SERVICE_LISTEN_TYPE -> APPX_SERVER_LISTEN_TYPE`;
        - `APPX_SERVICE_LISTEN_SOCK -> APPX_SERVER_LISTEN_SOCK`;
        - `APPX_SERVICE_BIND -> APPX_SERVER_LISTEN_BIND`;
        - `APPX_SERVICE_PORT -> APPX_SERVER_LISTEN_PORT`;

## 2024-01-25
### Added
- Внедрены парсеры на основе интерфейсов `mrserver.RequestParserFile` и
  `mrserver.RequestParserImage` для получения файлов и изображений из `multipart` формы.
  - заменено `mrreq.File -> ht.parser.FormImage`;
  - в `CategoryImageService` изменён тип `mrtype.File -> mrtype.Image`;

### Changed
- Переименовано `ConvertImageMetaToInfo -> ImageMetaToInfoPointer`;

### Removed
- `mrserver.RequestParserPath` удалён вместо него используется
  `mrserver.RequestParserString` и `mrserver.RequestParserParamFunc`;

## 2024-01-22
### Changed
- Расформирован объект `ClientContext` и его одноименный интерфейс, в результате:
    - Изменена сигнатура обработчиков с `func(c mrcore.ClientContext)` на `func(w http.ResponseWriter, r *http.Request) error`;
    - С помощью интерфейсов `RequestDecoder`, `ResponseEncoder` можно задавать различные форматы
      принимаемых и отправляемых данных (сейчас реализован только формат `JSON`);
    - Запросы обрабатываются встраиваемыми в обработчики объектов `mrparser.*` через интерфейсы:
      `mrserver.RequestParserPath`, `RequestParser`, `RequestParserItemStatus`, `RequestParserKeyInt32`,
      `RequestParserSortPage`, `RequestParserUUID`, `RequestParserValidate`;
    - Ответы отправляются встраиваемыми в обработчики объекты `mrresponse.*` через интерфейсы:
      `mrserver.ResponseSender`, `FileResponseSender`, `ErrorResponseSender`;
    - Вместо метода `Validate(structRequest any)` используется объект `mrparser.Validator`;
- Произведены следующие замены:
    - `HttpController.AddHandlers -> Handlers() []HttpHandler`
      убрана зависимость контроллера от роутера и секции,
      для установки стандартных разрешений добавлены следующие методы
      `mrfactory.WithPermission`, `mrfactory.WithMiddlewareCheckAccess`;
    - `ModulesAccess -> AccessControl` (`modules_access -> access_control`) и добавлен интерфейс `mrcore.AccessControl`;
    - `ClientSection -> AppSection` (`client_section -> app_section`) удалена зависимость от `AccessControl`;
- При внедрении новой версии библиотеки `go-sysmess` было заменено:
    - `mrerr.FieldErrorList -> CustomErrorList`;

## 2024-01-19
### Changed
- Enum тип БД gs_catalog.item_status заменён на int2 и удалён. Доработано, чтобы `enum` типы сохранялись в виде `int`;
- Код получения файла в обработчике заменён на `mrreq.File`;
- Переименован метод `checkProduct` в `usecase` на более абстрактный `checkItem`
  (проверяет возможность добавления, сохранения записи);
- Добавлен метод `prepareItem` в `usecase` для подготовки записи перед её отправкой в ответе;
- Все языковые идентификаторы приведены к типу uint16;

## 2024-01-16
### Added
- Добавлены новые `OpenAPI` компоненты: `App.Measure*`, `App.Response.Model.FileInfo`,
  `App.Response.Model.ImageInfo`;
- Для каждой секции добавлены настройки `AuthSecret` и `AuthAudience`;
- Добавлены системные обработчики (`RegisterSystemHandlers`);
- Добавлена фильтрация поля цены товара (`Custom.Request.Query.Filter.Price*`);
- Добавлено поле `Config.AppStartedAt` для отслеживания времени запуска сервиса;

### Changed
- Поле `categories.image_path` заменено на `categories.image_meta` типа `jsonb`,
  в котором теперь хранится мета информация о файле;
- В каждом модуле теперь собственные Options, которые отделены от общего конфига;

## 2023-12-10
### Changed
- Внедрено использование `CallStack` в `mrerr.AppError` и `mrcore.Logger`, а также функция `CallerEnabledFunc` для отключения избыточной информации;
- Внедрён переработанный механизм работы с виртуальными хранилищами данных;
- Изменения в REST API документации:
    - добавлена общая информация о пользовательских ограничениях полей и ошибках
      используемая во всех документах в виде отдельного файла;
    - в самих методах таких как добавление, сохранение приводится только краткая запись пользовательских ограничений полей и ошибок;
    - добавлено поле `UpdatedAt`, исправлены мелкие ошибки;
    - добавлен пакет measures со стандартными компонентами;
- Добавлена обработка и вывод поля `UpdatedAt` с учётом того, что оно необязательное;
- Схема БД `gosample` заменена на `gs_catalog`, у таблиц удалён префикс `catalog_`;
- Заменено `mrerr.Arg -> mrmsg.Data`, `S3Pool -> FileProviderPool`;

## 2023-12-07
### Changed
- Доработаны настройки конфига, добавлены `APPX_SERVER_*` и другие переменные;
- Добавлено управление callstack в отладочном режиме (cм. `mrerr.SetCallStackOptions`);
- Префикс в записях логов теперь стал необязательным параметром;
- Переработана работа с ошибками под новое ядро системы, в связи с этим:
    - из `mrcore.FactoryErrServiceEntityNotFound` удалены параметры;
    - для обработки `tagVersion` теперь используется `mrcore.FactoryErrServiceEntityVersionInvalid`;
    - `mrcore.FactoryErrServiceTemporarilyUnavailable` заменено на `uc.serviceHelper.WrapErrorFailed`;
    - `uc.serviceHelper.WrapErrorForSelect -> uc.serviceHelper.WrapErrorEntityNotFoundOrFailed`;
    - `mrcore.FactoryErrServiceIncorrectSwitchStatus -> mrcore.FactoryErrServiceSwitchStatusRejected`;
    - добавлены ошибки: `FactoryErrCategoryImageNotFound`, `FactoryErrProductNotFound`;
- Добавлен механизм виртуальных файловых провайдеров (S3 хранилищ).
  В конфиге прописан виртуальный файловый провайдер `imageStorage` с привязкой к реальному бакету S3 хранилища.
  Его уже используют другие модули системы. Тем самым убрана прямая зависимость модулей от реального бакета.
- В контроллере из `category.go` функционал работы с изображениями перенесён в `category_images.go`;
- Доработана структура БД, добавлены тестовые данные. Схема `public` переименована в `gosample`.
  Для каждого юнита добавлена своя константа с названием схемы БД;
- Изменения в REST API документации:
    - добавлен `App.Response.Model.ErrorList` (массив `App.Response.Model.ErrorAttribute`);
    - переименован `App.Response.Model.Error -> App.Response.Model.ErrorDetails`;
    - переименовано поле у `App.Response.Model.ErrorDetails`: `detail -> details`;
- Теперь юниты используют обращение друг к другу через API.
  Для этого вместо storageCategory и storageTrademark ранее используемые юнитом product добавлены соответствующие CategoryServiceAPI и TrademarkServiceAPI;

## 2023-11-23
### Changed
- Отладочная информация о загрузке картинки теперь вызывается из пакета mrdebug;

### Fixed
- Добавлен пропущенный cursor.Err();
- поправлено неправильная замена module.DBSchema в GetMetaData; 

## 2023-11-23
### Changed
- Обновлены зависимости библиотеки;

## 2023-11-20
### Changed
- Обновлена REST API документация, подключена её сборка в OpenAPI формате;
- Доработана работа с тегом версии сущности, которая позволяет избежать сохранения устаревших данных;
- Большинство объектов mrcore.FactoryErrServiceIncorrectInputData заменено на mrcore.FactoryErrServiceEntityNotFound;
- Название схемы данных БД вынесена в константу module.DBSchema;
- Доработан конфиг приложения, добавлен ModulesSettings для настроек модулей;
- Переименовано:
    - `mrcore.ClientData -> mrcore.ClientContext`;
    - `ClientContext::ParseAndValidate -> ClientContext::Validate` (удалён `ClientContext::Parse`);
- Обновлены зависимости от библиотек, доработаны компоненты под новое ядро системы;
- Обновлён `.editorconfig`;

## 2023-11-13
### Changed
- В системе выделены два модуля, которые не зависимы друг от друга. В связи с этим удалён префикс `Catalog` у всех сущностей, т.к. теперь и так понятно в каком модуле они находятся.
- Переименованы некоторые переменные и функции (типа Id -> ID) в соответствии с code style языка go;
- Доработаны компоненты работающие с файлами под новый интерфейс `FileProviderAPI`;
- Тексты SQL запросов приведены к единому стилю;
- Обновлены зависимости от библиотек, доработаны компоненты под новое ядро системы;
- Все файлы библиотеки были пропущены через `gofmt`;

## 2023-11-01
### Added
- Внедрена работа с пользовательскими разрешениями и привилегиями (ролевая модель) для каждого модуля, добавлены настройки для ролей в `config`;
- Реализована поддержка фильтров и компонентов `Sorter` и `Pager` для списков API (обработка запроса, построение SQL фрагментов);
- Логика каждого модуля разделена на независимые API (для примера `public` и `admin-api`);
- Добавлен контроллер для загрузки картинок в `public API`;
- В `.gitignore` добавлена директория `/golang` которую создаёт докер;

### Changed
- Переименованы следующие сущности:
    - `mrcom_status -> mrenum`;
    - `mrcom_orderer -> mrorderer`;
- Обновлены зависимости библиотеки;

## 2023-10-09
### Changed
- Обновлены зависимости библиотеки;
- Обработка ошибок приведена к более компактному виду;
- Переработаны фильтры списков;
- Доработана работа со статусами и их переключениями;

### Fixed
- Добавлены настройки `APPX_REDIS_*` и `APPX_S3_*` для правильного запуска сервиса в докере;

## 2023-09-20
### Changed
- Обновлены зависимости библиотеки;
- Фиксация зависимостей инфраструктуры;
- Заменён адаптер `*mrpostgres.ConnAdapter` на интерфейс `mrstorage.DbConn`;
- Заменены tabs на пробелы в коде;

## 2023-09-16
### Add
- К категориям продукции добавлена возможность прикрепления изображений;
- Добавлено поле UpdateAt для всех объектов;

### Changed
- Все объекты, которые создаются при запуске сервиса перенесены в пакет `factory`;
