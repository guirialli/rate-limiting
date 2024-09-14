DROP TABLE IF EXISTS books;

CREATE TABLE IF NOT EXISTS books
(
    id          varchar(255) PRIMARY KEY,
    title       varchar(255) NOT NULL,
    pages       int,
    description TEXT,
    author_id   varchar(255) NOT NULL
);
