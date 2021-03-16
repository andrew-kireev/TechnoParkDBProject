CREATE TABLE IF NOT EXISTS users
(
    nickname varchar(100) NOT NULL primary key,
    fullname varchar(100) NOT NULL,
    about    text,
    email    varchar(100) NOT NULL unique
);

CREATE TABLE IF NOT EXISTS forum
(
    tittle        text         NOT NULL,
    user_nickname varchar(100) NOT NULL REFERENCES users (nickname),
    slug          text         NOT NULL PRIMARY KEY,
    posts         bigint,
    threads       bigint
);