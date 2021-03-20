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

CREATE TABLE IF NOT EXISTS votes
(
    nickname  citext not null references users (nickname),
    thread_id int    not null references threads (id),
    voice     int    not null,
    unique (nickname, thread_id)
);

CREATE OR REPLACE FUNCTION insert_votes_threads()
    RETURNS TRIGGER AS
$insert_votes_threads$
BEGIN
    UPDATE threads
    SET votes = NEW.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$insert_votes_threads$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_votes_threads ON votes;
CREATE TRIGGER insert_votes_threads
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes_threads();