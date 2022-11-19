package db

import (
	"blackforestbytes.com/simplecloudnotifier/models"
	"database/sql"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"time"
)

func (db *Database) CreateUser(ctx TxContext, readKey string, sendKey string, adminKey string, protoken *string, username *string) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO users (username, read_key, send_key, admin_key, is_pro, pro_token, timestamp_created) VALUES (?, ?, ?, ?, ?, ?, ?)",
		username,
		readKey,
		sendKey,
		adminKey,
		bool2DB(protoken != nil),
		protoken,
		time2DB(now))
	if err != nil {
		return models.User{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		UserID:            liid,
		Username:          username,
		ReadKey:           readKey,
		SendKey:           sendKey,
		AdminKey:          adminKey,
		TimestampCreated:  now,
		TimestampLastRead: nil,
		TimestampLastSent: nil,
		MessagesSent:      0,
		QuotaUsed:         0,
		QuotaUsedDay:      nil,
		IsPro:             protoken != nil,
		ProToken:          protoken,
	}, nil
}

func (db *Database) CreateClient(ctx TxContext, userid int64, ctype models.ClientType, fcmToken string, agentModel string, agentVersion string) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO clients (user_id, type, fcm_token, timestamp_created, agent_model, agent_version) VALUES (?, ?, ?, ?, ?, ?)",
		userid,
		string(ctype),
		fcmToken,
		time2DB(now),
		agentModel,
		agentVersion)
	if err != nil {
		return models.Client{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Client{}, err
	}

	return models.Client{
		ClientID:         liid,
		UserID:           userid,
		Type:             ctype,
		FCMToken:         langext.Ptr(fcmToken),
		TimestampCreated: now,
		AgentModel:       agentModel,
		AgentVersion:     agentVersion,
	}, nil
}

