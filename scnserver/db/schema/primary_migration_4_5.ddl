
ALTER TABLE clients ADD COLUMN "name" TEXT NULL;

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

    PRIMARY KEY (client_id)
) STRICT;

INSERT INTO clients_new SELECT * FROM clients;
DROP TABLE clients;
ALTER TABLE clients_new RENAME TO clients;


CREATE        INDEX "idx_clients_userid"   ON clients (user_id);
CREATE UNIQUE INDEX "idx_clients_fcmtoken" ON clients (fcm_token);


