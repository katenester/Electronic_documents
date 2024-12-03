-- Таблица пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(255) NOT NULL UNIQUE,  -- Логин пользователя
                       password_hash TEXT NOT NULL,         -- Хеш пароля
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица документов
CREATE TABLE documents (
                           id SERIAL PRIMARY KEY,
                           owner_id INT NOT NULL REFERENCES users(id), -- Владелец документа (ссылка на пользователя)
                           name VARCHAR(255) NOT NULL,                 -- Имя документа (например, photo.jpg)
                           mime VARCHAR(100),                          -- MIME-тип документа
                           file BOOLEAN NOT NULL DEFAULT TRUE,         -- Флаг наличия файла
                           public BOOLEAN NOT NULL DEFAULT FALSE,      -- Флаг публичности документа
                           json_data JSONB,                            -- Дополнительные метаданные документа (необязательные)
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Дата создания документа
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица доступа к документам
CREATE TABLE document_grants (
                                 id SERIAL PRIMARY KEY,
                                 document_id INT NOT NULL REFERENCES documents(id), -- Ссылка на документ
                                 granted_to INT NOT NULL REFERENCES users(id),      -- Пользователь, получивший доступ
                                 UNIQUE (document_id, granted_to)                   -- Уникальность доступа к каждому документу
);

-- Таблица сессий пользователей
CREATE TABLE sessions (
                          id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL REFERENCES users(id),    -- Ссылка на пользователя
                          token TEXT NOT NULL UNIQUE,                   -- Токен сессии
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания сессии
                          expired_at TIMESTAMP                          -- Время истечения сессии (если используется)
);
