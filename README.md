# URL Shortener Golang

Простой API-сервис для сокращения ссылок, написанный на Go.  
Сервис принимает длинный URL и возвращает короткую ссылку, которая перенаправляет на оригинальный адрес.

---

## Функциональность

Сервис поддерживает следующие возможности:

- Создание короткой ссылки из оригинального URL
- Редирект по короткой ссылке на оригинальный URL
- Валидация входных URL
- Один оригинальный URL соответствует одной короткой ссылке
- In-memory хранилище
- HTTP middleware для логирования запросов
- Unit-тесты 

---

## Технологии

- Go
- net/http
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
  "short_url": "http://localhost:8080/AAAAAAAAAB"
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

---

## Установка

### Клонировать репозиторий

```bash
git clone https://github.com/your-username/url-shortener-golang.git
cd url-shortener-golang
```

---

## Запуск сервера

```bash
go run ./cmd/url-shortener/main.go
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
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
Location: https://google.com
```

---

## Тестирование

Запуск всех тестов:

```bash
go test ./...
```

---

## Будущие улучшения

- Добавление PostgreSQL
- Docker контейнеризация