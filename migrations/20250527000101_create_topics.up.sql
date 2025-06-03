CREATE TABLE IF NOT EXISTS backend_schema.topics (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    Description TEXT DEFAULT ''
);
