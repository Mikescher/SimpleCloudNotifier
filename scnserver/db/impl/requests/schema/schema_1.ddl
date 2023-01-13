
CREATE TABLE `requests`
(
    request_id            TEXT                                   NOT NULL,

    method                TEXT                                   NOT NULL,
    uri                   TEXT                                   NOT NULL,
    user_agent            TEXT                                       NULL,
    authentication        TEXT                                       NULL,
    request_body          TEXT                                       NULL,
    request_body_size     INTEGER                                NOT NULL,
    request_content_type  TEXT                                   NOT NULL,
    remote_ip             TEXT                                   NOT NULL,

    userid                TEXT                                       NULL,
    permissions           TEXT                                       NULL,

    response_statuscode   INTEGER                                    NULL,
    response_body_size    INTEGER                                    NULL,
    response_body         TEXT                                       NULL,
    response_content_type TEXT                                   NOT NULL,
    processing_time       INTEGER                                NOT NULL,
    retry_count           INTEGER                                NOT NULL,
    panicked              INTEGER    CHECK(panicked IN (0, 1))   NOT NULL,
    panic_str             TEXT                                       NULL,

    timestamp_created     INTEGER                                NOT NULL,
    timestamp_start       INTEGER                                NOT NULL,
    timestamp_finish      INTEGER                                NOT NULL,

    PRIMARY KEY (request_id)
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