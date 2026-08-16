# Решение

## Ожидаемое решение (ориентир)

```sql
CREATE TABLE users (
    id            bigserial PRIMARY KEY,
    name          text NOT NULL,
    registered_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE chats (
    id         bigserial PRIMARY KEY,
    name       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- связь many-to-many между пользователями и чатами
CREATE TABLE chat_users (
    chat_id   bigint NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id   bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
    id         bigserial PRIMARY KEY,
    chat_id    bigint NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    author_id  bigint NOT NULL REFERENCES users(id),
    text       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON chat_users (user_id);
CREATE INDEX ON messages (chat_id, created_at DESC);
```

Запрос:

```sql
SELECT c.id AS chat_id, c.name AS chat_name
FROM chats c
JOIN chat_users cu ON cu.chat_id = c.id
JOIN users u ON u.id = cu.user_id
WHERE u.name = 'Вася';
```

Обсуждаемые темы: нормализация, индексы, что делать при миллионах сообщений
(партиционирование по `chat_id`/времени), почему `message` ссылается на чат
напрямую (правило 2 и 3).

Схема и запрос отдельным файлом: [`schema.sql`](schema.sql).
