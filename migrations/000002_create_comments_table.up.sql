CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    parent_comment_id INTEGER REFERENCES comments(id),
    text TEXT NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);