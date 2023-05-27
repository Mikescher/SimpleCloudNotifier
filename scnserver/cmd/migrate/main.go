package main

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"blackforestbytes.com/simplecloudnotifier/logic"
	"blackforestbytes.com/simplecloudnotifier/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"gogs.mikescher.com/BlackForestBytes/goext/sq"
	"os"
	"regexp"
	"strings"
	"time"
)

type OldUser struct {
	UserId            int64      `db:"user_id"`
	UserKey           string     `db:"user_key"`
	FcmToken          *string    `db:"fcm_token"`
	MessagesSent      int64      `db:"messages_sent"`
	TimestampCreated  time.Time  `db:"timestamp_created"`
	TimestampAccessed *time.Time `db:"timestamp_accessed"`
	QuotaToday        int64      `db:"quota_today"`
	QuotaDay          *time.Time `db:"quota_day"`
	IsPro             bool       `db:"is_pro"`
	ProToken          *string    `db:"pro_token"`
}

type OldMessage struct {
	ScnMessageId  int64     `db:"scn_message_id"`
	SenderUserId  int64     `db:"sender_user_id"`
	TimestampReal time.Time `db:"timestamp_real"`
	Ack           []uint8   `db:"ack"`
	Title         string    `db:"title"`
	Content       *string   `db:"content"`
	Priority      int64     `db:"priority"`
	Sendtime      int64     `db:"sendtime"`
	FcmMessageId  *string   `db:"fcm_message_id"`
	UsrMessageId  *string   `db:"usr_message_id"`
}

type SCNExport struct {
	Messages []SCNExportMessage `json:"cmessagelist"`
}

type SCNExportMessage struct {
	MessageID int64 `json:"scnid"`
}

