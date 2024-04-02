CREATE SEQUENCE test_seq;
CREATE TABLE test_tbl (
    ID NUMERIC PRIMARY KEY NOT NULL DEFAULT NEXTVAL('test_seq'),
    label VARCHAR
);

INSERT INTO test_tbl (label) VALUES ('One');
INSERT INTO test_tbl (label) VALUES ('Two');
