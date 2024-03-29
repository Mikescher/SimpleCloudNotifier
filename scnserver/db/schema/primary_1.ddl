DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `user_id`            INT(11)         NOT NULL   AUTO_INCREMENT,
    `user_key`           VARCHAR(64)     NOT NULL,
    `fcm_token`          VARCHAR(256)        NULL   DEFAULT NULL,
    `messages_sent`      INT(11)         NOT NULL   DEFAULT '0',
    `timestamp_created`  DATETIME        NOT NULL   DEFAULT CURRENT_TIMESTAMP,
    `timestamp_accessed` DATETIME            NULL   DEFAULT NULL,

    `quota_today`        INT(11)         NOT NULL   DEFAULT '0',
    `quota_day`          DATE                NULL   DEFAULT NULL,

    `is_pro`             BIT             NOT NULL   DEFAULT 0,
    `pro_token`          VARCHAR(256)        NULL   DEFAULT NULL,

    PRIMARY KEY (`user_id`)
);

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `message_id`     INT(11)         NOT NULL   AUTO_INCREMENT,
    `sender_user_id`     INT(11)         NOT NULL,

    `timestamp_real`     DATETIME        NOT NULL   DEFAULT CURRENT_TIMESTAMP,
    `ack`                TINYINT(1)      NOT NULL   DEFAULT 0,

    `title`              VARCHAR(256)    NOT NULL,
    `content`            LONGTEXT            NULL,
    `priority`           INT(11)         NOT NULL,
    `sendtime`           BIGINT UNSIGNED NOT NULL,

    `fcm_message_id`     VARCHAR(256)        NULL,
    `usr_message_id`     VARCHAR(256)        NULL,

    PRIMARY KEY (`message_id`)
);