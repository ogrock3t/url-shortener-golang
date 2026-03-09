# URL Shortener Golang

Простой API-сервис для сокращения ссылок, написанный на Golang.  
Сервис принимает длинный URL, создает для него короткую ссылку длиной 10 символов и поддерживает перенаправление на оригинальный адрес.

Проект реализован с поддержкой двух хранилищ:

- in-memory storage
- PostgreSQL

Тип хранилища выбирается параметром при запуске сервиса.

---

## Функциональность

Сервис поддерживает следующие возможности:

- Создание короткой ссылки из оригинального URL
- Редирект по короткой ссылке на оригинальный URL
- Получение оригинального URL по короткой ссылке через отдельный endpoint
- Валидация входных URL
- Один оригинальный URL соответствует одной короткой ссылке
- Поддержка двух реализаций хранилища:
    - in-memory
    - PostgreSQL
- Выбор хранилища через параметр запуска
- HTTP middleware для логирования запросов
- Unit-тесты для `domain`, `service` и `handlers`
- Docker-образ для запуска сервиса

---

## Технологии

- Go
- net/http
- PostgreSQL
- pgx / pgxpool
- Docker
- sync.RWMutex
- context
- httptest

---

## API

### Создать короткую ссылку

**POST** `/shorten`

#### Request

```json
{
  "url": "https://google.com"
}
```

#### Response

```json
{
  "short_url": "http://localhost:8080/AAAAAAAAAA"
}
```

---

### Редирект по короткой ссылке

**GET** `/{short_code}`

#### Пример

```http
GET /pa_5Af7cAX
```

#### Ответ

```http
302 Found
Location: https://github.com/ogrock3t
```

> Для endpoint `GET /{short_code}` выбрано поведение через HTTP redirect, так как это соответствует типичному сценарию использования URL shortener-сервисов.

---


### Получить оригинальный URL по короткой ссылке

**GET** `/resolve/{short_code}`

#### Пример

```http
GET /resolve/pa_5Af7cAX
```

#### Response

```json
{
  "original_url": "https://github.com/ogrock3t"
}
```

---

## Установка

### Клонировать репозиторий

```bash
git clone https://github.com/ogrock3t/url-shortener-golang.git
cd url-shortener-golang
```

---

## Конфигурация

Пример `.env` файла:

```env
PORT=8080
BASE_URL=http://localhost:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable
```

---

## Запуск сервера

### In-memory storage

```bash
go run ./cmd/url-shortener/main.go -storage=in-memory
```

### PostgreSQL storage

```bash
go run ./cmd/url-shortener/main.go -storage=postgres
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
```

---

## Запуск PostgreSQL через Docker

```bash
docker compose up -d
```

После этого можно запускать сервис с PostgreSQL storage.

---

## Запуск Docker-образа

### Сборка образа

```bash
docker build -t url-shortener .
```

### Запуск с in-memory storage

```bash
docker run -p 8080:8080 url-shortener -storage=in-memory
```

### Запуск с PostgreSQL storage

```bash
docker run --env-file .env -p 8080:8080 url-shortener -storage=postgres
```

---

## Пример использования

### Создать короткую ссылку

```bash
curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{"url":"https://google.com"}'
```

#### Ответ

```json
{
  "short_url": "http://localhost:8080/AAAAAAAAAA"
}
```

---

### Перейти по короткой ссылке

```bash
curl -v http://localhost:8080/AAAAAAAAAA
```

#### Ответ

```http
HTTP/1.1 302 Found
Location: https://google.com
```

---

### Получить оригинальный URL через API

```bash
curl http://localhost:8080/resolve/AAAAAAAAAA
```

#### Ответ

```json
{
  "original_url": "https://google.com"
}
```

---

## Алгоритм генерации короткой ссылки

Короткая ссылка не хранится в базе данных как отдельное поле.  
Сервис генерирует её детерминированно из числового `id`.

Используется следующий подход:

- `id` кодируется в строку длиной 10 символов
- используется алфавит:
    - `A-Z`
    - `a-z`
    - `0-9`
    - `_`
- декодирование работает в обратную сторону:
    - `short_url -> id`
    - `id -> original_url`

Преимущества такого подхода:

- ссылка всегда соответствует требованиям задания
- не нужно хранить `short_url` отдельно
- алгоритм простой и воспроизводимый
- короткая ссылка однозначно восстанавливает `id`

---