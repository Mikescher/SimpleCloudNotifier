


DROP INDEX "idx_clients_fcmtoken";

CREATE UNIQUE INDEX "idx_clients_fcmtoken" ON clients (fcm_token) WHERE deleted=0;

