CREATE TABLE `users`
(
	`user_id`            INT(11)      NOT NULL   AUTO_INCREMENT,
	`user_key`           VARCHAR(64)  NOT NULL,
	`fcm_token`          VARCHAR(256)     NULL   DEFAULT NULL,
	`messages_sent`      INT(11)      NOT NULL   DEFAULT '0',
	`timestamp_created`  DATETIME     NOT NULL   DEFAULT CURRENT_TIMESTAMP,
	`timestamp_accessed` DATETIME         NULL   DEFAULT NULL,

	`quota_today`        INT(11)      NOT NULL   DEFAULT '0',
	`quota_day`          DATE             NULL   DEFAULT NULL,

	`is_pro`             BIT          NOT NULL   DEFAULT 0,
	`pro_token`          VARCHAR(256)     NULL   DEFAULT NULL,

	PRIMARY KEY (`user_id`)
);
