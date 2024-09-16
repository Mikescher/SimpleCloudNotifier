
DROP INDEX "idx_deliveries_receiver";

CREATE TABLE deliveries_new
(
    delivery_id         TEXT                                                     NOT NULL,

    message_id          TEXT                                                     NOT NULL,
    receiver_user_id    TEXT                                                     NOT NULL,
    receiver_client_id  TEXT                                                     NOT NULL,

    timestamp_created   INTEGER                                                  NOT NULL,
    timestamp_finalized INTEGER                                                      NULL,


    status              TEXT     CHECK(status IN ('RETRY','SUCCESS','FAILED'))   NOT NULL,
    retry_count         INTEGER                                                  NOT NULL   DEFAULT 0,
    next_delivery       INTEGER                                                      NULL   DEFAULT NULL,

    fcm_message_id      TEXT                                                         NULL,

    PRIMARY KEY (delivery_id)
) STRICT;

UPDATE deliveries SET next_delivery  = NULL;

INSERT INTO deliveries_new
SELECT
    delivery_id,
    message_id,
    receiver_user_id,
    receiver_client_id,
    timestamp_created,
    timestamp_finalized,
    status,
    retry_count,
    next_delivery,
    fcm_message_id
FROM deliveries;


DROP TABLE deliveries;
ALTER TABLE deliveries_new RENAME TO deliveries;


CREATE INDEX "idx_deliveries_receiver" ON deliveries (message_id, receiver_client_id);



