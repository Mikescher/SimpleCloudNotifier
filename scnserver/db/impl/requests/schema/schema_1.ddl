
CREATE TABLE `requests`
(
    request_id           INTEGER        PRIMARY KEY,
    timestamp_created    INTEGER        NOT NULL

) STRICT;


CREATE TABLE `meta`
(
    meta_key     TEXT       NOT NULL,
    value_int    INTEGER        NULL,
    value_txt    TEXT           NULL,
    value_real   REAL           NULL,
    value_blob   BLOB           NULL,

    PRIMARY KEY (meta_key)
) STRICT;


INSERT INTO meta (meta_key, value_int) VALUES ('schema', 1)