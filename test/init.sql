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

DROP TABLE IF EXISTS postgres_database_test_empty_table;
CREATE TABLE postgres_database_test_empty_table(
    id integer not null,
    name text not null
);

DROP TABLE IF EXISTS postgres_database_test_table_with_one_item;
CREATE TABLE postgres_database_test_table_with_one_item(
    id integer not null PRIMARY KEY,
    name text not null
);
INSERT INTO postgres_database_test_table_with_one_item VALUES(1, 'Mort');

DROP TABLE IF EXISTS postgres_database_test_table_with_two_items;
CREATE TABLE postgres_database_test_table_with_two_items(
    id integer not null PRIMARY KEY,
    name text not null
);
INSERT INTO postgres_database_test_table_with_two_items VALUES(1, 'Mort');
INSERT INTO postgres_database_test_table_with_two_items VALUES(2, 'Luna');

DROP TABLE IF EXISTS postgres_database_test_insertion_table;
CREATE TABLE postgres_database_test_insertion_table(
    id integer not null PRIMARY KEY,
    name text not null
);

DROP TABLE IF EXISTS postgres_database_test_table_to_drop;
CREATE TABLE postgres_database_test_table_to_drop(
    id integer not null PRIMARY KEY,
    name text not null
);

DROP TABLE IF EXISTS postgres_database_test_table_to_update;
CREATE TABLE postgres_database_test_table_to_update(
    id integer not null PRIMARY KEY,
    name text not null
);
INSERT INTO postgres_database_test_table_to_update VALUES(1, 'orginal_name');

DROP TABLE IF EXISTS postgres_database_test_table_to_delete;
CREATE TABLE postgres_database_test_table_to_delete(
    id integer not null PRIMARY KEY,
    name text not null
);
INSERT INTO postgres_database_test_table_to_delete VALUES(1, 'name');
