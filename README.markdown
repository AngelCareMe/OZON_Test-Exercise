# Ozon Test

## Описание
Это приложение представляет собой простую социальную сеть, реализованную на Go, которая позволяет создавать посты, комментировать их, а также отключать комментарии для определённых постов. Поддерживаются два типа хранилищ: in-memory и PostgreSQL.

## Структура проекта
- **cmd/**: Точка входа приложения (`main.go`).
- **config/**: Управление конфигурацией приложения через YAML.
- **internal/api/**: Обработчики HTTP-запросов для постов и комментариев.
- **internal/models/**: Определения структур данных (`Post`, `Comment`).
- **internal/services/**: Бизнес-логика для работы с постами и комментариями.
- **internal/storage/**: Реализация хранилищ (in-memory и PostgreSQL).
- **migrations/**: SQL-миграции для PostgreSQL.

## Зависимости
- Go 1.24
- PostgreSQL 13 (для использования PostgreSQL-хранилища)
- Docker и Docker Compose (для контейнеризации)
- Библиотеки Go:
  - `github.com/spf13/viper` – для работы с конфигурацией
  - `github.com/jackc/pgx/v4` – для работы с PostgreSQL
  - `github.com/Masterminds/squirrel` – для построения SQL-запросов
  - `github.com/golang-migrate/migrate/v4` – для миграций базы данных

## Установка и запуск

### Локальный запуск (in-memory)
1. Убедитесь, что Go установлен (`go version`).
2. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/AngelCareMe/OZON_Test-Exercise
   cd ozon_test
   ```
3. Установите зависимости:
   ```bash
   go mod download
   ```
4. Скомпилируйте и запустите приложение:
   ```bash
   go build -o main ./cmd/main.go
   ./main -storage=inmemory
   ```
5. Сервер будет доступен на `http://localhost:8080`.

### Локальный запуск (PostgreSQL)
1. Убедитесь, что Docker и Docker Compose установлены.
2. Создайте файл `.env` на основе `.env.example` и задайте переменные окружения:
   ```bash
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=yourdb
   ```
3. Запустите приложение с помощью Docker Compose:
   ```bash
   docker-compose up --build
   ```
4. Сервер будет доступен на `http://localhost:8080`.

### Тестирование
Для запуска тестов выполните:
```bash
go test ./...
```

## API-эндпоинты
### Посты
- **GET /posts**  
  Получить список всех постов.  
  **Ответ**: JSON-массив постов (`id`, `title`, `text`, `allow_comments`, `author`, `created_at`).
  **Пример**:
  ```json
  [
    {
      "id": 1,
      "title": "Test Post",
      "text": "This is a test post",
      "allow_comments": true,
      "author": "Author1",
      "created_at": "2025-06-09T13:15:00Z"
    }
  ]
  ```

- **POST /posts/create**  
  Создать новый пост.  
  **Тело запроса**:
  ```json
  {
    "title": "Test Post",
    "text": "This is a test post",
    "author": "Author1"
  }
  ```
  **Ответ**: JSON созданного поста.

- **POST /posts/disable-comments**  
  Отключить комментарии для поста.  
  **Тело запроса**:
  ```json
  {
    "post_id": 1
  }
  ```
  **Ответ**: Статус 200 при успехе.

### Комментарии
- **GET /comments?post_id=<ID>&limit=<N>&offset=<M>**  
  Получить комментарии для поста с пагинацией.  
  **Параметры**:
  - `post_id`: ID поста (обязательный).
  - `limit`: Количество комментариев (по умолчанию 10).
  - `offset`: Смещение для пагинации.  
  **Ответ**: JSON-массив комментариев (`id`, `post_id`, `parent_comment_id`, `text`, `author`, `created_at`).
  **Пример**:
  ```json
  [
    {
      "id": 1,
      "post_id": 1,
      "parent_comment_id": null,
      "text": "Great post!",
      "author": "User1",
      "created_at": "2025-06-09T13:15:00Z"
    }
  ]
  ```

- **POST /comments/create**  
  Создать новый комментарий.  
  **Тело запроса**:
  ```json
  {
    "post_id": 1,
    "parent_comment_id": null,
    "text": "Great post!",
    "author": "User1"
  }
  ```
  **Ответ**: JSON созданного комментария.  
  **Ограничения**: 
  - Текст не должен превышать 2000 символов.
  - Комментарии не создаются, если для поста отключены комментарии.

## Конфигурация
Конфигурация задаётся через `config.yaml` или переменные окружения:
- **server.host**: Хост сервера (по умолчанию `localhost`).
- **server.port**: Порт сервера (по умолчанию `8080`).
- **database.host**: Хост базы данных (например, `db` в Docker).
- **database.port**: Порт базы данных (по умолчанию `5432`).
- **database.user**: Пользователь базы данных.
- **database.password**: Пароль базы данных.
- **database.dbname**: Имя базы данных.

## Примечания
- Для PostgreSQL-хранилища требуется настроенная база данных и применённые миграции (выполняется автоматически при запуске с флагом `-storage=postgres`).
- In-memory хранилище подходит для тестирования и разработки, но не сохраняет данные после перезапуска.
- API возвращает соответствующие коды ошибок (400 для неверных запросов, 500 для внутренних ошибок).
