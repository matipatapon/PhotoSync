SET lc_messages TO 'en_US.UTF-8';

DROP TABLE IF EXISTS files CASCADE;
CREATE TABLE files
(
    id bigserial,
    user_id bigint REFERENCES users(id) NOT NULL,
    creation_date timestamp NOT NULL,
    filename text NOT NULL,
    mime_type smallint NOT NULL,
    file bytea NOT NULL,
    hash text NOT NULL,
    size bigint NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT "file is unique" UNIQUE (user_id, hash, size)
);

ALTER TABLE IF EXISTS files
    OWNER to postgres;

DROP TABLE IF EXISTS users CASCADE;
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
