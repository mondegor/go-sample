# Описание Go Sample v0.13.7
Этот репозиторий содержит описание сервиса Go Sample.

## Статус сервиса
Сервис находится в рабочем состоянии, используется в качестве примера.

## Описание сервиса
Web сервис каталога продукции, который решает типовые задачи разработки подобных web сервисов.
В сервисе используется слоистая и гексагональная архитектура, выделено несколько архитектурных
слоёв, которые облегчают сопровождение сервиса и упрощают его расширение в дальнейшем.
В качестве хранилища данных выбрана postgres. Для хранения изображений используется S3 minio,
но можно настроить хранение файлов в обычной файловой системе
(для этого нужно заменить вызов с `NewS3Minio` на `NewFileStorage`).
Для распределённой блокировки ресурсов используется redis.

Подробнее смотри [документацию проекта](./docs/README.md)

> Перед запуском консольных скриптов сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

### Команды для сборки API документации v0.1.11
- `mrcmd openapi help` - помощь по командам плагина openapi;
- `mrcmd openapi build-all` - сборка документации всех API;

### Примеры запуска сборки документации из консоли Windows:
- GitBash (cmd): `"C:\Program Files\Git\git-bash.exe" --cd=d:\mrwork\go-sample mrcmd openapi build-all`;
- WSL (PowerShell): `cd D:\workdir\go-sample; wsl -d Ubuntu-20.04 -e mrcmd openapi build-all`;

## Разворачивание, установка и запуск сервиса

### Разворачивание сервиса
- Выбрать рабочую директорию, где должен быть расположен сервис
- `mkdir go-sample && cd go-sample` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/go-sample.git .`
- `cp .env.dist .env`
- `mrcmd state` // проверка состояния сервиса
- `mrcmd config` // проверка установленных переменных сервиса

### Установка сервиса и его первый запуск
- `mrcmd docker ps` // убеждаемся, что Docker daemon запущен
- `mrcmd install`
- `mrcmd start`
- `mrcmd docker-compose ps` // проверка всех запущенных ресурсов сервиса
- `mrcmd go-migrate up` // загрузка дампа с данными в БД
- `mrcmd go logs` // проверка, что основной сервис запущен

### Запуск и остановка сервиса
- `mrcmd start`
- `mrcmd stop`

### Остановка сервиса и удаление всех его установленных ресурсов
- `mrcmd uninstall`

### Часто используемые команды
- `mrcmd help` - помощь в контексте текущего сервиса;
- `mrcmd state` - общее состояние текущего сервиса;
- `mrcmd docker-compose conf` // отображает список `.yaml` файлов из которых собрана конфигурация и саму конфигурацию;
- `mrcmd docker-compose ps` - текущее состояние запущенных ресурсов;
- `mrcmd docker-compose logs` - логи всех запущенных ресурсов;
- `mrcmd go-migrate help` - все команды сервиса go-migrate;
- `mrcmd postgres help` - все команды сервиса postgres;
- `mrcmd go help` - выводит список всех доступных go команд (docker версия);
- `mrcmd go-dev help` // выводит список всех доступных go-dev команд (локальная версия);
- `mrcmd go-dev check` // статический анализ кода библиотеки (линтеры: govet, staticcheck, errcheck)
- `mrcmd go-dev test` // запуск тестов библиотеки
- `mrcmd golangci-lint check` // запуск линтеров для проверки кода (на основе `.golangci.yaml`)
- `mrcmd plantuml build-all` // генерирует файлы изображений из `.puml` [подробнее](https://github.com/mondegor/mrcmd-plugins/blob/master/plantuml/README.md#%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D0%B5%D0%B9-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0-markdown--plantuml)

> Более подробную информацию по использованию утилиты Mrcmd
> смотрите [здесь](https://github.com/mondegor/mrcmd#readme).

## Панели управления развёрнутой инфраструктурой
- MINIO: http://127.0.0.1:9984/ (admin 12345678)