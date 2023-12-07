CREATE TABLE IF NOT EXISTS users
(
    uid bytea NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL DEFAULT '',
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
    );