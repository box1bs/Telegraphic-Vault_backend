# Telegraphic Vault

Telegraphic Vault - это персональный менеджер знаний с открытым исходным кодом, разработанный для эффективного сохранения и организации информации. Проект сочетает в себе функции менеджера закладок и заметок, помогая систематизировать знания и идеи в процессе обучения и личностного роста.

## Основные возможности

- **Умное управление закладками:**
  - Сохранение веб-страниц с описанием и тегами
  - Поиск по содержимому и тегам
  - Организация закладок по категориям

- **Система заметок:**
  - Создание и редактирование заметок
  - Поддержка тегов для категоризации

- **Безопасность:**
  - Шифрование данных
  - Система аутентификации и авторизации
  - Защита от брутфорс-атак

## Для кого этот проект

- **Студентов:**
  - Организация учебных материалов
  - Сохранение полезных статей и исследований
  - Создание конспектов и заметок по предметам

- **Начинающих исследователей:**
  - Документирование собственных наблюдений и идей
  - Сбор и систематизация информации по интересующим темам
  - Ведение дневника исследований

- **Самоучек:**
  - Структурирование материалов для самообразования
  - Отслеживание прогресса в обучении
  - Создание персональной базы знаний

# API Documentation

## Auth Handlers

### `GET /auth`
Получение ключа шифрования для регистрации.

### `POST /auth`
Регистрация нового пользователя.

**Request Body:**
```json
{
  "username": "string",
  "password": "string",
  "key": "string"
}
```

### `POST /auth/login`
Авторизация пользователя.

**Request Body:**
```json
{
  "username": "string",
  "password": "string",
  "key": "string"
}
```

### `POST /auth/refresh`
Обновление токенов.

**Headers:**
- `Authorization: Bearer <refresh_token>`

## Bookmark Handlers

### `GET /app/bookmarks`
Получение всех закладок пользователя.

**Headers:**
- `Authorization: Bearer <token>`

### `POST /app/bookmarks`
Создание новой закладки.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "url": "string",
  "title": "string",
  "description": "string",
  "tags": ["string"]
}
```

### `PUT /app/bookmarks`
Обновление закладки.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "url": "string",
  "title": "string",
  "description": "string",
  "tags": ["string"]
}
```

### `DELETE /app/bookmarks`
Удаление закладки.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `uri`: URL закладки

### `GET /app/bookmarks/search`
Поиск закладок.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `q`: поисковый запрос

## Note Handlers

### `GET /app/notes`
Получение всех заметок пользователя.

**Headers:**
- `Authorization: Bearer <token>`

### `POST /app/notes`
Создание новой заметки.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "tags": ["string"]
}
```

### `PUT /app/notes`
Обновление заметки.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "current_title": "string",
  "new_title": "string",
  "content": "string",
  "tags": ["string"]
}
```

### `DELETE /app/notes`
Удаление заметки.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `title`: заголовок заметки

### `GET /app/notes/search`
Поиск заметок.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `q`: поисковый запрос
