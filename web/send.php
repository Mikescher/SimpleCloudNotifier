<?php

include_once 'model.php';

//------------------------------------------------------------------
sleep(1);
//------------------------------------------------------------------

$INPUT = array_merge($_GET, $_POST);

if (!isset($INPUT['user_id']))  die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key'])) die(json_encode(['success' => false, 'errhighlight' => 102, 'message' => 'Missing parameter [[user_token]]']));
if (!isset($INPUT['title']))    die(json_encode(['success' => false, 'errhighlight' => 103, 'message' => 'Missing parameter [[title]]']));

//------------------------------------------------------------------

$user_id  = $INPUT['user_id'];
$user_key = $INPUT['user_key'];
$message  = $INPUT['title'];
$content  = isset($INPUT['content'])  ? $INPUT['content']  : '';
$priority = isset($INPUT['priority']) ? $INPUT['priority'] : '1';

//------------------------------------------------------------------

if ($priority !== '0' && $priority !== '1' && $priority !== '2') die(json_encode(['success' => false, 'errhighlight' => 105, 'message' => 'Invalid priority']));

if (strlen(trim($message)) == 0) die(json_encode(['success' => false, 'errhighlight' => 103, 'message' => 'No title specified']));
if (strlen($message) > 120)      die(json_encode(['success' => false, 'errhighlight' => 103, 'message' => 'Title too long (120 characters)']));
if (strlen($content) > 10000)    die(json_encode(['success' => false, 'errhighlight' => 104, 'message' => 'Content too long (10000 characters)']));

//------------------------------------------------------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, fcm_token, messages_sent, quota_today, is_pro, quota_day FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'errhighlight' => 101, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'errhighlight' => 102, 'message' => 'Authentification failed']));

$fcm = $data['fcm_token'];

$new_quota = $data['quota_today'] + 1;
if ($data['quota_day'] === null || $data['quota_day'] !== date("Y-m-d")) $new_quota=1;
if ($new_quota > Statics::quota_max($data['is_pro'])) die(json_encode(['success' => false, 'errhighlight' => -1, 'message' => 'Daily quota reached ('.Statics::quota_max($data['is_pro']).')']));

//------------------------------------------------------------------

$url = "https://fcm.googleapis.com/fcm/send";
$payload = json_encode(
[
	'to' => $fcm,
	//'dry_run' => true,
	'android' => [ 'priority' => 'high' ],
	//'notification' =>
	//[
	//	'title' => $message,
	//	'body' => $content,
	//],
	'data' =>
	[
		'title' => $message,
		'body' => $content,
		'priority' => $priority,
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

$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), messages_sent=messages_sent+1, quota_today=:q, quota_day=NOW() WHERE user_id = :uid');
$stmt->execute(['uid' => $user_id, 'q' => $new_quota]);

echo (json_encode(
[
	'success'      => true,
	'message'      => 'Message sent',
	'response'     => $httpresult,
	'messagecount' => $data['messages_sent']+1,
	'quota'        => $new_quota,
	'is_pro'       => $data['is_pro'],
	'quota_max'    => Statics::quota_max($data['is_pro']),
]));
return 0;