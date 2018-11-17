<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="/css/mini-default.min.css"> <!-- https://minicss.org/docs -->
    <title>Simple Cloud Notifications - API</title>
    <!--<link rel="stylesheet" href="/css/mini-nord.min.css">-->
    <!--<link rel="stylesheet" href="/css/mini-dark.min.css">-->
	<link rel="stylesheet" href="/css/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="icon" type="image/png" href="/favicon.png"/>
    <link rel="icon" type="image/png" href="/favicon.ico"/>
</head>
<body>

    <div id="copyinfo">
        <a tabindex="-1" href="https://www.blackforestbytes.com">&#169; blackforestbytes</a>
        <a tabindex="-1" href="https://www.mikescher.com">made by Mike Schw&ouml;rer</a>
    </div>

	<div id="mainpnl">
        <a tabindex="-1" href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier" class="button bordered" id="tl_link"><span class="icn-google-play"></span></a>
        <a tabindex="-1" href="/index.php" class="button bordered" id="tr_link">Send</a>

        <a tabindex="-1" href="/" class="linkcaption"><h1>Simple Cloud Notifier</h1></a>

        <p>Get your user-id and user-key from the app and send notifications to your phone by performing a POST request against <code>https://simplecloudnotifier.blackforestbytes.com/send.php</code></p>
        <pre>curl                                          \
    --data "user_id={userid}"                 \
    --data "user_key={userkey}"               \
    --data "title={message_title}"            \
    --data "content={message_body}"           \
    --data "priority={0|1|2}"                 \
    --data "msg_id={unique_message_id}"       \
    https://scn.blackforestbytes.com/send.php</pre>
    <p>The <code>content</code>, <code>priority</code> and <code>msg_id</code> parameters are optional, you can also send message with only a title and the default priority</p>
        <pre>curl                                          \
    --data "user_id={userid}"                 \
    --data "user_key={userkey}"               \
    --data "title={message_title}"            \
    https://scn.blackforestbytes.com/send.php</pre>

        <a href="/api_more.php" class="button bordered tertiary" style="float: right; min-width: 100px; text-align: center">More</a>

    </div>
</body>
</html>