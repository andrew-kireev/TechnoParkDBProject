CREATE TABLE IF NOT EXISTS users
(
    user_id  bigserial PRIMARY KEY,
    nickname varchar(100) NOT NULL unique,
    fullname varchar(100) NOT NULL,
    email    varchar(100) NOT NULL unique
);