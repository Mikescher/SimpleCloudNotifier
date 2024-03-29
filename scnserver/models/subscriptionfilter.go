package models

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strings"
	"time"
)

type SubscriptionFilter struct {
	AnyUserID                *UserID
	SubscriberUserID         *[]UserID
	SubscriberUserID2        *[]UserID // Used to filter <SubscriberUserID> again
	ChannelOwnerUserID       *[]UserID
	ChannelOwnerUserID2      *[]UserID // Used to filter <ChannelOwnerUserID> again
	ChannelID                *[]ChannelID
	Confirmed                *bool
	SubscriberIsChannelOwner *bool
	Timestamp                *time.Time
	TimestampAfter           *time.Time
	TimestampBefore          *time.Time
}

func (f SubscriptionFilter) SQL() (string, string, sq.PP, error) {

	joinClause := ""

	sqlClauses := make([]string, 0)

	params := sq.PP{}

	if f.AnyUserID != nil {
		sqlClauses = append(sqlClauses, "(subscriber_user_id = :anyuid1 OR channel_owner_user_id = :anyuid2)")
		params["anyuid1"] = *f.AnyUserID
		params["anyuid2"] = *f.AnyUserID
	}

	if f.SubscriberUserID != nil {
		filter := make([]string, 0)
		for i, v := range *f.SubscriberUserID {
			filter = append(filter, fmt.Sprintf("(subscriber_user_id = :subscriber_uid_1_%d)", i))
			params[fmt.Sprintf("subscriber_uid_1_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SubscriberUserID2 != nil {
		filter := make([]string, 0)
		for i, v := range *f.SubscriberUserID2 {
			filter = append(filter, fmt.Sprintf("(subscriber_user_id = :subscriber_uid_2_%d)", i))
			params[fmt.Sprintf("subscriber_uid_2_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelOwnerUserID != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelOwnerUserID {
			filter = append(filter, fmt.Sprintf("(channel_owner_user_id = :chanowner_uid_1_%d)", i))
			params[fmt.Sprintf("chanowner_uid_1_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelOwnerUserID2 != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelOwnerUserID2 {
			filter = append(filter, fmt.Sprintf("(channel_owner_user_id = :chanowner_uid_2_%d)", i))
			params[fmt.Sprintf("chanowner_uid_2_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelID != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelID {
			filter = append(filter, fmt.Sprintf("(channel_id = :chanid_%d)", i))
			params[fmt.Sprintf("chanid_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.Confirmed != nil {
		if *f.Confirmed {
			sqlClauses = append(sqlClauses, "(confirmed=1)")
		} else {
			sqlClauses = append(sqlClauses, "(confirmed=0)")
		}
	}

	if f.SubscriberIsChannelOwner != nil {
		if *f.SubscriberIsChannelOwner {
			sqlClauses = append(sqlClauses, "(subscriber_user_id  = channel_owner_user_id)")
		} else {
			sqlClauses = append(sqlClauses, "(subscriber_user_id != channel_owner_user_id)")
		}
	}

	if f.Timestamp != nil {
		sqlClauses = append(sqlClauses, "(timestamp_created = :ts_equals)")
		params["ts_equals"] = (*f.Timestamp).UnixMilli()
	}

	if f.TimestampAfter != nil {
		sqlClauses = append(sqlClauses, "(timestamp_created > :ts_after)")
		params["ts_after"] = (*f.TimestampAfter).UnixMilli()
	}

	if f.TimestampBefore != nil {
		sqlClauses = append(sqlClauses, "(timestamp_created < :ts_before)")
		params["ts_before"] = (*f.TimestampBefore).UnixMilli()
	}

	sqlClause := ""
	if len(sqlClauses) > 0 {
		sqlClause = strings.Join(sqlClauses, " AND ")
	} else {
		sqlClause = "1=1"
	}

	return sqlClause, joinClause, params, nil
}

func (f SubscriptionFilter) Hash() string {
	bh, err := dataext.StructHash(f, dataext.StructHashOptions{HashAlgo: sha512.New()})
	if err != nil {
		return "00000000"
	}

	str := hex.EncodeToString(bh)
	return str[0:mathext.Min(8, len(bh))]
}
