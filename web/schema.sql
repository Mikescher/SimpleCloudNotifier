CREATE TABLE `users`
(
	`id`                 INT(11)      NOT NULL   AUTO_INCREMENT,
	`auth_token`         VARCHAR(64)  NOT NULL,
	`fcm_token`          VARCHAR(256)     NULL   DEFAULT NULL,
	`messages_sent`      INT(11)      NOT NULL   DEFAULT '0',
	`timestamp_created`  DATETIME     NOT NULL   DEFAULT CURRENT_TIMESTAMP,
	`timestamp_accessed` DATETIME         NULL   DEFAULT NULL,

	`quota_today`        INT(11)      NOT NULL   DEFAULT '0',
	`quota_day`          DATE             NULL   DEFAULT NULL,
	`quota_max`          INT(11)      NOT NULL   DEFAULT '100',

	PRIMARY KEY (`id`)
);
