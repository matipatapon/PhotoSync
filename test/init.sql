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