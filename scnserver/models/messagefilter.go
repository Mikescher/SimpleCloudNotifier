package models

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"strconv"
	"strings"
	"time"
)

type MessageFilter struct {
	ConfirmedSubscriptionBy *UserID
	SearchString            *[]string
	Sender                  *[]UserID
	ChannelNameCS           *[]string // case-sensitive
	ChannelNameCI           *[]string // case-insensitive
	ChannelID               *[]ChannelID
	SenderNameCS            *[]string // case-sensitive
	SenderNameCI            *[]string // case-insensitive
	HasSenderName           *bool
	SenderIP                *[]string
	TimestampCoalesce       *time.Time
	TimestampCoalesceAfter  *time.Time
	TimestampCoalesceBefore *time.Time
	TimestampReal           *time.Time
	TimestampRealAfter      *time.Time
	TimestampRealBefore     *time.Time
	TimestampClient         *time.Time
	TimestampClientAfter    *time.Time
	TimestampClientBefore   *time.Time
	TitleCS                 *string // case-sensitive
	TitleCI                 *string // case-insensitive
	Priority                *[]int
	UserMessageID           *[]string
	OnlyDeleted             bool
	IncludeDeleted          bool
	CompatAcknowledged      *bool
	UsedKeyID               *[]KeyTokenID
}

