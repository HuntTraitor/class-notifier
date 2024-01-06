CREATE TABLE IF NOT EXISTS users (
    userid SERIAL PRIMARY KEY,
    name CHAR(256) NOT NULL,
    email CHAR(256) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);