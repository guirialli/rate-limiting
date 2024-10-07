DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS books
(
    id          varchar(255) PRIMARY KEY,
    title       varchar(255) NOT NULL,
    pages       int,
    description TEXT,
    author_id   varchar(255) NOT NULL /* I removed the author-book relationship for isolated testing. */
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
