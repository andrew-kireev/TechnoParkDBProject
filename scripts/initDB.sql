CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    id       bigserial,
    nickname citext NOT NULL primary key,
    fullname text   NOT NULL,
    about    text,
    email    citext NOT NULL unique
);

CREATE UNLOGGED TABLE IF NOT EXISTS forum
(
    title         text   NOT NULL,
    user_nickname citext NOT NULL REFERENCES users (nickname),
    slug          citext NOT NULL PRIMARY KEY,
    posts         bigint default 0,
    threads       int    default 0
);


CREATE UNLOGGED TABLE IF NOT EXISTS threads
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

create unlogged table posts
(
    id        bigserial
        constraint posts_pkey primary key,
    parent    integer                  default 0,
    author    citext not null
        constraint posts_author_fkey references users,
    message   text   not null,
    is_edited boolean                  default false,
    forum     citext
        constraint posts_forum_fkey references forum,
    thread    integer
        constraint posts_thread_fkey references threads,
    created   timestamp with time zone default now(),
    path      bigint[]                 default ARRAY []::integer[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes
(
    nickname  citext not null references users (nickname),
    thread_id int    not null references threads (id),
    voice     int    not null,
    unique (nickname, thread_id)
);

create unlogged table if not exists users_to_forums
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
    VALUES (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
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


CREATE OR REPLACE FUNCTION path() RETURNS TRIGGER AS
$path$
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
    RETURN new;
END;
$path$ LANGUAGE plpgsql;

CREATE TRIGGER path_tri
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE path();


create index if not exists thread_slug_hash on threads using hash (slug);
create index if not exists thread_forum_created on threads (forum, created);
create index if not exists thread_user on threads using hash (author);
create index if not exists thread_full on threads (slug, id, forum);
create index if not exists thread_all on threads (title, author, forum, message, votes, slug, created);

create index if not exists post_pathparent on posts ((path[1]));
create index if not exists posts_thread_thread_id on posts (thread, id);
create index if not exists post_path_id on posts (id, (path[1]));
create index if not exists post_parent on posts (thread, id, (path[1]), parent);
create index if not exists posts_sorting on posts using gin ((path[1]) desc, path, id);
create index if not exists post_thread on posts (thread);
create index if not exists post_gin on posts using gin (path);
create index if not exists posts_thread_path_id on posts (thread, path, id);

create index if not exists users_nickname_hash on users using hash (nickname);
create index if not exists users_full on users  (nickname, fullname, about, email);

create unique index if not exists votes on votes (nickname, thread_id, voice);

create index if not exists forum_slug on forum using hash (slug);
create index if not exists forum_full on forum (title, user_nickname, slug, posts, threads);
create index if not exists forum_user on forum using hash (user_nickname);
create index if not exists f_u_nickname on users_to_forums using hash (nickname);
create index if not exists users_to_forums_forum_nickname on users_to_forums (forum, nickname);
create index if not exists users_to_forums_forum on users_to_forums (forum);