func (db *Database) ClearFCMTokens(ctx TxContext, fcmtoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM clients WHERE fcm_token = ?", fcmtoken)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ClearProTokens(ctx TxContext, protoken string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET is_pro=0, pro_token=NULL WHERE pro_token = ?", protoken)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByKey(ctx TxContext, key string) (*models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE admin_key = ? OR send_key = ? OR read_key = ? LIMIT 1", key, key, key)
	if err != nil {
		return nil, err
	}

	user, err := models.DecodeUser(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetChannelByKey(ctx TxContext, key string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE subscribe_key = ? OR send_key = ? LIMIT 1", key, key)
	if err != nil {
		return nil, err
	}

	channel, err := models.DecodeChannel(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (db *Database) GetUser(ctx TxContext, userid int64) (models.User, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.User{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM users WHERE user_id = ? LIMIT 1", userid)
	if err != nil {
		return models.User{}, err
	}

	user, err := models.DecodeUser(rows)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *Database) UpdateUserUsername(ctx TxContext, userid int64, username *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET username = ? WHERE user_id = ?", username, userid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUserProToken(ctx TxContext, userid int64, protoken *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET pro_token = ? AND is_pro = ? WHERE user_id = ?", protoken, bool2DB(protoken != nil), userid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) ListClients(ctx TxContext, userid int64) ([]models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM clients WHERE user_id = ?", userid)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeClients(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetClient(ctx TxContext, userid int64, clientid int64) (models.Client, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Client{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM clients WHERE user_id = ? AND client_id = ? LIMIT 1", userid, clientid)
	if err != nil {
		return models.Client{}, err
	}

	client, err := models.DecodeClient(rows)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

func (db *Database) DeleteClient(ctx TxContext, clientid int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM clients WHERE client_id = ?", clientid)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetMessageByUserMessageID(ctx TxContext, usrMsgId string) (*models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM messages WHERE usr_message_id = ? LIMIT 1", usrMsgId)
	if err != nil {
		return nil, err
	}

	msg, err := models.DecodeMessage(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (db *Database) GetChannelByName(ctx TxContext, userid int64, chanName string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ? OR name = ? LIMIT 1", userid, chanName)
	if err != nil {
		return nil, err
	}

	channel, err := models.DecodeChannel(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (db *Database) CreateChannel(ctx TxContext, userid int64, name string, subscribeKey string, sendKey string) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO channels (owner_user_id, name, subscribe_key, send_key, timestamp_created) VALUES (?, ?, ?, ?, ?)",
		userid,
		name,
		subscribeKey,
		sendKey,
		time2DB(now))
	if err != nil {
		return models.Channel{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Channel{}, err
	}

	return models.Channel{
		ChannelID:         liid,
		OwnerUserID:       userid,
		Name:              name,
		SubscribeKey:      subscribeKey,
		SendKey:           sendKey,
		TimestampCreated:  now,
		TimestampLastRead: nil,
		TimestampLastSent: nil,
		MessagesSent:      0,
	}, nil
}

func (db *Database) CreateSubscribtion(ctx TxContext, subscriberUID int64, ownerUID int64, chanName string, chanID int64, confirmed bool) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO subscriptions (subscriber_user_id, channel_owner_user_id, channel_name, channel_id, timestamp_created, confirmed) VALUES (?, ?, ?, ?, ?, ?)",
		subscriberUID,
		ownerUID,
		chanName,
		chanID,
		time2DB(now),
		confirmed)
	if err != nil {
		return models.Subscription{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Subscription{}, err
	}

	return models.Subscription{
		SubscriptionID:     liid,
		SubscriberUserID:   subscriberUID,
		ChannelOwnerUserID: ownerUID,
		ChannelID:          chanID,
		ChannelName:        chanName,
		TimestampCreated:   now,
		Confirmed:          confirmed,
	}, nil
}

func (db *Database) CreateMessage(ctx TxContext, senderUserID int64, channel models.Channel, timestampSend *time.Time, title string, content *string, priority int, userMsgId *string) (models.Message, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Message{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO messages (sender_user_id, owner_user_id, channel_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		senderUserID,
		channel.OwnerUserID,
		channel.Name,
		channel.ChannelID,
		time2DB(now),
		time2DBOpt(timestampSend),
		title,
		content,
		priority,
		userMsgId)
	if err != nil {
		return models.Message{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		SCNMessageID:    liid,
		SenderUserID:    senderUserID,
		OwnerUserID:     channel.OwnerUserID,
		ChannelName:     channel.Name,
		ChannelID:       channel.ChannelID,
		TimestampReal:   now,
		TimestampClient: timestampSend,
		Title:           title,
		Content:         content,
		Priority:        priority,
		UserMessageID:   userMsgId,
	}, nil
}

func (db *Database) ListSubscriptionsByChannel(ctx TxContext, channelID int64) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE channel_id = ?", channelID)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) ListSubscriptionsByOwner(ctx TxContext, ownerUserID int64) ([]models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE channel_owner_user_id = ?", ownerUserID)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeSubscriptions(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) CreateRetryDelivery(ctx TxContext, client models.Client, msg models.Message) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now().UTC()
	next := now.Add(5 * time.Second)

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?)",
		msg.SCNMessageID,
		client.UserID,
		client.ClientID,
		time2DB(now),
		nil,
		models.DeliveryStatusRetry,
		nil,
		time2DB(next))
	if err != nil {
		return models.Delivery{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Delivery{}, err
	}

	return models.Delivery{
		DeliveryID:         liid,
		SCNMessageID:       msg.SCNMessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   now,
		TimestampFinalized: nil,
		Status:             models.DeliveryStatusRetry,
		RetryCount:         0,
		NextDelivery:       langext.Ptr(next),
		FCMMessageID:       nil,
	}, nil
}

func (db *Database) CreateSuccessDelivery(ctx TxContext, client models.Client, msg models.Message, fcmDelivID string) (models.Delivery, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Delivery{}, err
	}

	now := time.Now().UTC()

	res, err := tx.ExecContext(ctx, "INSERT INTO deliveries (scn_message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (?, ?, ?, ?, ?, ?, ?)",
		msg.SCNMessageID,
		client.UserID,
		client.ClientID,
		time2DB(now),
		time2DB(now),
		models.DeliveryStatusSuccess,
		fcmDelivID,
		nil)
	if err != nil {
		return models.Delivery{}, err
	}

	liid, err := res.LastInsertId()
	if err != nil {
		return models.Delivery{}, err
	}

	return models.Delivery{
		DeliveryID:         liid,
		SCNMessageID:       msg.SCNMessageID,
		ReceiverUserID:     client.UserID,
		ReceiverClientID:   client.ClientID,
		TimestampCreated:   now,
		TimestampFinalized: langext.Ptr(now),
		Status:             models.DeliveryStatusSuccess,
		RetryCount:         0,
		NextDelivery:       nil,
		FCMMessageID:       langext.Ptr(fcmDelivID),
	}, nil
}

func (db *Database) ListChannels(ctx TxContext, userid int64) ([]models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ?", userid)
	if err != nil {
		return nil, err
	}

	data, err := models.DecodeChannels(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *Database) GetChannel(ctx TxContext, userid int64, channelid int64) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM channels WHERE owner_user_id = ? AND channel_id = ? LIMIT 1", userid, channelid)
	if err != nil {
		return models.Channel{}, err
	}

	client, err := models.DecodeChannel(rows)
	if err != nil {
		return models.Channel{}, err
	}

	return client, nil
}

func (db *Database) GetSubscription(ctx TxContext, subid int64) (models.Subscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Subscription{}, err
	}

	rows, err := tx.QueryContext(ctx, "SELECT * FROM subscriptions WHERE subscription_id = ? LIMIT 1", subid)
	if err != nil {
		return models.Subscription{}, err
	}

	sub, err := models.DecodeSubscription(rows)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (db *Database) DeleteSubscription(ctx TxContext, subid int64) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM subscriptions WHERE subscription_id = ?", subid)
	if err != nil {
		return err
	}

	return nil
}
