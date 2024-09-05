# Прокси-сервер для загрузки превьюшек видеороликов с Youtube
Разработать gRPC прокси-сервис для загрузки превьюшек c видеороликов YouTube. При повторном запросе на тот же видеоролик, сервис
должен отдать закэшированный ответ (можно использовать рантайм кэш, но будет плюсом
если будет использоваться временное хранилище, например sqlite). Так же предлагается
написать клиентскую часть как утилиту командной строки, которой передается в качестве
параметров ссылки на видеоролики. В командной утилите предусмотреть ключ --async,
который позволяет скачивать большое количество файлов асинхронно.

## Содержание
- [Технологии](#технологии)
- [Начало работы](#использование)
- [Работа с утилитой](#описание)

## Технологии
- [Golang1.22.5](https://go.dev/doc/install)
- [Viper](https://github.com/spf13/viper)
- [Protobuf](https://protobuf.dev/)
- [Docker](https://www.docker.com/)
- [Memcached](https://memcached.org/)
- [Sqlite](https://www.sqlite.org/)
- [Minio](https://min.io/)

## Использование
Для начала работы с проектом необходимо выполнить ряд команд, описание ниже.
1) Установите Docker c официального источника
2) Установить компилятор, инструменты Golang версии, не ниже указанной в блоке [технологий](#технологии).
3) Настройте [git](https://git-scm.com/downloads) у себя на компьютере.

### Инструкция для запуска:

Выполняем в первом окне терминала:
```
$ git clone git@github.com:cantylv/thumbnail-loader.git
$ make init
$ make run 
```
Открываем другое окно терминала и пользуемся программой:

Пример
```
go run cmd/main/main.go --cache_inmemory=true --async=true --cache_timeout=100s https://www.youtube.com/watch\?v\=42M3esYyHdw https://www.youtube.com/watch\?v\=3EhxLz1EFIM https://www.youtube.com/watch\?v\=HjGDJk4kOYk https://www.youtube.com/watch\?v\=MYjgu-IzGEI https://www.youtube.com/watch\?v\=_MBdQhW22VI
```

## Описание
При запуске программы вы можете передать флаги, которые выполняют определенные функции:
```
-a, --async                    configure whether asynchronous loading is required
-c, --cache_inmemory           determines 'type' of cache; if true, cache data will be stored in ram, in another way in winchester
-t, --cache_timeout duration   the duration for which cache instance will store data (default 30s)
-u, --upload_folder string     the destination folder for uploading files from youtube (default "uploads")
```

Пример использования вы можете увидеть выше. Также в программе предусмотрена одна переменная окружения, которая отвечает за расположение папки загрузки изображений, ее название `UPLOAD_FOLDER`. 

**Примечание:** параметр --cache_inmemory по умолчанию имеет значение false, что означает кэширование в БД sqlite. Если значение будет true, то кэширование будет работать на основе in-memory хранилища Memcached.