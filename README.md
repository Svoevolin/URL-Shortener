# REST API "Укорачиватель ссылок"

## Подготовка

### Пример
1. Создать конфиг .yaml, установите значения
    ```zsh
    mkdir config && nano config/local.yaml
    ```
    ```yaml
    env: local                    # local, dev, prod
    http_server:                
      address: localhost:8080     
      timeout: 4s
      idle_timeout: 60s 
      alias_length: 10            # Длина алиасов для генерации по умолчанию
      user: "myuser"              # BasicAuth
      password: "mypass"          
    database:
      dsn: host=localhost port=5432 user=postgres password=postgres dbname=url sslmode=disable
    ```
2. Установите DSN и переменную окружения CONFIG_PATH в Makefile

    ```env
    DSN="host=localhost port=5432 user=postgres password=postgres dbname=url sslmode=disable"
    CONFIG_PATH=${CURDIR}/config/local.yaml
    ```
## Запуск

Запустить бота локально:
 ```zsh
 make goose-up && make run
```

## Описание

**Запросы**
- HTTP Basic Auth: **POST**, **PUT**, **DELETE**
- Без аутентификации: **GET**

#### POST:

- `/url/` - создаст укороченную ссылку

    - body:
      ```json
      {"url": "https://github.com/Svoevolin", "alias": "thebest"}
      ```
    - или alias будет сгенерирован сам:
      ```json
      {"url": "https://github.com/Svoevolin"}
      ```
#### PUT:

- `/url/` - изменит укороченную ссылку

    - body:
      ```json
      {"url": "https://github.com/Svoevolin", "alias": "numberOne"}
      ```

    - или alias будет сгенерирован сам:
      ```json
      {"url": "https://github.com/Svoevolin"}
      ```

#### DELETE:

- `/url/{alias}` - удалит укороченную ссылку

#### GET:

- `/{alias}` - перенаправит на полную ссылку по укороченной

## Архитектура
```
.
├── bin
├── cmd
├── config
│     └── local.yaml                           - конфигурационный файл
├── internal
│     ├── config
│     ├── database
│     │     ├── database.go                    - вынесены общие объекты (на случай добавления другой базы данных)
│     │     └── postgres                       - методы к базе данных
│     ├── http-server
│     │     ├── handlers
│     │     │     └── url
│     │     │           ├── delete             - хэндлер для DELETE запроса
│     │     │           ├── redirect           - хэндлер для GET    запроса
│     │     │           ├── save               - хэндлер для POST   запроса
│     │     │           └── update             - хэндлер для PUT    запроса
│     │     └── middleware
│     │           └── logger                   - миддлварь для логгирования запросов
│     └── lib
│           ├── api                            - содержит чекер редиректа и функцию для
│           │                                    получения URL параметров (для тестов chi роутера)
│           ├── logger
│           │     └── handlers
│           │           ├── slogpretty         - опции логгера (из пакета slog), форматирует json, делает красивыми логи 
│           │           └── slogdiscard        - mock-логгер, выбрасывает все логи (для тестов)
│           ├── random                         - функция генерации случайного алиаса
│           └── response                       - пакет для работы с ответами
|
├── tests                                      - функциональные тесты (черный ящик)
├── migrations                                 - миграции для базы данных  
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### Особенности
- Функционал покрыт unit и функциональными тестами