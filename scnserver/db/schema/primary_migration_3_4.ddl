
DROP INDEX idx_messages_owner_channel;


DROP INDEX idx_messages_owner_channel_nc;


DROP INDEX idx_messages_idempotency;
CREATE UNIQUE INDEX "idx_messages_idempotency" ON messages (sender_user_id, usr_message_id COLLATE BINARY);


DROP INDEX idx_messages_usedkey;
CREATE INDEX "idx_messages_usedkey" ON messages (sender_user_id, used_key_id);


ALTER TABLE messages DROP COLUMN owner_user_id;




