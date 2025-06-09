CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    allow_comments BOOLEAN NOT NULL,
    author TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);