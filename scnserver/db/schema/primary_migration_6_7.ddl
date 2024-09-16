


DROP INDEX "idx_clients_userid";
DROP INDEX "idx_clients_fcmtoken";


CREATE TABLE clients_new
(
    client_id          TEXT                                                                    NOT NULL,

    user_id            TEXT                                                                    NOT NULL,
    type               TEXT       CHECK(type IN ('ANDROID','IOS','LINUX','MACOS','WINDOWS'))   NOT NULL,
    fcm_token          TEXT                                                                    NOT NULL,
    name               TEXT                                                                        NULL,

    timestamp_created  INTEGER                                                                 NOT NULL,

    agent_model        TEXT                                                                    NOT NULL,
    agent_version      TEXT                                                                    NOT NULL,


    deleted               INTEGER  CHECK(deleted IN (0, 1))        NOT NULL    DEFAULT '0',

    PRIMARY KEY (client_id)
) STRICT;


INSERT INTO clients_new
SELECT
    client_id,
    user_id,
    type,
    fcm_token,
    name,
    timestamp_created,
    agent_model,
    agent_version,
    0 AS deleted
FROM clients;


DROP TABLE clients;
ALTER TABLE clients_new RENAME TO clients;


CREATE        INDEX "idx_clients_userid"   ON clients (user_id);
CREATE        INDEX "idx_clients_deleted"  ON clients (deleted);
CREATE UNIQUE INDEX "idx_clients_fcmtoken" ON clients (fcm_token);



