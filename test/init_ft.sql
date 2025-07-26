DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id bigserial,
    username text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT "username is unique" UNIQUE (username)
);

ALTER TABLE IF EXISTS users
    OWNER to postgres;
