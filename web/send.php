<?php

include_once 'model.php';


if (!isset($_GET['user_id']))       die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($_GET['user_key']))      die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_token]]']));
if (!isset($_GET['message_title'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[message_title]]']));

$user_id  = $_GET['user_id'];
$user_key = $_GET['user_key'];
$message  = $_GET['message_title'];
$content  = isset($_POST['message_content']) ? $_POST['message_content'] : '';

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, fcm_token FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'message' => 'No User found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'message' => 'Authentification failed']));

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

echo (json_encode(['success' => true, 'message' => 'Message sent', 'response' => $httpresult]));
return 0;