CREATE TABLE users
(
    user_id            INTEGER                             PRIMARY KEY AUTOINCREMENT,

    username           TEXT                                    NULL  DEFAULT NULL,

    send_key           TEXT                                NOT NULL,
    read_key           TEXT                                NOT NULL,
    admin_key          TEXT                                NOT NULL,

    timestamp_created  INTEGER                             NOT NULL,
    timestamp_lastread INTEGER                                 NULL   DEFAULT NULL,
    timestamp_lastsent INTEGER                                 NULL   DEFAULT NULL,

    messages_sent      INTEGER                             NOT NULL   DEFAULT '0',

    quota_used         INTEGER                             NOT NULL   DEFAULT '0',
    quota_used_day     TEXT                                    NULL   DEFAULT NULL,

    is_pro             INTEGER   CHECK(is_pro IN (0, 1))   NOT NULL   DEFAULT 0,
    pro_token          TEXT                                    NULL   DEFAULT NULL
) STRICT;
CREATE  UNIQUE INDEX "idx_users_protoken" ON users (pro_token) WHERE pro_token IS NOT NULL;


CREATE TABLE clients
(
    client_id          INTEGER                                        PRIMARY KEY AUTOINCREMENT,

    user_id            INTEGER                                        NOT NULL,
    type               TEXT       CHECK(type IN ('ANDROID', 'IOS'))   NOT NULL,
    fcm_token          TEXT                                               NULL,

    timestamp_created  INTEGER                                        NOT NULL,

    agent_model        TEXT                                           NOT NULL,
    agent_version      TEXT                                           NOT NULL
) STRICT;
CREATE        INDEX "idx_clients_userid"   ON clients (user_id);
CREATE UNIQUE INDEX "idx_clients_fcmtoken" ON clients (fcm_token);


CREATE TABLE channels
(
    channel_id         INTEGER     PRIMARY KEY AUTOINCREMENT,

    owner_user_id      INTEGER     NOT NULL,

    internal_name      TEXT        NOT NULL,
    display_name       TEXT        NOT NULL,

    subscribe_key      TEXT        NOT NULL,
    send_key           TEXT        NOT NULL,

    timestamp_created  INTEGER     NOT NULL,
    timestamp_lastsent INTEGER         NULL   DEFAULT NULL,

    messages_sent      INTEGER     NOT NULL   DEFAULT '0'
) STRICT;
CREATE UNIQUE INDEX "idx_channels_identity" ON channels (owner_user_id, internal_name);

CREATE TABLE subscriptions
(
    subscription_id        INTEGER                                PRIMARY KEY AUTOINCREMENT,

    subscriber_user_id     INTEGER                                NOT NULL,
    channel_owner_user_id  INTEGER                                NOT NULL,
    channel_internal_name  TEXT                                   NOT NULL,
    channel_id             INTEGER                                NOT NULL,

    timestamp_created      INTEGER                                NOT NULL,

    confirmed              INTEGER   CHECK(confirmed IN (0, 1))   NOT NULL
) STRICT;
CREATE UNIQUE INDEX "idx_subscriptions_ref"     ON subscriptions (subscriber_user_id, channel_owner_user_id, channel_internal_name);
CREATE        INDEX "idx_subscriptions_chan"    ON subscriptions (channel_id);
CREATE        INDEX "idx_subscriptions_subuser" ON subscriptions (subscriber_user_id);
CREATE        INDEX "idx_subscriptions_ownuser" ON subscriptions (channel_owner_user_id);
CREATE        INDEX "idx_subscriptions_tsc"     ON subscriptions (timestamp_created);
CREATE        INDEX "idx_subscriptions_conf"    ON subscriptions (confirmed);


