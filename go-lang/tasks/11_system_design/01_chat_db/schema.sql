-- Схема БД для чата: пользователи, чаты, сообщения.
--
-- Связи:
--   users <-> chats  : many-to-many через chat_users
--                      (правило 1: пользователь состоит в нескольких чатах);
--   messages -> chats: many-to-one, chat_id NOT NULL
--                      (правила 2 и 3: сообщение принадлежит ровно одному чату);
--   messages -> users: автор сообщения.

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

-- many-to-many: участие пользователей в чатах
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

-- Индексы под типовые запросы:
CREATE INDEX chat_users_user_idx ON chat_users (user_id);                -- «чаты пользователя»
CREATE INDEX messages_chat_time_idx ON messages (chat_id, created_at DESC); -- «лента чата»

-- Задание 2: все чаты пользователя "Вася" в формате (chat_id, chat_name)
SELECT c.id AS chat_id, c.name AS chat_name
FROM chats c
JOIN chat_users cu ON cu.chat_id = c.id
JOIN users u ON u.id = cu.user_id
WHERE u.name = 'Вася';
