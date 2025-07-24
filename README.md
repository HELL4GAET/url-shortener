# ⚒️ URL Shortener ![Go](https://img.shields.io/badge/go-1.23-blue) ![Last Commit](https://img.shields.io/github/last-commit/HELL4GAET/url-shortener)

URL Shortener — это микросервис на Go, который позволяет сокращать длинные URL‑адреса в компактные "алиасы" и отслеживать статистику переходов. Проект упакован в Docker и полностью документирован через спецификацию OpenAPI, что упрощает интеграцию в любые CI/CD‑цепочки и оркестрацию через Docker Compose.


## 💻 Технологический стек
- Язык: Go
- База данных: PostgreSQL
- Контейнеризация: Docker + Docker Compose
- Документация API: OpenAPI 3.0 (docs/openapi.yaml) + Swagger UI
- Маршрутизация: chi



## 📦 Эндпоинты

    POST /api/v1/shorten
    Тело запроса: text/plain с URL
    Ответ: text/plain — короткий URL.

    GET /{alias}
    Перенаправление по коду (HTTP 307).

    GET /api/v1/stats/{alias}
    JSON со статистикой (clicks, created_at, original_url).

## 🚀 Быстрый старт
Склонируйте репозиторий и перейдите в него:

```bash
git clone https://github.com/HELL4GAET/url-shortener.git
cd url-shortener
```
Переименуйте или создайте файл .env с содержимым:
- POSTGRES_USER=postgres
- POSTGRES_PASSWORD=postgres
- POSTGRES_DB=postgres
- DATABASE_URL=postgresql://{POSTGRES_USER}:{POSTGRES_PASSWORD}@db:5432/{POSTGRES_DB}?sslmode=disable
- PORT=8080

Запустите все сервисы через Docker Compose:
```bash
docker-compose up --build
```
Подождите, пока поднимется база и приложение, затем откройте:
- API: http://localhost:8080/api/v1/shorten и передайте валидный url в теле запроса text/plain
- Swagger UI: http://localhost:8081


## 🔍 Swagger UI

После запуска доступен по адресу:
http://localhost:8081

## 🧪 Примеры curl
```bash
# Сократить URL
curl -X POST http://localhost:8080/api/v1/shorten -d "https://example.com"

# Перейти по короткому URL
curl -v http://localhost:8080/abc123

# Получить статистику
curl http://localhost:8080/api/v1/stats/abc123
```

## 📝 Лицензия ![License](https://img.shields.io/badge/license-MIT-green)

Этот проект распространяется под лицензией [MIT](LICENSE).