func (f MessageFilter) SQL() (string, string, sq.PP, error) {

	joinClause := ""
	if f.ConfirmedSubscriptionBy != nil {
		joinClause += " LEFT JOIN subscriptions AS subs on messages.channel_id = subs.channel_id "
	}
	if f.SearchString != nil {
		joinClause += " JOIN messages_fts AS mfts on (mfts.rowid = messages.rowid) "
	}

	sqlClauses := make([]string, 0)

	params := sq.PP{}

	if f.OnlyDeleted {
		sqlClauses = append(sqlClauses, "(deleted=1)")
	} else if f.IncludeDeleted {
		// nothing, return all
	} else {
		sqlClauses = append(sqlClauses, "(deleted=0)") // default
	}

	if f.ConfirmedSubscriptionBy != nil {
		sqlClauses = append(sqlClauses, "(subs.subscriber_user_id = :sub_uid AND subs.confirmed = 1)")
		params["sub_uid"] = *f.ConfirmedSubscriptionBy
	}

	if f.Sender != nil {
		filter := make([]string, 0)
		for i, v := range *f.Sender {
			filter = append(filter, fmt.Sprintf("(sender_user_id = :sender_%d)", i))
			params[fmt.Sprintf("sender_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelNameCI != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelNameCI {
			filter = append(filter, fmt.Sprintf("(messages.channel_internal_name = :channelnameci_%d COLLATE NOCASE)", i))
			params[fmt.Sprintf("channelnameci_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelNameCS != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelNameCS {
			filter = append(filter, fmt.Sprintf("(messages.channel_internal_name = :channelnamecs_%d COLLATE BINARY)", i))
			params[fmt.Sprintf("channelnamecs_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelID != nil {
		filter := make([]string, 0)
		for i, v := range *f.ChannelID {
			filter = append(filter, fmt.Sprintf("(messages.channel_id = :channelid_%d)", i))
			params[fmt.Sprintf("channelid_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SenderNameCI != nil {
		filter := make([]string, 0)
		for i, v := range *f.SenderNameCI {
			filter = append(filter, fmt.Sprintf("(sender_name = :sendernameci_%d COLLATE NOCASE)", i))
			params[fmt.Sprintf("sendernameci_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "(sender_name IS NOT NULL AND ("+strings.Join(filter, " OR ")+"))")
	}

	if f.SenderNameCS != nil {
		filter := make([]string, 0)
		for i, v := range *f.SenderNameCS {
			filter = append(filter, fmt.Sprintf("(sender_name = :sendernamecs_%d COLLATE BINARY)", i))
			params[fmt.Sprintf("sendernamecs_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "(sender_name IS NOT NULL AND ("+strings.Join(filter, " OR ")+"))")
	}

	if f.HasSenderName != nil {
		if *f.HasSenderName {
			sqlClauses = append(sqlClauses, "(sender_name IS NOT NULL)")
		} else {
			sqlClauses = append(sqlClauses, "(sender_name IS     NULL)")
		}
	}

	if f.SenderIP != nil {
		filter := make([]string, 0)
		for i, v := range *f.SenderIP {
			filter = append(filter, fmt.Sprintf("(sender_ip = :senderip_%d)", i))
			params[fmt.Sprintf("senderip_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.TimestampCoalesce != nil {
		sqlClauses = append(sqlClauses, "(COALESCE(timestamp_client, timestamp_real) = :ts_equals)")
		params["ts_equals"] = (*f.TimestampCoalesce).UnixMilli()
	}

	if f.TimestampCoalesceAfter != nil {
		sqlClauses = append(sqlClauses, "(COALESCE(timestamp_client, timestamp_real) > :ts_after)")
		params["ts_after"] = (*f.TimestampCoalesceAfter).UnixMilli()
	}

	if f.TimestampCoalesceBefore != nil {
		sqlClauses = append(sqlClauses, "(COALESCE(timestamp_client, timestamp_real) < :ts_before)")
		params["ts_before"] = (*f.TimestampCoalesceBefore).UnixMilli()
	}

	if f.TimestampReal != nil {
		sqlClauses = append(sqlClauses, "(timestamp_real = :ts_real_equals)")
		params["ts_real_equals"] = (*f.TimestampRealAfter).UnixMilli()
	}

	if f.TimestampRealAfter != nil {
		sqlClauses = append(sqlClauses, "(timestamp_real > :ts_real_after)")
		params["ts_real_after"] = (*f.TimestampRealAfter).UnixMilli()
	}

	if f.TimestampRealBefore != nil {
		sqlClauses = append(sqlClauses, "(timestamp_real < :ts_real_before)")
		params["ts_real_before"] = (*f.TimestampRealBefore).UnixMilli()
	}

	if f.TimestampClient != nil {
		sqlClauses = append(sqlClauses, "(timestamp_client IS NOT NULL AND timestamp_client = :ts_client_equals)")
		params["ts_client_equals"] = (*f.TimestampClient).UnixMilli()
	}

	if f.TimestampClientAfter != nil {
		sqlClauses = append(sqlClauses, "(timestamp_client IS NOT NULL AND timestamp_client > :ts_client_after)")
		params["ts_client_after"] = (*f.TimestampClientAfter).UnixMilli()
	}

	if f.TimestampClientBefore != nil {
		sqlClauses = append(sqlClauses, "(timestamp_client IS NOT NULL AND timestamp_client < :ts_client_before)")
		params["ts_client_before"] = (*f.TimestampClientBefore).UnixMilli()
	}

	if f.TitleCI != nil {
		sqlClauses = append(sqlClauses, "(title = :titleci COLLATE NOCASE)")
		params["titleci"] = *f.TitleCI
	}

	if f.TitleCS != nil {
		sqlClauses = append(sqlClauses, "(title = :titleci COLLATE BINARY)")
		params["titleci"] = *f.TitleCI
	}

	if f.Priority != nil {
		prioList := "(" + strings.Join(langext.ArrMap(*f.Priority, func(p int) string { return strconv.Itoa(p) }), ", ") + ")"
		sqlClauses = append(sqlClauses, "(priority IN "+prioList+")")
	}

	if f.UserMessageID != nil {
		filter := make([]string, 0)
		for i, v := range *f.UserMessageID {
			filter = append(filter, fmt.Sprintf("(usr_message_id = :usermessageid_%d)", i))
			params[fmt.Sprintf("usermessageid_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "(usr_message_id IS NOT NULL AND ("+strings.Join(filter, " OR ")+"))")
	}

	if f.CompatAcknowledged != nil {
		joinClause += " LEFT JOIN compat_acks AS filter_compatack_compat_acks on messages.message_id = filter_compatack_compat_acks.message_id "

		if *f.CompatAcknowledged {
			sqlClauses = append(sqlClauses, "(filter_compatack_compat_acks.message_id IS NOT NULL)")
		} else {
			sqlClauses = append(sqlClauses, "(filter_compatack_compat_acks.message_id IS     NULL)")
		}
	}

	if f.UsedKeyID != nil {
		filter := make([]string, 0)
		for i, v := range *f.UsedKeyID {
			filter = append(filter, fmt.Sprintf("(used_key_id = :usedkeyid_%d)", i))
			params[fmt.Sprintf("usedkeyid_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SearchString != nil {
		filter := make([]string, 0)
		for i, v := range *f.SearchString {
			filter = append(filter, fmt.Sprintf("(messages_fts match :searchstring_%d)", i))
			params[fmt.Sprintf("searchstring_%d", i)] = v
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	sqlClause := ""
	if len(sqlClauses) > 0 {
		sqlClause = strings.Join(sqlClauses, " AND ")
	} else {
		sqlClause = "1=1"
	}

	return sqlClause, joinClause, params, nil
}

func (f MessageFilter) Hash() string {
	bh, err := dataext.StructHash(f, dataext.StructHashOptions{HashAlgo: sha512.New()})
	if err != nil {
		return "00000000"
	}

	str := hex.EncodeToString(bh)
	return str[0:mathext.Min(8, len(bh))]
}
