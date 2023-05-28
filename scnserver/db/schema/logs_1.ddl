
CREATE TABLE `logs`
(
    log_id               TEXT           NOT NULL,
    timestamp_created    INTEGER        NOT NULL,

    PRIMARY KEY (log_id)
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