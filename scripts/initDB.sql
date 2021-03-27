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

CREATE TABLE IF NOT EXISTS posts
(
    id        bigserial not null primary key,
    parent    int                      default 0,
    author    citext    not null references users (nickname),
    message   text      not null,
    is_edited boolean                  default false,
    forum     citext    not null references forum (slug),
    thread    int       not null references threads (id),
    created   timestamp with time zone default now()
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
    slug    citext,
    created timestamp with time zone default now()
);

CREATE TABLE IF NOT EXISTS posts
(
    id        bigserial not null primary key,
    parent    int                      default 0,
    author    citext    not null references users (nickname),
    message   text      not null,
    is_edited boolean                  default false,
    forum     citext    not null references forum (slug),
    thread    int       not null references threads (id),
    created   timestamp with time zone default now()
);


SELECT t.id,
       t.title,
       t.author,
       t.forum,
       t.message,
       t.votes,
       t.slug,
       t.created
from threads as t
         LEFT JOIN forum f on t.forum = f.slug
WHERE f.slug = 'lm6sJN9tBRIfSe'
  and t.created >= '2020-03-21T08:03:48.617Z'
ORDER BY t.created desc
LIMIT NULLIF(4, 0);

SELECT t.id,
       t.title,
       t.author,
       t.forum,
       t.message,
       t.votes,
       t.slug,
       t.created
from threads as t
         LEFT JOIN forum f on t.forum = f.slug
WHERE f.slug = 'JL894__GBR-5k'
  and t.created >= '2020-10-28 01:10:32.234'
ORDER BY t.created asc
LIMIT NULLIF(4, 0);


SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads as t
LEFT JOIN forum f on t.forum = f.slug
WHERE f.slug = 'urJsg6V2PS-cr' and t.created >= '2020-08-18T16:46:40.418Z'
ORDER BY t.created asc
LIMIT NULLIF(4, 0);


INSERT INTO posts (parent, author, message, forum, thread)
VALUES (0, 'cuiusque.ntO4uzL1nRz571', 'fdfds', 'SeeML50Rpk-Jr', 625),
       (1, 'cuiusque.ntO4uzL1nRz571', 'SeeML50Rpk-Jr', 'SeeML50Rpk-Jr', 700)
returning id, parent, author, message, is_edited, forum, thread, created;


SELECT f.title, f.user_nickname, f.slug, f.posts, f.threads FROM forum as f
left join threads t on f.slug = t.forum
where t.id = 1743;


CREATE TABLE IF NOT EXISTS votes
(
    nickname  citext not null references users (nickname),
    thread_id int    not null references threads (id),
    voice     int    not null,
    unique (nickname, thread_id)
);

create table if not exists users_to_forums
(
    nickname citext not null references users (nickname),
    forum    citext not null references forum (slug),
    unique (nickname, forum)
);

CREATE OR REPLACE FUNCTION insert_votes_threads()
    RETURNS TRIGGER AS
$insert_votes_threads$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
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

CREATE OR REPLACE FUNCTION update_votes_threads()
    RETURNS TRIGGER AS
$update_votes_threads$
BEGIN
    UPDATE threads
    SET votes=votes + 2 * NEW.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$update_votes_threads$ LANGUAGE plpgsql;



DROP TRIGGER IF EXISTS update_votes_threads ON votes;
CREATE TRIGGER update_votes_threads
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes_threads();

SELECT id, parent, author, message, is_edited, forum, thread, created from posts
WHERE thread = 3471
ORDER BY created, id asc
LIMIT NULLIF(65, 0);


CREATE OR REPLACE FUNCTION count_posts()
    RETURNS TRIGGER AS
$count_posts$
BEGIN
    UPDATE forum
    SET posts = forum.posts + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$count_posts$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS count_posts ON forum;
CREATE TRIGGER count_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE count_posts();

CREATE OR REPLACE FUNCTION count_threads()
    RETURNS TRIGGER AS
$count_threads$
BEGIN
    UPDATE forum
    SET threads = forum.threads + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$count_threads$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS count_threads ON threads;
CREATE TRIGGER count_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE count_threads();


CREATE OR REPLACE FUNCTION update_users_forum()
    RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    INSERT INTO users_to_forums (nickname, forum)
    VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING ;
    RETURN NEW;
END;
$update_users_forum$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_users_forum ON threads;
CREATE TRIGGER update_users_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_users_forum();

DROP TRIGGER IF EXISTS update_users_forum ON posts;
CREATE TRIGGER update_users_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_users_forum();


