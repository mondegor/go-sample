# Описание Go Sample v0.1.1
Этот репозиторий содержит описание сервиса Go Sample.

## Статус сервиса
Сервис находится в стадии разработки.

## Описание сервиса
Web сервис каталога продукции, который решает типовые задачи разработки подобных web сервисов.
В сервисе выделено несколько архитектурных слоёв, которые облегчают сопровождение сервиса и упрощают его расширение в дальнейшем. В качестве хранилища данных выбрана postgres. Для хранения изображений используется S3 minio, но можно настроить хранение файлов в обычной файловой системе (для этого нужно заменить вызов с NewS3Minio на NewFileStorage). Для распределённой блокировки ресурсов используется redis.

## REST API документация
- https://github.com/mondegor/go-sample/blob/master/docs/catalog.yaml

## Разворачивание, установка и запуск сервиса

### Разворачивание сервиса
> Перед разворачиванием сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- Выбрать рабочую директорию, где должен быть расположен сервис
- `mkdir go-sample && cd go-sample` // создать и перейти в директорию проекта
- `git clone -b latest git@github.com:mondegor/go-sample.git .`
- `cp .env.dist .env`
- `mrcmd state` // проверка состояния сервиса
- `mrcmd config` // проверка установленных переменных сервиса

> Более подробную информацию по использованию утилиты Mrcmd смотрите [здесь](https://github.com/mondegor/mrcmd#readme).

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
- `mrcmd docker-compose ps` - текущее состояние запущенных ресурсов;
- `mrcmd docker-compose logs` - логи всех запущенных ресурсов;
- `mrcmd go help` - все команды сервиса go;
- `mrcmd go-migrate help` - все команды сервиса go-migrate;
- `mrcmd postgres help` - все команды сервиса postgres;