DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;

CREATE TABLE IF NOT EXISTS books
(
    id          varchar(255) PRIMARY KEY,
    title       varchar(255) NOT NULL,
    pages       int,
    description TEXT,
    author_id   varchar(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS authors (
    id varchar(255) PRIMARY KEY,
    name varchar(255) NOT NULL,
    birthday datetime,
    description TEXT
);
CREATE TABLE IF NOT EXISTS users
(
    id varchar(255) PRIMARY KEY,
    username varchar(64) NOT NULL UNIQUE,
    password TEXT NOT NULL
);
