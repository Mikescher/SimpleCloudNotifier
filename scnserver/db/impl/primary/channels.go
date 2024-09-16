package primary

import (
	"blackforestbytes.com/simplecloudnotifier/db"
	"blackforestbytes.com/simplecloudnotifier/models"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"time"
)

func (db *Database) GetChannelByName(ctx db.TxContext, userid models.UserID, chanName string) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QuerySingleOpt[models.Channel](ctx, tx, "SELECT * FROM channels WHERE owner_user_id = :uid AND internal_name = :nam LIMIT 1", sq.PP{"uid": userid, "nam": chanName}, sq.SModeExtended, sq.Safe)
}

func (db *Database) GetChannelByID(ctx db.TxContext, chanid models.ChannelID) (*models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	return sq.QuerySingleOpt[models.Channel](ctx, tx, "SELECT * FROM channels WHERE channel_id = :cid LIMIT 1", sq.PP{"cid": chanid}, sq.SModeExtended, sq.Safe)
}

type CreateChanel struct {
	UserId       models.UserID
	DisplayName  string
	IntName      string
	SubscribeKey string
	Description  *string
}

func (db *Database) CreateChannel(ctx db.TxContext, userid models.UserID, dispName string, intName string, subscribeKey string, description *string) (models.Channel, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.Channel{}, err
	}

	entity := models.Channel{
		ChannelID:         models.NewChannelID(),
		OwnerUserID:       userid,
		DisplayName:       dispName,
		InternalName:      intName,
		SubscribeKey:      subscribeKey,
		DescriptionName:   description,
		TimestampCreated:  models.NowSCNTime(),
		TimestampLastSent: nil,
		MessagesSent:      0,
	}

	_, err = sq.InsertSingle(ctx, tx, "channels", entity)
	if err != nil {
		return models.Channel{}, err
	}

	return entity, nil
}

func (db *Database) ListChannelsByOwner(ctx db.TxContext, userid models.UserID, subUserID models.UserID) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	order := " ORDER BY channels.timestamp_created ASC, channels.channel_id ASC "

	sql := "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub ON channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE owner_user_id = :ouid" + order

	pp := sq.PP{
		"ouid":   userid,
		"subuid": subUserID,
	}

	return sq.QueryAll[models.ChannelWithSubscription](ctx, tx, sql, pp, sq.SModeExtended, sq.Safe)
}

func (db *Database) ListChannelsBySubscriber(ctx db.TxContext, userid models.UserID, confirmed *bool) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := ""
	if confirmed != nil && *confirmed {
		confCond = " AND sub.confirmed = 1"
	} else if confirmed != nil && !*confirmed {
		confCond = " AND sub.confirmed = 0"
	}

	order := " ORDER BY channels.timestamp_created ASC, channels.channel_id ASC "

	sql := "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE sub.subscription_id IS NOT NULL " + confCond + order

	pp := sq.PP{
		"subuid": userid,
	}

	return sq.QueryAll[models.ChannelWithSubscription](ctx, tx, sql, pp, sq.SModeExtended, sq.Safe)
}

func (db *Database) ListChannelsByAccess(ctx db.TxContext, userid models.UserID, confirmed *bool) ([]models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return nil, err
	}

	confCond := "OR (sub.subscription_id IS NOT NULL)"
	if confirmed != nil && *confirmed {
		confCond = "OR (sub.subscription_id IS NOT NULL AND sub.confirmed = 1)"
	} else if confirmed != nil && !*confirmed {
		confCond = "OR (sub.subscription_id IS NOT NULL AND sub.confirmed = 0)"
	}

	order := " ORDER BY channels.timestamp_created ASC, channels.channel_id ASC "

	sql := "SELECT channels.*, sub.* FROM channels LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid WHERE owner_user_id = :ouid " + confCond + order

	pp := sq.PP{
		"ouid":   userid,
		"subuid": userid,
	}

	return sq.QueryAll[models.ChannelWithSubscription](ctx, tx, sql, pp, sq.SModeExtended, sq.Safe)
}

func (db *Database) GetChannel(ctx db.TxContext, userid models.UserID, channelid models.ChannelID, enforceOwner bool) (models.ChannelWithSubscription, error) {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return models.ChannelWithSubscription{}, err
	}

	params := sq.PP{
		"cid":    channelid,
		"subuid": userid,
	}

	selectors := "channels.*, sub.*"

	join := "LEFT JOIN subscriptions AS sub on channels.channel_id = sub.channel_id AND sub.subscriber_user_id = :subuid"

	cond := "channels.channel_id = :cid"
	if enforceOwner {
		cond = "owner_user_id = :ouid AND channels.channel_id = :cid"
		params["ouid"] = userid
	}

	sql := "SELECT " + selectors + " FROM channels " + join + " WHERE " + cond + " LIMIT 1"

	return sq.QuerySingle[models.ChannelWithSubscription](ctx, tx, sql, params, sq.SModeExtended, sq.Safe)
}

func (db *Database) IncChannelMessageCounter(ctx db.TxContext, channel *models.Channel) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	now := time.Now()

	_, err = tx.Exec(ctx, "UPDATE channels SET messages_sent = messages_sent+1, timestamp_lastsent = :ts WHERE channel_id = :cid", sq.PP{
		"ts":  time2DB(now),
		"cid": channel.ChannelID,
	})
	if err != nil {
		return err
	}

	channel.MessagesSent += 1
	channel.TimestampLastSent = models.NewSCNTimePtr(&now)

	return nil
}

func (db *Database) UpdateChannelSubscribeKey(ctx db.TxContext, channelid models.ChannelID, newkey string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET subscribe_key = :key WHERE channel_id = :cid", sq.PP{
		"key": newkey,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateChannelDisplayName(ctx db.TxContext, channelid models.ChannelID, dispname string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET display_name = :nam WHERE channel_id = :cid", sq.PP{
		"nam": dispname,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateChannelDescriptionName(ctx db.TxContext, channelid models.ChannelID, descname *string) error {
	tx, err := ctx.GetOrCreateTransaction(db)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE channels SET description_name = :nam WHERE channel_id = :cid", sq.PP{
		"nam": descname,
		"cid": channelid,
	})
	if err != nil {
		return err
	}

	return nil
}
