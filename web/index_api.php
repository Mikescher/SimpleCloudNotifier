<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="/css/mini-default.min.css">
	<link rel="stylesheet" href="/css/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>


	<form id="mainpnl">

        <a href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier" class="button bordered" id="tl_link"><span class="icn-google-play"></span></a>
        <a href="/index.php" class="button bordered" id="tr_link">Send</a>

		<h1>Simple Cloud Notifier</h1>

        <p>Get your user-id and user-key from the app and send notifications to your phone by performing a POST request against <code>https://simplecloudnotifier.blackforestbytes.com/send.php</code></p>
        <pre>curl                                     \
    --data "user_id={userid}"            \
    --data "user_key={userkey}"          \
    --data "message={message_title}"     \
    --data "content={message_content}"   \
    https://simplecloudnotifier.blackforestbytes.com/send.php</pre>
    <p>The <code>content</code> parameter is optional, you can also send message with only a title</p>
        <pre>curl                                     \
    --data "user_id={userid}"            \
    --data "user_key={userkey}"          \
    --data "message={message_title}"     \
    https://simplecloudnotifier.blackforestbytes.com/send.php</pre>
    </form>

	<div id="copyinfo">
		<a href="https://www.blackforestbytes.com">&#169; blackforestbytes</a>
	</div>

</body>
</html>