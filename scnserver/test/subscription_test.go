package test

import (
	"blackforestbytes.com/simplecloudnotifier/api/apierr"
	tt "blackforestbytes.com/simplecloudnotifier/test/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"testing"
)

func TestListSubscriptionsOfUser(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, c int, sel string) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?selector=%s", u.UID, sel))
		tt.AssertEqual(t, sel+".len", c, len(slist.Subscriptions))
	}

	assertCount(data.User[16], 3, "outgoing_all")
	assertCount(data.User[16], 3, "outgoing_confirmed")
	assertCount(data.User[16], 0, "outgoing_unconfirmed")
	assertCount(data.User[16], 3, "incoming_all")
	assertCount(data.User[16], 3, "incoming_confirmed")
	assertCount(data.User[16], 0, "incoming_unconfirmed")

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	sub1 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})
	sub2 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan2.SubscribeKey), gin.H{
		"channel_id": chan2.ChannelId,
	})
	sub3 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan3.SubscribeKey), gin.H{
		"channel_id": chan3.ChannelId,
	})

	tt.AssertNotDefaultAny(t, "sub1", sub1)
	tt.AssertNotDefaultAny(t, "sub2", sub2)
	tt.AssertNotDefaultAny(t, "sub3", sub3)

	tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1["subscription_id"]), gin.H{
		"confirmed": true,
	})

	tt.RequestAuthDelete[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub3["subscription_id"]), gin.H{})

	assertCount(data.User[16], 3, "outgoing_all")
	assertCount(data.User[16], 3, "outgoing_confirmed")
	assertCount(data.User[16], 0, "outgoing_unconfirmed")
	assertCount(data.User[16], 5, "incoming_all")
	assertCount(data.User[16], 4, "incoming_confirmed")
	assertCount(data.User[16], 1, "incoming_unconfirmed")

	tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1["subscription_id"]), gin.H{
		"confirmed": false,
	})

	assertCount(data.User[16], 5, "incoming_all")
	assertCount(data.User[16], 3, "incoming_confirmed")
	assertCount(data.User[16], 2, "incoming_unconfirmed")

	assertCount(data.User[0], 7, "outgoing_all")
	assertCount(data.User[0], 5, "outgoing_confirmed")
	assertCount(data.User[0], 2, "outgoing_unconfirmed")
	assertCount(data.User[0], 5, "incoming_all")
	assertCount(data.User[0], 5, "incoming_confirmed")
	assertCount(data.User[0], 0, "incoming_unconfirmed")
}

func TestListSubscriptionsOfChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, cid *chanobj, c int) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", u.UID, cid.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", c, len(slist.Subscriptions))
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	sub1 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	assertCount(data.User[16], chan1, 2)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	sub2 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan2.SubscribeKey), gin.H{
		"channel_id": chan2.ChannelId,
	})

	assertCount(data.User[16], chan1, 2)
	assertCount(data.User[16], chan2, 2)
	assertCount(data.User[16], chan3, 1)

	sub3 := tt.RequestAuthPost[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan3.SubscribeKey), gin.H{
		"channel_id": chan3.ChannelId,
	})

	assertCount(data.User[16], chan1, 2)
	assertCount(data.User[16], chan2, 2)
	assertCount(data.User[16], chan3, 2)

	tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1["subscription_id"]), gin.H{
		"confirmed": true,
	})

	assertCount(data.User[16], chan1, 2)
	assertCount(data.User[16], chan2, 2)
	assertCount(data.User[16], chan3, 2)

	tt.RequestAuthDelete[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub3["subscription_id"]), gin.H{})

	assertCount(data.User[16], chan1, 2)
	assertCount(data.User[16], chan2, 2)
	assertCount(data.User[16], chan3, 1)

	tt.RequestAuthDelete[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1["subscription_id"]), gin.H{})

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 2)
	assertCount(data.User[16], chan3, 1)

	tt.RequestAuthDelete[gin.H](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, sub2["subscription_id"]), gin.H{})

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)
}

func TestCreateSubscriptionToOwnChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, cid *chanobj, c int) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", u.UID, cid.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", c, len(slist.Subscriptions))
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan1.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan3.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	assertCount(data.User[16], chan1, 0)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)

	{
		sub0 := tt.RequestAuthPost[subobj](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[16].UID, chan1.SubscribeKey), gin.H{
			"channel_id": chan1.ChannelId,
		})
		tt.AssertEqual(t, "Confirmed", true, sub0.Confirmed) // sub to own channel == auto confirm
	}

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)
}

func TestCreateDoubleSubscriptionToOwnChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, cid *chanobj, c int) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", u.UID, cid.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", c, len(slist.Subscriptions))
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan1.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan3.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	assertCount(data.User[16], chan1, 0)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)

	sub0 := tt.RequestAuthPost[subobj](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions", data.User[16].UID), gin.H{
		"channel_id": chan1.ChannelId,
	})
	tt.AssertEqual(t, "Confirmed", true, sub0.Confirmed) // sub to own channel == auto confirm

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)

	sub1 := tt.RequestAuthPost[subobj](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[16].UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})
	tt.AssertEqual(t, "Confirmed", true, sub1.Confirmed)                     // sub to own channel == auto confirm
	tt.AssertEqual(t, "Confirmed", sub0.SubscriptionId, sub1.SubscriptionId) // same sub

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)
}

func TestCreateSubscriptionToForeignChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	gsub1 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId))

	gsub2 := tt.RequestAuthGet[subobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId))

	tt.AssertEqual(t, "SubscriptionId", gsub1.SubscriptionId, gsub2.SubscriptionId)

}

func TestGetSubscriptionToOwnChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data1.UID), gin.H{
		"name": "Chan1",
	})

	slist := tt.RequestAuthGet[sublist](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data1.UID, chan1.ChannelId))
	tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))

	gsub1 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, slist.Subscriptions[0].SubscriptionId))

	tt.AssertEqual(t, "Confirmed", true, gsub1.Confirmed)
	tt.AssertEqual(t, "ChannelId", chan1.ChannelId, gsub1.ChannelId)
}

func TestGetSubscriptionToForeignChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, c int, sel string) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?selector=%s", u.UID, sel))
		tt.AssertEqual(t, sel+".len", c, len(slist.Subscriptions))
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	sub1 := tt.RequestAuthPost[subobj](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})
	sub2 := tt.RequestAuthPost[subobj](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan2.SubscribeKey), gin.H{
		"channel_id": chan2.ChannelId,
	})
	sub3 := tt.RequestAuthPost[subobj](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data.User[0].UID, chan3.SubscribeKey), gin.H{
		"channel_id": chan3.ChannelId,
	})

	tt.AssertNotDefaultAny(t, "sub1", sub1)
	tt.AssertNotDefaultAny(t, "sub2", sub2)
	tt.AssertNotDefaultAny(t, "sub3", sub3)

	tt.RequestAuthPatch[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	})

	tt.RequestAuthDelete[gin.H](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, sub3.SubscriptionId), gin.H{})

	assertCount(data.User[16], 3, "outgoing_all")
	assertCount(data.User[16], 3, "outgoing_confirmed")
	assertCount(data.User[16], 0, "outgoing_unconfirmed")
	assertCount(data.User[16], 5, "incoming_all")
	assertCount(data.User[16], 4, "incoming_confirmed")
	assertCount(data.User[16], 1, "incoming_unconfirmed")

	assertCount(data.User[0], 7, "outgoing_all")
	assertCount(data.User[0], 6, "outgoing_confirmed")
	assertCount(data.User[0], 1, "outgoing_unconfirmed")
	assertCount(data.User[0], 5, "incoming_all")
	assertCount(data.User[0], 5, "incoming_confirmed")
	assertCount(data.User[0], 0, "incoming_unconfirmed")

	gsub1 := tt.RequestAuthGet[subobj](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, sub1.SubscriptionId))
	tt.AssertEqual(t, "SubscriptionId", sub1.SubscriptionId, gsub1.SubscriptionId)
	tt.AssertEqual(t, "Confirmed", true, gsub1.Confirmed)

	gsub2 := tt.RequestAuthGet[subobj](t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, sub2.SubscriptionId))
	tt.AssertEqual(t, "SubscriptionId", sub2.SubscriptionId, gsub2.SubscriptionId)
	tt.AssertEqual(t, "Confirmed", false, gsub2.Confirmed)

	tt.RequestAuthGetShouldFail(t, data.User[0].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[0].UID, sub3.SubscriptionId), 404, apierr.SUBSCRIPTION_NOT_FOUND)
}

func TestCancelSubscriptionToForeignChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	gsub1 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId))

	tt.RequestAuthDelete[tt.Void](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, gsub1.SubscriptionId), gin.H{})

	tt.RequestAuthGetShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, gsub1.SubscriptionId), 404, apierr.SUBSCRIPTION_NOT_FOUND)
}

func TestCancelSubscriptionToOwnChannel(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data := tt.InitDefaultData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}
	type chanlist struct {
		Channels []chanobj `json:"channels"`
	}

	assertCount := func(u tt.Userdat, cid *chanobj, c int) {
		slist := tt.RequestAuthGet[sublist](t, u.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", u.UID, cid.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", c, len(slist.Subscriptions))
	}

	clist := tt.RequestAuthGet[chanlist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data.User[16].UID))
	chan1 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan1" })
	chan2 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan2" })
	chan3 := langext.ArrFirstOrNil(clist.Channels, func(v chanobj) bool { return v.InternalName == "Chan3" })

	assertCount(data.User[16], chan1, 1)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan1.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	assertCount(data.User[16], chan1, 0)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 1)

	{
		slist := tt.RequestAuthGet[sublist](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data.User[16].UID, chan3.ChannelId))
		tt.RequestAuthDelete[tt.Void](t, data.User[16].AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data.User[16].UID, slist.Subscriptions[0].SubscriptionId), gin.H{})
	}

	assertCount(data.User[16], chan1, 0)
	assertCount(data.User[16], chan2, 1)
	assertCount(data.User[16], chan3, 0)
}

func TestDenySubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
	}

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthDelete[tt.Void](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), gin.H{})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
		tt.RequestAuthGetShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), 404, apierr.SUBSCRIPTION_NOT_FOUND)
	}
}

func TestConfirmSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
	}

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatchShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	}, 401, apierr.USER_AUTH_FAILED)

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatch[gin.H](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
		for _, v := range slist.Subscriptions {
			tt.AssertEqual(t, "Confirmed", true, v.Confirmed)
		}
		{
			subAfter1 := tt.RequestAuthGet[subobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", true, subAfter1.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter1.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter1.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter1.ChannelOwnerUserId)
		}
		{
			subAfter2 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", true, subAfter2.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter2.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter2.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter2.ChannelOwnerUserId)
		}
	}
}

func TestUnconfirmSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
	}

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatchShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	}, 401, apierr.USER_AUTH_FAILED)

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatch[gin.H](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
		for _, v := range slist.Subscriptions {
			tt.AssertEqual(t, "Confirmed", true, v.Confirmed)
		}
		{
			subAfter1 := tt.RequestAuthGet[subobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", true, subAfter1.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter1.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter1.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter1.ChannelOwnerUserId)
		}
		{
			subAfter2 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", true, subAfter2.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter2.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter2.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter2.ChannelOwnerUserId)
		}
	}

	tt.RequestAuthPatch[gin.H](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{
		"confirmed": false,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
		for _, v := range slist.Subscriptions {
			tt.AssertEqual(t, "Confirmed", v.SubscriberUserId == data2.UID, v.Confirmed)
		}
		{
			subAfter1 := tt.RequestAuthGet[subobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", false, subAfter1.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter1.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter1.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter1.ChannelOwnerUserId)
		}
		{
			subAfter2 := tt.RequestAuthGet[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId))
			tt.AssertEqual(t, "Confirmed", false, subAfter2.Confirmed)
			tt.AssertEqual(t, "ChannelId", chan1.ChannelId, subAfter2.ChannelId)
			tt.AssertEqual(t, "SubscriberUserId", data1.UID, subAfter2.SubscriberUserId)
			tt.AssertEqual(t, "ChannelOwnerUserId", data2.UID, subAfter2.ChannelOwnerUserId)
		}
	}
}

func TestCancelIncomingSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
	}

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatch[gin.H](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthDelete[tt.Void](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
		tt.RequestAuthGetShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), 404, apierr.SUBSCRIPTION_NOT_FOUND)
	}
}

func TestCancelOutgoingSubscription(t *testing.T) {
	ws, baseUrl, stop := tt.StartSimpleWebserver(t)
	defer stop()

	data1 := tt.InitSingleData(t, ws)
	data2 := tt.InitSingleData(t, ws)

	type subobj struct {
		ChannelId           string `json:"channel_id"`
		ChannelInternalName string `json:"channel_internal_name"`
		ChannelOwnerUserId  string `json:"channel_owner_user_id"`
		Confirmed           bool   `json:"confirmed"`
		SubscriberUserId    string `json:"subscriber_user_id"`
		SubscriptionId      string `json:"subscription_id"`
		TimestampCreated    string `json:"timestamp_created"`
	}
	type sublist struct {
		Subscriptions []subobj `json:"subscriptions"`
	}
	type chanobj struct {
		ChannelId       string `json:"channel_id"`
		DescriptionName string `json:"description_name"`
		DisplayName     string `json:"display_name"`
		InternalName    string `json:"internal_name"`
		MessagesSent    int    `json:"messages_sent"`
		OwnerUserId     string `json:"owner_user_id"`
		SubscribeKey    string `json:"subscribe_key"`
		Subscription    struct {
			ChannelId           string `json:"channel_id"`
			ChannelInternalName string `json:"channel_internal_name"`
			ChannelOwnerUserId  string `json:"channel_owner_user_id"`
			Confirmed           bool   `json:"confirmed"`
			SubscriberUserId    string `json:"subscriber_user_id"`
			SubscriptionId      string `json:"subscription_id"`
			TimestampCreated    string `json:"timestamp_created"`
		} `json:"subscription"`
		TimestampCreated  string `json:"timestamp_created"`
		TimestampLastsent string `json:"timestamp_lastsent"`
	}

	chan1 := tt.RequestAuthPost[chanobj](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels", data2.UID), gin.H{
		"name": "Chan1",
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
	}

	sub1 := tt.RequestAuthPost[subobj](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions?chan_subscribe_key=%s", data1.UID, chan1.SubscribeKey), gin.H{
		"channel_id": chan1.ChannelId,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthPatch[gin.H](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data2.UID, sub1.SubscriptionId), gin.H{
		"confirmed": true,
	})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 2, len(slist.Subscriptions))
	}

	tt.RequestAuthDelete[tt.Void](t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), gin.H{})

	{
		slist := tt.RequestAuthGet[sublist](t, data2.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/channels/%s/subscriptions", data2.UID, chan1.ChannelId))
		tt.AssertEqual(t, "channel.subs.len", 1, len(slist.Subscriptions))
		tt.RequestAuthGetShouldFail(t, data1.AdminKey, baseUrl, fmt.Sprintf("/api/v2/users/%s/subscriptions/%s", data1.UID, sub1.SubscriptionId), 404, apierr.SUBSCRIPTION_NOT_FOUND)
	}
}
