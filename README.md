# Wishlist API

REST API сервис для создания и управления вишлистами (списками желаемых подарков).
Пользователь может зарегистрироваться, создать вишлист к событию, добавить в него подарки и поделиться публичной ссылкой, по которой другие пользователи могут просмотреть список и забронировать подарок.

## Функциональность

### Авторизация

* Регистрация пользователя по **email + пароль**
* Вход в систему
* Пароли хранятся в хэшированном виде
* Закрытые эндпоинты доступны только авторизованным пользователям

### Вишлисты

Авторизованный пользователь может:

* создавать вишлисты
* просматривать свои вишлисты
* редактировать вишлисты
* удалять вишлисты

Каждый вишлист содержит:

* название события
* описание
* дату события

### Позиции вишлиста

Каждый вишлист содержит список подарков.

У позиции есть:

* название
* описание
* ссылка на товар
* приоритет (степень желаемости подарка)
* статус бронирования

Пользователь может:

* добавлять позиции
* редактировать позиции
* удалять позиции

### Публичный доступ

При создании вишлиста генерируется **уникальный токен**.

По этому токену можно:

* получить публичный вишлист без авторизации
* забронировать подарок

Если подарок уже забронирован — API возвращает ошибку.

---

# Технологии

* **Go 1.25** (gomock для тестирования)
* **PostgreSQL** (pgx, squirrel)
* **Docker / Docker Compose**
* **SQL migrations** (golang-migrate)
* **REST API** (chi, golang-jwt)
* **JSON** (validator)

---

# Запуск проекта

Для запуска требуется только **Docker**.

```bash
make up
```

После запуска сервис будет доступен по адресу:

```
http://localhost:8080
```

Просмотр логов.

```
make logs
```

Остановка (graceful shutdown).

```bash
make down
```

Запуск unit-тестов с отчетом о покрытии.

```bash
make test
```


---

# Конфигурация

Все параметры конфигурации задаются через переменные окружения.

Пример конфигурации находится в файле:

```
.env.example
```

Перед запуском можно создать `.env` на основе примера.

Основные параметры:

* `POSTGRES_USER`
* `POSTGRES_PASSWORD`
* `POSTGRES_DB`

* `DB_HOST`
* `DB_PORT`
* `DB_NAME`
* `DB_USER`
* `DB_PASSWORD`

* `HTTP_ADDR`
* `JWT_SECRET`

---

# Структура проекта

```
.
├── cmd/                    # Точка входа приложения
├── internal/
│   ├── application/        # Бизнес логика
│   ├── domain/             # Бизнес сущности
│   ├── config/             # Конфигурация
│   └── infrastructure/     # Внешние компоненты (Хендлеры, репозитории, мидлвары, ...)
├── migrations/             # SQL миграции
├── docker-compose.yml
├── Dockerfile
└── README.md
```

---

# API

## Авторизация

### Регистрация

```
POST /api/auth/register
```

```
curl -X POST http://localhost:8080/api/auth/register \ 
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@mail.com",
    "password": "123456"
  }'
```

### Вход

```
POST /api/auth/login
```

```
curl -X POST http://localhost:8080/api/auth/login \ 
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@mail.com",
    "password": "123456"
  }'
```

---

## Вишлисты

### Создать вишлист

```
POST    /api/wishlists
```

```
curl -X POST http://localhost:8080/api/wishlists \                                              
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Birthday",     
    "description": "My birthday wishlist",
    "event_date": "2026-12-01" 
  }'
```

### Посмотреть все вишлисти

```
GET     /api/wishlists
```

```
curl -X GET http://localhost:8080/api/wishlists \                                 
  -H "Authorization: Bearer $TOKEN"
```

### Посмотреть конректный вишлист по ID

```
GET     /api/wishlists/{id}
```

```
curl -X GET http://localhost:8080/api/wishlists/1 \                             
  -H "Authorization: Bearer $TOKEN"
  ```

### Обновить вишлист

```
PUT     /api/wishlists/{id}
```

```
curl -X PUT http://localhost:8080/api/wishlists/1 \
  -H "Authorization: Bearer $TOKEN" \ 
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Birthday",
    "description": "Updated desc",
    "event_date": "2026-12-10" 
  }'
```

### Удалить вишлист

```
DELETE  /api/wishlists/{id}
```

```
curl -X DELETE http://localhost:8080/api/wishlists/1 \            
  -H "Authorization: Bearer $TOKEN"
```

---

## Позиции вишлиста

### Создать позицию

```
POST    /api/items/{wishlist_id}
```

```
curl -X POST http://localhost:8080/api/items/1 \                                     
  -H "Authorization: Bearer $TOKEN" \  
  -H "Content-Type: application/json" \
  -d '{
    "title": "iPhone",
    "description": "Pro Max",
    "url": "https://apple.com",
    "priority": 10
  }'
```

### Посмотреть все позиции вишлиста

```
GET     /api/items/{wishlist_id}
```

```
curl -X GET http://localhost:8080/api/items/1 \                        
  -H "Authorization: Bearer $TOKEN"
```

### Обновить позицию

```
PUT     /api/items/{id}
```

```
curl -X PUT http://localhost:8080/api/items/1 \               
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "iPhone 15",       
    "description": "Updated",     
    "url": "https://apple.com",
    "priority": 5
  }'
```

### Удалить позицию

```
DELETE  /api/items/{id}
```

```
curl -X DELETE http://localhost:8080/api/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## Публичные эндпоинты

### Получить вишлист по токену

```
GET api/public/{token}
```

```
curl -X GET http://localhost:8080/api/public/$PUBLIC_WISHLIST_TOKEN
```

### Забронировать подарок

```
POST api/public/{token}/items/{id}/reserve
```

```
curl -X POST http://localhost:8080/api/public/$PUBLIC_WISHLIST_TOKEN/items/1/reserve
```

---

# Дополнительные улучшения

В проекте реализованы:

* Unit-тесты для бизнес-логики
* Graceful shutdown