func main() {
	ctx := context.Background()

	conf, _ := scn.GetConfig("local-host")
	conf.DBMain.File = ".run-data/migrate_main.sqlite3"
	conf.DBMain.EnableLogger = false

	if _, err := os.Stat(".run-data/migrate_main.sqlite3"); err == nil {
		err = os.Remove(".run-data/migrate_main.sqlite3")
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(".run-data/migrate_main.sqlite3-shm"); err == nil {
		err = os.Remove(".run-data/migrate_main.sqlite3-shm")
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(".run-data/migrate_main.sqlite3-wal"); err == nil {
		err = os.Remove(".run-data/migrate_main.sqlite3-wal")
		if err != nil {
			panic(err)
		}
	}

	sqlite, err := logic.NewDBPool(conf)
	if err != nil {
		panic(err)
	}

	err = sqlite.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	connstr := os.Getenv("SQL_CONN_STR")
	if connstr == "" {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Enter DB URL [127.0.0.1:3306]: ")
		scanner.Scan()
		host := scanner.Text()
		if host == "" {
			host = "127.0.0.1:3306"
		}

		fmt.Print("Enter DB Username [root]: ")
		scanner.Scan()
		username := scanner.Text()
		if host == "" {
			host = "root"
		}

		fmt.Print("Enter DB Password []: ")
		scanner.Scan()
		pass := scanner.Text()
		if host == "" {
			host = ""
		}

		connstr = fmt.Sprintf("%s:%s@tcp(%s)", username, pass, host)
	}

	_dbold, err := sqlx.Open("mysql", connstr+"/simple_cloud_notifier?parseTime=true")
	if err != nil {
		panic(err)
	}
	dbold := sq.NewDB(_dbold)

	rowsUser, err := dbold.Query(ctx, "SELECT * FROM users", sq.PP{})
	if err != nil {
		panic(err)
	}

	var export SCNExport
	exfn, err := os.ReadFile("scn_export.json")
	err = json.Unmarshal(exfn, &export)
	if err != nil {
		panic(err)
	}

	appids := make(map[int64]int64)
	for _, v := range export.Messages {
		appids[v.MessageID] = v.MessageID
	}

	users := make([]OldUser, 0)
	for rowsUser.Next() {
		var u OldUser
		err = rowsUser.StructScan(&u)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	fmt.Printf("\n")

	for _, v := range users {
		fmt.Printf("========================================\n")
		fmt.Printf("            MIGRATE USER %d\n", v.UserId)
		fmt.Printf("========================================\n")
		migrateUser(ctx, sqlite.Primary.DB(), dbold, v, appids)
		fmt.Printf("========================================\n")
		fmt.Printf("\n")
		fmt.Printf("\n")
	}

	err = sqlite.Stop(context.Background())
	if err != nil {
		panic(err)
	}
}

var rexTitleChannel = rext.W(regexp.MustCompile("^\\[(?P<channel>[A-Za-z\\-0-9_ ]+)] (?P<title>(.|\\r|\\n)+)$"))

var usedFCM = make(map[string]models.ClientID)

func migrateUser(ctx context.Context, dbnew sq.DB, dbold sq.DB, user OldUser, appids map[int64]int64) {

	rowsMessages, err := dbold.Query(ctx, "SELECT * FROM messages WHERE sender_user_id = :uid ORDER BY timestamp_real ASC", sq.PP{"uid": user.UserId})
	if err != nil {
		panic(err)
	}

	messages := make([]OldMessage, 0)
	for rowsMessages.Next() {
		var m OldMessage
		err = rowsMessages.StructScan(&m)
		if err != nil {
			panic(err)
		}
		messages = append(messages, m)
	}

	fmt.Printf("Found %d messages\n", len(messages))

	userid := models.NewUserID()

	fmt.Printf("New UserID: %s\n", userid)

	tokKeyID := models.NewKeyTokenID()
	tokKeySec := user.UserKey

	protoken := user.ProToken
	if protoken != nil {
		protoken = langext.Ptr("ANDROID|v1|" + *protoken)
	}

	_, err = dbnew.Exec(ctx, "INSERT INTO users (user_id, username, is_pro, pro_token, timestamp_created) VALUES (:uid, :un, :pro, :tok, :ts)", sq.PP{
		"uid": userid,
		"un":  nil,
		"pro": langext.Conditional(user.IsPro, 1, 0),
		"tok": protoken,
		"ts":  user.TimestampCreated.UnixMilli(),
	})
	if err != nil {
		panic(err)
	}

	_, err = dbnew.Exec(ctx, "INSERT INTO compat_ids (old, new, type) VALUES (:old, :new, :typ)", sq.PP{
		"old": user.UserId,
		"new": userid,
		"typ": "userid",
	})
	if err != nil {
		panic(err)
	}

	_, err = dbnew.Exec(ctx, "INSERT INTO keytokens (keytoken_id, name, timestamp_created, owner_user_id, all_channels, channels, token, permissions) VALUES (:tid, :nam, :tsc, :owr, :all, :cha, :tok, :prm)", sq.PP{
		"tid": tokKeyID,
		"nam": "AdminKey (migrated)",
		"tsc": user.TimestampCreated.UnixMilli(),
		"owr": userid,
		"all": 1,
		"cha": "",
		"tok": tokKeySec,
		"prm": "A",
	})
	if err != nil {
		panic(err)
	}

	var clientid *models.ClientID = nil

	if user.FcmToken != nil && *user.FcmToken != "BLACKLISTED" {

		if _, ok := usedFCM[*user.FcmToken]; ok {

			fmt.Printf("Skip Creating Client (fcm token reuse)\n")

		} else {
			_clientid := models.NewClientID()

			_, err = dbnew.Exec(ctx, "INSERT INTO clients (client_id, user_id, type, fcm_token, timestamp_created, agent_model, agent_version) VALUES (:cid, :uid, :typ, :fcm, :ts, :am, :av)", sq.PP{
				"cid": _clientid,
				"uid": userid,
				"typ": "ANDROID",
				"fcm": *user.FcmToken,
				"ts":  user.TimestampCreated.UnixMilli(),
				"am":  "[migrated]",
				"av":  "[migrated]",
			})
			if err != nil {
				panic(err)
			}

			fmt.Printf("Created Client %s\n", _clientid)

			clientid = &_clientid

			usedFCM[*user.FcmToken] = _clientid

			_, err = dbnew.Exec(ctx, "INSERT INTO compat_clients (client_id) VALUES (:cid)", sq.PP{"cid": _clientid})
			if err != nil {
				panic(err)
			}
		}
	}

	mainChannelID := models.NewChannelID()
	_, err = dbnew.Exec(ctx, "INSERT INTO channels (channel_id, owner_user_id, display_name, internal_name, description_name, subscribe_key, timestamp_created) VALUES (:cid, :ouid, :dnam, :inam, :hnam, :subkey, :ts)", sq.PP{
		"cid":    mainChannelID,
		"ouid":   userid,
		"dnam":   "main",
		"inam":   "main",
		"hnam":   nil,
		"subkey": scn.RandomAuthKey(),
		"ts":     user.TimestampCreated.UnixMilli(),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created (Main) Channel [%s]: %s\n", "main", mainChannelID)

	_, err = dbnew.Exec(ctx, "INSERT INTO subscriptions (subscription_id, subscriber_user_id, channel_owner_user_id, channel_internal_name, channel_id, timestamp_created, confirmed) VALUES (:sid, :suid, :ouid, :cnam, :cid, :ts, :conf)", sq.PP{
		"sid":  models.NewSubscriptionID(),
		"suid": user.UserId,
		"ouid": user.UserId,
		"cnam": "main",
		"cid":  mainChannelID,
		"ts":   user.TimestampCreated.UnixMilli(),
		"conf": true,
	})
	if err != nil {
		panic(err)
	}

	channelMap := make(map[string]models.ChannelID)

	lastTitle := ""
	lastChannel := models.NewChannelID()
	lastContent := langext.Ptr("")
	lastSendername := langext.Ptr("")
	lastTimestamp := time.Time{}

	for _, oldmessage := range messages {

		messageid := models.NewMessageID()

		title := oldmessage.Title

		channelInternalName := "main"
		channelID := mainChannelID

		if oldmessage.UsrMessageId != nil && strings.TrimSpace(*oldmessage.UsrMessageId) == "" {
			oldmessage.UsrMessageId = nil
		}

		if match, ok := rexTitleChannel.MatchFirst(title); ok {

			chanNameTitle := match.GroupByName("channel").Value()

			if strings.HasPrefix(chanNameTitle, "VBOARD ERROR") {
				chanNameTitle = "VBOARD-ERROR"
			}

			if chanNameTitle != "status" {
				title = match.GroupByName("title").Value()

				dummyApp := logic.Application{}

				dispName := dummyApp.NormalizeChannelDisplayName(chanNameTitle)
				intName := dummyApp.NormalizeChannelInternalName(chanNameTitle)

				if v, ok := channelMap[intName]; ok {
					channelID = v
					channelInternalName = intName
				} else {

					channelID = models.NewChannelID()
					channelInternalName = intName

					_, err = dbnew.Exec(ctx, "INSERT INTO channels (channel_id, owner_user_id, display_name, internal_name, description_name, subscribe_key, timestamp_created) VALUES (:cid, :ouid, :dnam, :inam, :hnam, :subkey, :ts)", sq.PP{
						"cid":    channelID,
						"ouid":   userid,
						"dnam":   dispName,
						"inam":   intName,
						"hnam":   nil,
						"subkey": scn.RandomAuthKey(),
						"ts":     oldmessage.TimestampReal.UnixMilli(),
					})
					if err != nil {
						panic(err)
					}

					_, err = dbnew.Exec(ctx, "INSERT INTO subscriptions (subscription_id, subscriber_user_id, channel_owner_user_id, channel_internal_name, channel_id, timestamp_created, confirmed) VALUES (:sid, :suid, :ouid, :cnam, :cid, :ts, :conf)", sq.PP{
						"sid":  models.NewSubscriptionID(),
						"suid": user.UserId,
						"ouid": user.UserId,
						"cnam": intName,
						"cid":  channelID,
						"ts":   oldmessage.TimestampReal.UnixMilli(),
						"conf": true,
					})
					if err != nil {
						panic(err)
					}

					channelMap[intName] = channelID

					fmt.Printf("Auto Created Channel [%s]: %s\n", dispName, channelID)

				}
			}
		}

		sendername := determineSenderName(user, oldmessage, title, oldmessage.Content, channelInternalName)

		if lastTitle == title && channelID == lastChannel &&
			langext.PtrEquals(lastContent, oldmessage.Content) &&
			langext.PtrEquals(lastSendername, sendername) && oldmessage.TimestampReal.Sub(lastTimestamp) < 5*time.Second {

			lastTitle = title
			lastChannel = channelID
			lastContent = oldmessage.Content
			lastSendername = sendername
			lastTimestamp = oldmessage.TimestampReal

			fmt.Printf("Skip message [%d] \"%s\" (fast-duplicate)\n", oldmessage.ScnMessageId, oldmessage.Title)

			continue
		}

		var sendTimeMillis *int64 = nil
		if oldmessage.Sendtime > 0 && (oldmessage.Sendtime*1000) != oldmessage.TimestampReal.UnixMilli() {
			sendTimeMillis = langext.Ptr(oldmessage.Sendtime * 1000)
		}

		if user.UserId == 56 && oldmessage.ScnMessageId >= 15729 {
			if _, ok := appids[oldmessage.ScnMessageId]; !ok {

				lastTitle = title
				lastChannel = channelID
				lastContent = oldmessage.Content
				lastSendername = sendername
				lastTimestamp = oldmessage.TimestampReal

				fmt.Printf("Skip message [%d] \"%s\" (locally deleted in app)\n", oldmessage.ScnMessageId, oldmessage.Title)
				continue
			}
		}

		pp := sq.PP{
			"mid":  messageid,
			"suid": userid,
			"ouid": user.UserId,
			"cnam": channelInternalName,
			"cid":  channelID,
			"tsr":  oldmessage.TimestampReal.UnixMilli(),
			"tsc":  sendTimeMillis,
			"tit":  title,
			"cnt":  oldmessage.Content,
			"prio": oldmessage.Priority,
			"umid": oldmessage.UsrMessageId,
			"ip":   "",
			"snam": sendername,
			"ukid": tokKeyID,
		}
		_, err = dbnew.Exec(ctx, "INSERT INTO messages (message_id, sender_user_id, owner_user_id, channel_internal_name, channel_id, timestamp_real, timestamp_client, title, content, priority, usr_message_id, sender_ip, sender_name, used_key_id) VALUES (:mid, :suid, :ouid, :cnam, :cid, :tsr, :tsc, :tit, :cnt, :prio, :umid, :ip, :snam, :ukid)", pp)
		if err != nil {
			jv, _ := json.MarshalIndent(pp, "", "  ")
			fmt.Printf("%s", string(jv))
			panic(err)
		}

		_, err = dbnew.Exec(ctx, "INSERT INTO compat_ids (old, new, type) VALUES (:old, :new, :typ)", sq.PP{
			"old": oldmessage.ScnMessageId,
			"new": messageid,
			"typ": "messageid",
		})
		if err != nil {
			panic(err)
		}

		if len(oldmessage.Ack) == 1 && oldmessage.Ack[0] == 1 {

			if clientid != nil {
				_, err = dbnew.Exec(ctx, "INSERT INTO deliveries (delivery_id, message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (:did, :mid, :ruid, :rcid, :tsc, :tsf, :stat, :fcm, :next)", sq.PP{
					"did":  models.NewDeliveryID(),
					"mid":  messageid,
					"ruid": user.UserId,
					"rcid": *clientid,
					"tsc":  oldmessage.TimestampReal.UnixMilli(),
					"tsf":  oldmessage.TimestampReal.UnixMilli(),
					"stat": models.DeliveryStatusSuccess,
					"fcm":  *user.FcmToken,
					"next": nil,
				})
				if err != nil {
					panic(err)
				}
			}

			_, err = dbnew.Exec(ctx, "INSERT INTO compat_acks (user_id, message_id) VALUES (:uid, :mid)", sq.PP{
				"uid": userid,
				"mid": messageid,
			})
			if err != nil {
				panic(err)
			}

		} else if len(oldmessage.Ack) == 1 && oldmessage.Ack[0] == 0 {

			if clientid != nil {
				_, err = dbnew.Exec(ctx, "INSERT INTO deliveries (delivery_id, message_id, receiver_user_id, receiver_client_id, timestamp_created, timestamp_finalized, status, fcm_message_id, next_delivery) VALUES (:did, :mid, :ruid, :rcid, :tsc, :tsf, :stat, :fcm, :next)", sq.PP{
					"did":  models.NewDeliveryID(),
					"mid":  messageid,
					"ruid": user.UserId,
					"rcid": *clientid,
					"tsc":  oldmessage.TimestampReal.UnixMilli(),
					"tsf":  oldmessage.TimestampReal.UnixMilli(),
					"stat": models.DeliveryStatusFailed,
					"fcm":  *user.FcmToken,
					"next": nil,
				})
				if err != nil {
					panic(err)
				}

				fmt.Printf("Create failed-delivery for message %d (no ack)\n", oldmessage.ScnMessageId)
			}

		} else {
			panic("cannot parse ack")
		}

		lastTitle = title
		lastChannel = channelID
		lastContent = oldmessage.Content
		lastSendername = sendername
		lastTimestamp = oldmessage.TimestampReal
	}

}

func determineSenderName(user OldUser, oldmessage OldMessage, title string, content *string, channame string) *string {
	if user.UserId != 56 {
		return nil
	}

	if channame == "t-ctrl" {
		return langext.Ptr("sbox")
	}

	if channame == "torr" {
		return langext.Ptr("sbox")
	}

	if channame == "yt-dl" {
		return langext.Ptr("mscom")
	}

	if channame == "ncc-upload" {
		return langext.Ptr("mscom")
	}

	if channame == "cron" {
		if strings.Contains(title, "error on bfb") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "error on mscom") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "error on niflheim-3") {
			return langext.Ptr("niflheim-3")
		}

		if strings.Contains(*content, "on mscom") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "on bfb") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "gogitmirror_cron") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "comic_downloader") {
			return langext.Ptr("mscom")
		}
	}

	if channame == "sshguard" {
		if strings.Contains(*content, "logged in to mscom") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "logged in to bfb") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "logged in to statussrv") {
			return langext.Ptr("statussrv")
		}
	}

	if channame == "docker-watch" {
		if strings.Contains(title, "on plantafelstaging") {
			return langext.Ptr("plantafelstaging")
		}
		if strings.Contains(title, "@ mscom") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "@ bfb") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/scn_server:latest") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "archivebox/archivebox:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "antoniomika/sish:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "binwiederhier/ntfy:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/kgserver:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/mikescher/kgserver:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "jenkins/jenkins:lts") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "mikescher/youtube-dl-viewer:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "etherpad/etherpad:latest") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "teamcity_agent") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "teamcity_server") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/inoshop/") {
			return langext.Ptr("inoshop")
		}
		if strings.Contains(*content, "inopart_mongo_") {
			return langext.Ptr("inoshop")
		}
		if strings.Contains(*content, "Image: wkk_") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/holz100") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/bewirto") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/bfb-website") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/bfb/website") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/psycho/backend") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/vereinsboard") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/isiproject") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/ar-app-supportchat-server") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/planitec/ar-app-supportchat-server") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "docker_registry") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/balu") && strings.Contains(*content, "prod") {
			return langext.Ptr("lbxprod")
		}
		if strings.Contains(*content, "registry.blackforestbytes.com/balu") && strings.Contains(*content, "dev") {
			return langext.Ptr("lbxdev")
		}
		if strings.Contains(*content, "Server: bfb-testserver") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(*content, "wptest_") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "balu-db") {
			return langext.Ptr("lbprod")
		}
	}

	if channame == "certbot" {
		if strings.Contains(title, "Update cert_badennet_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update cert_badennet_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update cert_bfbugs_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update bfbugs_0001") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update inoshop_bfb") {
			return langext.Ptr("inoshop")
		}
		if strings.Contains(title, "Update cert_bfb_0001") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "Update cert_bugkultur_0001") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "Update cert_public_0001") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "Update cert_korbers_0001") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "Update cert_wkk_staging_external") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(title, "Update cert_wkk_production_external") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(title, "Update cert_wkk_develop_external") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(title, "Update cert_wkk_internal") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(title, "Update bfb_de_wildcard") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Update cannonconquest") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Update isiproject_wildcard") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Update vereinsboard_demo") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Update vereinsboard_wildcard") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Update cert_bewirto_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update cert_badennet_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update cert_mampfkultur_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(title, "Update cert_psycho_main") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(*content, "DNS:*.blackforestbytes.com") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "DNS:*.mikescher.com") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "plantafel-digital.de") {
			return langext.Ptr("plan-web-prod")
		}
		if strings.Contains(title, "plantafeldev.de") {
			return langext.Ptr("plantafeldev")
		}
		if strings.Contains(title, "plantafelstaging.de") {
			return langext.Ptr("plantafeldev")
		}
		if strings.Contains(*content, "DNS:*.plantafeldev.de") {
			return langext.Ptr("plantafeldev")
		}
		if strings.Contains(*content, "plantafel-digital.de") {
			return langext.Ptr("plan-web-prod")
		}
		if strings.Contains(*content, "plantafeldev.de") {
			return langext.Ptr("plantafeldev")
		}
		if strings.Contains(*content, "plantafelstaging.de") {
			return langext.Ptr("plantafeldev")
		}
	}

	if channame == "space-warning" {
		if title == "bfb" {
			return langext.Ptr("bfb")
		}
		if title == "mscom" {
			return langext.Ptr("mscom")
		}
		if title == "plan-web-prod" {
			return langext.Ptr("plan-web-prod")
		}
		if title == "statussrv" {
			return langext.Ptr("statussrv")
		}
	}

	if channame == "srv-backup" {
		if strings.Contains(*content, "Server: bfb-testserver") {
			return langext.Ptr("bfb-testserver")
		}
		if strings.Contains(*content, "Server: bfb") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(*content, "Server: mscom") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(*content, "Server: statussrv") {
			return langext.Ptr("statussrv")
		}
	}

	if title == "[status] Updating uptime-kuma image" {
		return langext.Ptr("statussrv")
	}

	if channame == "omv-backup" {
		return langext.Ptr("omv")
	}

	if channame == "omv-rcheck" {
		return langext.Ptr("omv")
	}

	if channame == "tfin" {
		return langext.Ptr("sbox")
	}

	if channame == "vboard-error" {
		return langext.Ptr("bfb")
	}

	if channame == "vboard" {
		return langext.Ptr("bfb")
	}

	if channame == "cubox" {
		return langext.Ptr("cubox")
	}

	if channame == "sys" {
		if strings.Contains(title, "h2896063") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "h2516246") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "h2770024") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Reboot plan-web-prod") {
			return langext.Ptr("plan-web-prod")
		}
		if strings.Contains(title, "Reboot mikescher.com") {
			return langext.Ptr("mscom")
		}
		if strings.Contains(title, "Reboot blackforestbytes.com") {
			return langext.Ptr("bfb")
		}
		if strings.Contains(title, "Reboot plan-web-dev") {
			return langext.Ptr("plan-web-dev")
		}
		if strings.Contains(title, "Reboot plan-web-staging") {
			return langext.Ptr("plan-web-staging")
		}
		if strings.Contains(title, "Reboot virmach-01") {
			return langext.Ptr("statussrv")
		}
		if strings.Contains(title, "Reboot wkk-1") {
			return langext.Ptr("wkk")
		}
		if strings.Contains(title, "Reboot lbxprod") {
			return langext.Ptr("lbxprod")
		}
	}

	if channame == "yt-tvc" {
		return langext.Ptr("mscom")
	}

	if channame == "gdapi" {
		return langext.Ptr("bfb")
	}

	if channame == "ttrss" {
		return langext.Ptr("mscom")
	}

	if title == "NCC Upload failed" || title == "NCC Upload successful" {
		return langext.Ptr("mscom")
	}

	if oldmessage.ScnMessageId == 7975 {
		return langext.Ptr("mscom")
	}

	if strings.Contains(title, "bfbackup job") {
		return langext.Ptr("bfbackup")
	}

	if strings.Contains(title, "Repo migration of /volume1") {
		return langext.Ptr("bfbackup")
	}

	//fmt.Printf("Failed to determine sender of [%d] '%s' '%s'\n", oldmessage.ScnMessageId, oldmessage.Title, langext.Coalesce(oldmessage.Content, "<NULL>"))
	fmt.Printf("Failed to determine sender of [%d] '%s'\n", oldmessage.ScnMessageId, oldmessage.Title)

	return nil
}
