<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);


if (!isset($INPUT['user_id']))  die(json_encode(['success' => false, 'errid'=>101, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key'])) die(json_encode(['success' => false, 'errid'=>102, 'message' => 'Missing parameter [[user_key]]']));

$user_id   = $INPUT['user_id'];
$user_key  = $INPUT['user_key'];

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, quota_today, is_pro, quota_day, fcm_token FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'errid'=>201, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'errid'=>202, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'errid'=>203, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'errid'=>204, 'message' => 'Authentification failed']));

//-------------------

$stmt = $pdo->prepare('SELECT * FROM messages WHERE ack=0 AND sender_user_id=:uid ORDER BY `timestamp` DESC LIMIT 16');
$stmt->execute(['uid' => $user_id]);
$nonacks_sql = $stmt->fetchAll(PDO::FETCH_ASSOC);

$nonacks = [];

foreach ($nonacks_sql as $nack)
{
	$nonacks []=
	[
		'title'      => $nack['title'],
		'body'       => $nack['content'],
		'priority'   => $nack['priority'],
		'timestamp'  => $nack['timestamp'],
		'usr_msg_id' => $nack['usr_message_id'],
		'scn_msg_id' => $nack['scn_message_id'],
	];
}

api_return(200,
[
	'success'        => true,
	'message'        => 'ok',
	'count'          => count($nonacks),
	'data'           => $nonacks,
]);
