CREATE TABLE `users`
(
    `user_id`            INTEGER                AUTO_INCREMENT,
    `user_key`           TEXT         NOT NULL,
    `fcm_token`          TEXT             NULL  DEFAULT NULL,
    `messages_sent`      INTEGER     NOT NULL   DEFAULT '0',
    `timestamp_created`  TEXT        NOT NULL   DEFAULT CURRENT_TIMESTAMP,
    `timestamp_accessed` TEXT            NULL   DEFAULT NULL,

    `quota_today`        INTEGER     NOT NULL   DEFAULT '0',
    `quota_day`          TEXT            NULL   DEFAULT NULL,

    `is_pro`             INTEGER     NOT NULL   DEFAULT 0,
    `pro_token`          TEXT            NULL   DEFAULT NULL,

    PRIMARY KEY (`user_id`)
);

CREATE TABLE `messages`
(
    `scn_message_id`     INTEGER                AUTO_INCREMENT,
    `sender_user_id`     INTEGER     NOT NULL,

    `timestamp_real`     TEXT        NOT NULL   DEFAULT CURRENT_TIMESTAMP,
    `ack`                INTEGER     NOT NULL   DEFAULT 0,

    `title`              TEXT        NOT NULL,
    `content`            TEXT            NULL,
    `priority`           INTEGER     NOT NULL,
    `sendtime`           INTEGER     NOT NULL,

    `fcm_message_id`     TEXT            NULL,
    `usr_message_id`     TEXT            NULL,

    PRIMARY KEY (`scn_message_id`)
);

CREATE TABLE `meta`
(
    `key`         TEXT     NOT NULL,
    `value_int`   INTEGER      NULL,
    `value_txt`   TEXT         NULL,

    PRIMARY KEY (`key`)
);

INSERT INTO meta (key, value_int) VALUES ('schema', 2)