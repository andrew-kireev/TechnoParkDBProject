CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    nickname citext       NOT NULL primary key,
    fullname varchar(100) NOT NULL,
    about    text,
    email    citext       NOT NULL unique
);

CREATE TABLE IF NOT EXISTS forum
(
    title         text   NOT NULL,
    user_nickname citext NOT NULL REFERENCES users (nickname),
    slug          citext NOT NULL PRIMARY KEY,
    posts         bigint,
    threads       int
);

CREATE TABLE IF NOT EXISTS threads
(
    id      bigserial not null primary key,
    title   text      not null,
    author  citext    not null references users (nickname),
    forum   citext    not null references forum (slug),
    message text      not null,
    votes   integer                  default 0,
    slug    citext unique,
    created timestamp with time zone default now()
);