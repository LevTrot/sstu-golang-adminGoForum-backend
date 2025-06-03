CREATE TABLE IF NOT EXISTS backend_schema.posts (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES backend_schema.topics(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    username TEXT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);
