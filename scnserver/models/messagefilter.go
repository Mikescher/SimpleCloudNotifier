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
	SearchStringFTS         *[]string
	SearchStringPlain       *[]string
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
	if f.SearchStringFTS != nil {
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
		sqlClauses = append(sqlClauses, fmt.Sprintf("(subs.subscriber_user_id = :%s AND subs.confirmed = 1)", params.Add(*f.ConfirmedSubscriptionBy)))
	}

	if f.Sender != nil {
		filter := make([]string, 0)
		for _, v := range *f.Sender {
			filter = append(filter, fmt.Sprintf("(sender_user_id = :%s)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelNameCI != nil {
		filter := make([]string, 0)
		for _, v := range *f.ChannelNameCI {
			filter = append(filter, fmt.Sprintf("(messages.channel_internal_name = :%s COLLATE NOCASE)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelNameCS != nil {
		filter := make([]string, 0)
		for _, v := range *f.ChannelNameCS {
			filter = append(filter, fmt.Sprintf("(messages.channel_internal_name = :%s COLLATE BINARY)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.ChannelID != nil {
		filter := make([]string, 0)
		for _, v := range *f.ChannelID {
			filter = append(filter, fmt.Sprintf("(messages.channel_id = :%s)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SenderNameCI != nil {
		filter := make([]string, 0)
		for _, v := range *f.SenderNameCI {
			filter = append(filter, fmt.Sprintf("(sender_name = :%s COLLATE NOCASE)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "(sender_name IS NOT NULL AND ("+strings.Join(filter, " OR ")+"))")
	}

	if f.SenderNameCS != nil {
		filter := make([]string, 0)
		for _, v := range *f.SenderNameCS {
			filter = append(filter, fmt.Sprintf("(sender_name = :%s COLLATE BINARY)", params.Add(v)))
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
		for _, v := range *f.SenderIP {
			filter = append(filter, fmt.Sprintf("(sender_ip = :%s)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.TimestampCoalesce != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(COALESCE(timestamp_client, timestamp_real) = :%s)", params.Add((*f.TimestampCoalesce).UnixMilli())))
	}

	if f.TimestampCoalesceAfter != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(COALESCE(timestamp_client, timestamp_real) > :%s)", params.Add((*f.TimestampCoalesceAfter).UnixMilli())))
	}

	if f.TimestampCoalesceBefore != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(COALESCE(timestamp_client, timestamp_real) < :%s)", params.Add((*f.TimestampCoalesceBefore).UnixMilli())))
	}

	if f.TimestampReal != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_real = :%s)", params.Add((*f.TimestampRealAfter).UnixMilli())))
	}

	if f.TimestampRealAfter != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_real > :%s)", params.Add((*f.TimestampRealAfter).UnixMilli())))
	}

	if f.TimestampRealBefore != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_real < :%s)", params.Add((*f.TimestampRealBefore).UnixMilli())))
	}

	if f.TimestampClient != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_client IS NOT NULL AND timestamp_client = :%s)", params.Add((*f.TimestampClient).UnixMilli())))
	}

	if f.TimestampClientAfter != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_client IS NOT NULL AND timestamp_client > :%s)", params.Add((*f.TimestampClientAfter).UnixMilli())))
	}

	if f.TimestampClientBefore != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(timestamp_client IS NOT NULL AND timestamp_client < :%s)", params.Add((*f.TimestampClientBefore).UnixMilli())))
	}

	if f.TitleCI != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(title = :%s COLLATE NOCASE)", params.Add(*f.TitleCI)))
	}

	if f.TitleCS != nil {
		sqlClauses = append(sqlClauses, fmt.Sprintf("(title = :%s COLLATE BINARY)", params.Add(*f.TitleCI)))
	}

	if f.Priority != nil {
		prioList := "(" + strings.Join(langext.ArrMap(*f.Priority, func(p int) string { return strconv.Itoa(p) }), ", ") + ")"
		sqlClauses = append(sqlClauses, "(priority IN "+prioList+")")
	}

	if f.UserMessageID != nil {
		filter := make([]string, 0)
		for _, v := range *f.UserMessageID {
			filter = append(filter, fmt.Sprintf("(usr_message_id = :%s)", params.Add(v)))
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
		for _, v := range *f.UsedKeyID {
			filter = append(filter, fmt.Sprintf("(used_key_id = :%s)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SearchStringFTS != nil {
		filter := make([]string, 0)
		for _, v := range *f.SearchStringFTS {
			filter = append(filter, fmt.Sprintf("(messages_fts match :%s)", params.Add(v)))
		}
		sqlClauses = append(sqlClauses, "("+strings.Join(filter, " OR ")+")")
	}

	if f.SearchStringPlain != nil {
		filter := make([]string, 0)
		for _, v := range *f.SearchStringPlain {
			filter = append(filter, fmt.Sprintf("instr(lower(messages.channel_internal_name), lower(:%s))", params.Add(v)))
			filter = append(filter, fmt.Sprintf("instr(lower(messages.sender_name), lower(:%s))", params.Add(v)))
			filter = append(filter, fmt.Sprintf("instr(lower(messages.title), lower(:%s))", params.Add(v)))
			filter = append(filter, fmt.Sprintf("instr(lower(messages.content), lower(:%s))", params.Add(v)))

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
