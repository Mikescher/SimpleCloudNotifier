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
);
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
);
CREATE        INDEX "idx_clients_userid"   ON clients (user_id);
CREATE UNIQUE INDEX "idx_clients_fcmtoken" ON clients (fcm_token);


CREATE TABLE channels
(
    channel_id         INTEGER     PRIMARY KEY AUTOINCREMENT,

    owner_user_id      INTEGER     NOT NULL,

    name               TEXT        NOT NULL,

    subscribe_key      TEXT        NOT NULL,
    send_key           TEXT        NOT NULL,

    timestamp_created  INTEGER     NOT NULL,
    timestamp_lastread INTEGER         NULL   DEFAULT NULL,
    timestamp_lastsent INTEGER         NULL   DEFAULT NULL,

    messages_sent      INTEGER     NOT NULL   DEFAULT '0'
);
CREATE UNIQUE INDEX "idx_channels_identity" ON channels (owner_user_id, name);

CREATE TABLE subscriptions
(
    subscription_id        INTEGER                                PRIMARY KEY AUTOINCREMENT,

    subscriber_user_id     INTEGER                                NOT NULL,
    channel_owner_user_id  INTEGER                                NOT NULL,
    channel_name           TEXT                                   NOT NULL,
    channel_id             INTEGER                                NOT NULL,

    timestamp_created      INTEGER                                NOT NULL,

    confirmed              INTEGER   CHECK(confirmed IN (0, 1))   NOT NULL
);
CREATE UNIQUE INDEX "idx_subscriptions_ref" ON subscriptions (subscriber_user_id, channel_owner_user_id, channel_name);


CREATE TABLE messages
(
    scn_message_id     INTEGER                                  PRIMARY KEY AUTOINCREMENT,
    sender_user_id     INTEGER                                  NOT NULL,
    owner_user_id      INTEGER                                  NOT NULL,
    channel_name       TEXT                                     NOT NULL,
    channel_id         INTEGER                                  NOT NULL,

    timestamp_real     INTEGER                                  NOT NULL,
    timestamp_client   INTEGER                                      NULL,

    title              TEXT                                     NOT NULL,
    content            TEXT                                         NULL,
    priority           INTEGER  CHECK(priority IN (0, 1, 2))    NOT NULL,
    usr_message_id     TEXT                                         NULL
);
CREATE INDEX "idx_messages_channel"     ON messages (owner_user_id, channel_name);
CREATE INDEX "idx_messages_idempotency" ON messages (owner_user_id, usr_message_id);


CREATE TABLE deliveries
(
    delivery_id         INTEGER                                                  PRIMARY KEY AUTOINCREMENT,

    scn_message_id      INTEGER                                                  NOT NULL,
    receiver_user_id    INTEGER                                                  NOT NULL,
    receiver_client_id  INTEGER                                                  NOT NULL,

    timestamp_created   INTEGER                                                  NOT NULL,
    timestamp_finalized INTEGER                                                  NOT NULL,


    status              TEXT     CHECK(status IN ('RETRY','SUCCESS','FAILED'))   NOT NULL,
    retry_count         INTEGER                                                  NOT NULL   DEFAULT 0,
    next_delivery       INTEGER                                                      NULL   DEFAULT NULL,

    fcm_message_id      TEXT                                                         NULL
);
CREATE INDEX "idx_deliveries_receiver" ON deliveries (scn_message_id, receiver_client_id);


CREATE TABLE `meta`
(
    meta_key     TEXT       NOT NULL,
    value_int    INTEGER        NULL,
    value_txt    TEXT           NULL,
    value_real   REAL           NULL,
    value_blob   BLOB           NULL,

    PRIMARY KEY (meta_key)
);
INSERT INTO meta (meta_key, value_int) VALUES ('schema', 3)