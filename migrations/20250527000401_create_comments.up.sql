CREATE TABLE backend_schema.comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES backend_schema.posts(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