CREATE TABLE messages
(
    scn_message_id        INTEGER                                  PRIMARY KEY AUTOINCREMENT,
    sender_user_id        INTEGER                                  NOT NULL,
    owner_user_id         INTEGER                                  NOT NULL,
    channel_internal_name TEXT                                     NOT NULL,
    channel_id            INTEGER                                  NOT NULL,
    sender_ip             TEXT                                     NOT NULL,
    sender_name           TEXT                                         NULL,

    timestamp_real        INTEGER                                  NOT NULL,
    timestamp_client      INTEGER                                      NULL,

    title                 TEXT                                     NOT NULL,
    content               TEXT                                         NULL,
    priority              INTEGER  CHECK(priority IN (0, 1, 2))    NOT NULL,
    usr_message_id        TEXT                                         NULL,

    deleted               INTEGER  CHECK(deleted IN (0, 1))        NOT NULL    DEFAULT '0'
) STRICT;
CREATE        INDEX "idx_messages_owner_channel"    ON messages (owner_user_id, channel_internal_name COLLATE BINARY);
CREATE        INDEX "idx_messages_owner_channel_nc" ON messages (owner_user_id, channel_internal_name COLLATE NOCASE);
CREATE        INDEX "idx_messages_channel"          ON messages (channel_internal_name COLLATE BINARY);
CREATE        INDEX "idx_messages_channel_nc"       ON messages (channel_internal_name COLLATE NOCASE);
CREATE UNIQUE INDEX "idx_messages_idempotency"      ON messages (owner_user_id, usr_message_id COLLATE BINARY);
CREATE        INDEX "idx_messages_senderip"         ON messages (sender_ip COLLATE BINARY);
CREATE        INDEX "idx_messages_sendername"       ON messages (sender_name COLLATE BINARY);
CREATE        INDEX "idx_messages_sendername_nc"    ON messages (sender_name COLLATE NOCASE);
CREATE        INDEX "idx_messages_title"            ON messages (title COLLATE BINARY);
CREATE        INDEX "idx_messages_title_nc"         ON messages (title COLLATE NOCASE);
CREATE        INDEX "idx_messages_deleted"          ON messages (deleted);


CREATE VIRTUAL TABLE messages_fts USING fts5
(
    channel_internal_name,
    sender_name,
    title,
    content,

    tokenize = unicode61,
    content = 'messages',
    content_rowid = 'scn_message_id'
);

CREATE TRIGGER fts_insert AFTER INSERT ON messages BEGIN
    INSERT INTO messages_fts (rowid, channel_internal_name, sender_name, title, content) VALUES (new.scn_message_id, new.channel_internal_name, new.sender_name, new.title, new.content);
END;

CREATE TRIGGER fts_update AFTER UPDATE ON messages BEGIN
    INSERT INTO messages_fts (messages_fts, rowid, channel_internal_name, sender_name, title, content) VALUES ('delete', old.scn_message_id, old.channel_internal_name, old.sender_name, old.title, old.content);
    INSERT INTO messages_fts (              rowid, channel_internal_name, sender_name, title, content) VALUES (          new.scn_message_id, new.channel_internal_name, new.sender_name, new.title, new.content);
END;

CREATE TRIGGER fts_delete AFTER DELETE ON messages BEGIN
    INSERT INTO messages_fts (messages_fts, rowid, channel_internal_name, sender_name, title, content) VALUES ('delete', old.scn_message_id, old.channel_internal_name, old.sender_name, old.title, old.content);
END;



CREATE TABLE deliveries
(
    delivery_id         INTEGER                                                  PRIMARY KEY AUTOINCREMENT,

    scn_message_id      INTEGER                                                  NOT NULL,
    receiver_user_id    INTEGER                                                  NOT NULL,
    receiver_client_id  INTEGER                                                  NOT NULL,

    timestamp_created   INTEGER                                                  NOT NULL,
    timestamp_finalized INTEGER                                                      NULL,


    status              TEXT     CHECK(status IN ('RETRY','SUCCESS','FAILED'))   NOT NULL,
    retry_count         INTEGER                                                  NOT NULL   DEFAULT 0,
    next_delivery       INTEGER                                                      NULL   DEFAULT NULL,

    fcm_message_id      TEXT                                                         NULL
) STRICT;
CREATE INDEX "idx_deliveries_receiver" ON deliveries (scn_message_id, receiver_client_id);


CREATE TABLE `meta`
(
    meta_key     TEXT       NOT NULL,
    value_int    INTEGER        NULL,
    value_txt    TEXT           NULL,
    value_real   REAL           NULL,
    value_blob   BLOB           NULL,

    PRIMARY KEY (meta_key)
) STRICT;


INSERT INTO meta (meta_key, value_int) VALUES ('schema', 3)