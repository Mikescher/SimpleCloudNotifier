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

        <h2>Introduction</h2>
        <div class="section">
            <p>
                With this API you can send push notifications to your phone.
            </p>
            <p>
                To recieve them you will need to install the <a href="https://play.google.com/store/apps/details?id=com.blackforestbytes.simplecloudnotifier">SimpleCloudNotifier</a> app from the play store.
                When you open the app you can click on the account tab to see you unique <code>user_id</code> and <code>user_key</code>.
                These two values are used to identify and authenticate your device so that send messages can be routed to your phone.
            </p>
            <p>
                You can at any time generate a new <code>user_key</code> in the app and invalidate the old one.
            </p>
            <p>
                There is also a <a href="/index.php">web interface</a> for this API to manually send notifications to your phone or to test your setup.
            </p>
        </div>

        <h2>Quota</h2>
        <div class="section">
            <p>
                By default you can send up to 100 messages per day per device.
                If you need more you can upgrade your account in the app to get 1000 messages per day, this has the additional benefit of removing ads and supporting the development of the app (and making sure I can pay the server costs).
            </p>
        </div>

        <h2>API Requests</h2>
        <div class="section">
            <p>
                To send a new notification you send a <code>POST</code> request to the URL <code>https://scn.blackforestbytes.com/send.php</code>.
                All Parameters can either directly be submitted as URL parameters or they can be put into the POST body.
            </p>
            <p>
                You <i>need</i> to supply a valid <code>user_id</code> - <code>user_key</code> pair and a <code>title</code> for your message, all other parameter are optional.
            </p>
        </div>

        <h2>API Response</h2>
        <div class="section">
            <p>
                If the operation was successful the API will respond with an HTTP statuscode 200 and an JSON payload indicating the send message and your remaining quota
            </p>
            <pre class="red-code">{
    "success":true,
    "message":"Message sent",
    "response":
    {
        "multicast_id":8000000000000000006,
        "success":1,
        "failure":0,
        "canonical_ids":0,
        "results": [{"message_id":"0:10000000000000000000000000000000d"}]
    },
    "messagecount":623,
    "quota":17,
    "quota_max":100
}</pre>
            <p>
                If the operation is <b>not</b> successful the API will respond with an 4xx HTTP statuscode.
            </p>
            <table class="scode_table">
                <thead>
                <tr>
                    <th>Statuscode</th>
                    <th>Explanation</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td data-label="Statuscode">200 (OK)</td>
                    <td data-label="Explanation">Message sent</td>
                </tr>
                <tr>
                    <td data-label="Statuscode">400 (Bad Request)</td>
                    <td data-label="Explanation">The request is invalid (missing parameters or wrong values)</td>
                </tr>
                <tr>
                    <td data-label="Statuscode">401 (Unauthorized)</td>
                    <td data-label="Explanation">The user_id was not found or the user_key is wrong</td>
                </tr>
                <tr>
                    <td data-label="Statuscode">403 (Forbidden)</td>
                    <td data-label="Explanation">The user has exceeded its daily quota - wait 24 hours or upgrade your account</td>
                </tr>
                <tr>
                    <td data-label="Statuscode">412 (Precondition Failed)</td>
                    <td data-label="Explanation">There is no device connected with this account - open the app and press the refresh button in the account tab</td>
                </tr>
                <tr>
                    <td data-label="Statuscode">500 (Internal Server Error)</td>
                    <td data-label="Explanation">There was an internal error while sending your data - try again later</td>
                </tr>
                </tbody>
            </table>
            <p>
                There is also always a JSON payload with additional information.
                The <code>success</code> field is always there and in the error state you the <code>message</code> field to get a descritpion of the problem.
            </p>
            <pre class="red-code">{
    "success":false,
    "error":2101,
    "errhighlight":-1,
    "message":"Daily quota reached (100)"
}</pre>
        </div>

        <h2>Message Content</h2>
        <div class="section">
            <p>
                Every message must have a title set.
                But you also (optionally) add more content, while the title has a max length of 120 characters, the conntent can be up to 10.000 characters.
                You can see the whole message with title and content in the app or when clicking on the notification.
            </p>
            <p>
               If needed the content can be supplied in the <code>content</code> parameter.
            </p>
            <pre>curl                                          \
    --data "user_id={userid}"                 \
    --data "user_key={userkey}"               \
    --data "title={message_title}"            \
    --data "content={message_content}"        \
    https://scn.blackforestbytes.com/send.php</pre>
        </div>

        <h2>Message Priority</h2>
        <div class="section">
            <p>
                Currently you can send a message with three different priorities: 0 (low), 1 (normal) and 2 (high).
                In the app you can then configure a different behaviour for different priorities, e.g. only playing a sound if the notification is high priority.
            </p>
            <p>
                Priorites are either 0, 1 or 2 and are supplied in the <code>priority</code> parameter.
                If no priority is supplied the message will get the default priority of 1.
            </p>
            <pre>curl                                          \
    --data "user_id={userid}"                 \
    --data "user_key={userkey}"               \
    --data "title={message_title}"            \
    --data "priority={0|1|2}"                 \
    https://scn.blackforestbytes.com/send.php</pre>
        </div>

        <h2>Message Uniqueness</h2>
        <div class="section">
            <p>
                Sometimes your script can run in an environment with an unstable connection and you want to implement an automatic re-try mechanism to send a message again if the last try failed due to bad connectivity.
            </p>
            <p>
                To ensure that a message is only send once you can generate a unique id for your message (I would recommend a simple <code>uuidgen</code>).
                If you send a message with an UUID that was already used in the near past the API still returns OK, but no new message is sent.
            </p>
            <p>
                The message_id is optional - but if you want to use it you need to supply it via the <code>msg_id</code> parameter.
            </p>
            <pre>curl                                          \
    --data "user_id={userid}"                 \
    --data "user_key={userkey}"               \
    --data "title={message_title}"            \
    --data "msg_id={message_id}"              \
    https://scn.blackforestbytes.com/send.php</pre>
            <p>
                Be aware that the server only saves send messages for a short amount of time. Because of that you can only use this to prevent duplicates in a short time-frame, older messages with the same ID are probably already deleted and the message will be send again.
            </p>
        </div>
    </div>

</body>
</html>