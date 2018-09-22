<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);

if (!isset($INPUT['user_id']))       die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key']))      die(json_encode(['success' => false, 'errhighlight' => 102, 'message' => 'Missing parameter [[user_token]]']));
if (!isset($INPUT['message_title'])) die(json_encode(['success' => false, 'errhighlight' => 103, 'message' => 'Missing parameter [[message_title]]']));

$user_id  = $INPUT['user_id'];
$user_key = $INPUT['user_key'];
$message  = $INPUT['message_title'];
$content  = file_get_contents('php://input');
if ($content === null || $content === false) $content = '';

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, fcm_token, messages_sent FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'No User found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'errhighlight' => 102, 'message' => 'Authentification failed']));

$fcm = $data['fcm_token'];


$url = "https://fcm.googleapis.com/fcm/send";
$payload = json_encode(
[
	'to' => $fcm,
	//'dry_run' => true,
	'notification' =>
	[
		'title' => $message,
		'body' => $content,
	],
	'data' =>
	[
		'title' => $message,
		'body' => $content,
		'timestamp' => time(),
	]
]);
$header=
[
	'Authorization' => 'key=' . getConfig()['firebase']['server_key'],
	'Content-Type' => 'application/json',
];

try
{
	$httpresult = sendPOST($url, $payload, $header);
}
catch (Exception $e)
{
	die(json_encode(['success' => false, 'message' => 'Exception: ' . $e->getMessage()]));
}

$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), messages_sent=messages_sent+1 WHERE user_id = :uid');
$stmt->execute(['uid' => $user_id]);

echo (json_encode(
[
	'success' => true,
	'message' => 'Message sent',
	'response' => $httpresult,
	'messagecount' => $data['messages_sent']+1
]));
return 0;