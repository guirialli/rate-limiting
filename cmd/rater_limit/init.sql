CREATE TABLE IF NOT EXISTS authors
(
    id          varchar(255) PRIMARY KEY,
    name        varchar(255) NOT NULL,
    birthday    datetime,
    description text
);
CREATE TABLE IF NOT EXISTS books
(
    id          varchar(255) PRIMARY KEY,
    title       varchar(255) NOT NULL,
    pages       int          NOT NULL,
    description TEXT,
    author_id   varchar(255) NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors (id)
);
CREATE TABLE IF NOT EXISTS migrates
(
    id      int AUTO_INCREMENT primary key,
    version varchar(64)
);
