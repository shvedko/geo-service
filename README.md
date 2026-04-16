# Geo-Service (Golang + PostGIS)

Микросервис на Go для работы с геоинформационными данными. Позволяет сохранять географические точки в базу данных и выполнять пространственный поиск по прямоугольной области (Bounding Box).

## Функциональность

*   **REST API**: для создания и поиска гео-объектов.
*   **PostGIS Интеграция**: использование пространственных индексов GIST для быстрого поиска.
*   **Строгая валидация**: проверка координат на вхождение в диапазоны WGS84 (Широта: [-90, 90], Долгота: [-180, 180]).
*   **Docker-native**: готов к развертыванию через Docker Compose.

---

## Технологический стек

*   **Язык**: Go 1.22
*   **База данных**: PostgreSQL 15 + PostGIS 3.3
*   **Библиотеки**: `pgx/v5` (драйвер БД)
*   **Контейнеризация**: Docker, Docker Compose

---

## Быстрый старт

### 1. Требования
Установленный Docker и Docker Compose.

### 2. Запуск
Выполните команду в корне проекта:

```bash
docker compose up --build
```

Сервис будет доступен по адресу: http://localhost:8080

---

## API Эндпоинты

### 1. Создание точки
**POST** /points

Добавляет новую географическую точку.

```bash
curl -X POST http://localhost:8080/points \
-H "Content-Type: application/json" \
-d '{"name": "Eiffel Tower", "lat": 48.8584, "lon": 2.2945}'
```

### 2. Поиск в прямоугольнике
**GET** /points/{min_lon}/{min_lat}/{max_lon}/{max_lat}

Возвращает объекты внутри заданных границ.

Параметры запроса:

* min_lon, max_lon: Долгота (диапазон от -180 до 180).
* min_lat, max_lat: Широта (диапазон от -90 до 90).

```bash
curl http://localhost:8080/points/2.0/3.0/48.0/49.0
```

---

## Локальная разработка
Для запуска вне Docker используйте флаги командной строки:

        --addr string   bind address
        --base string   data base (default "postgres://postgres:postgres@postgres:5432/geo")
    -h, --help          help for geo
        --port string   bind port (default "8080")

Схема для создания БД: 
    
[sql/schema/geo.sql](sql/schema/geo.sql)

Обслуживание БД:

- `CLUSTER points USING points_point_idx`

или

- `CREATE EXTENSION pg_repack`
- `pg_repack -d geo -t points --order-by point`

Спецификация OpenAPI:

[doc/swagger.yaml](doc/swagger.yaml)

или

[http://localhost:8080/swagger.yaml](http://localhost:8080/swagger.yaml)
