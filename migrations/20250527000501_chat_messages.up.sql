CREATE TABLE backend_schema.chat_messages (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
