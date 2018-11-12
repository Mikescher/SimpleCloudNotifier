<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);


if (!isset($INPUT['user_id']))    die(json_encode(['success' => false, 'errid'=>101, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key']))   die(json_encode(['success' => false, 'errid'=>102, 'message' => 'Missing parameter [[user_key]]']));
if (!isset($INPUT['scn_msg_id'])) die(json_encode(['success' => false, 'errid'=>103, 'message' => 'Missing parameter [[scn_msg_id]]']));

$user_id    = $INPUT['user_id'];
$user_key   = $INPUT['user_key'];
$scn_msg_id = $INPUT['scn_msg_id'];

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, quota_today, is_pro, quota_day, fcm_token FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);
$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);

if (count($datas)<=0) die(json_encode(['success' => false, 'errid'=>201, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null)                     die(json_encode(['success' => false, 'errid'=>202, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'errid'=>203, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key)    die(json_encode(['success' => false, 'errid'=>204, 'message' => 'Authentification failed']));

$stmt = $pdo->prepare('SELECT ack FROM messages WHERE scn_message_id=:smid AND sender_user_id=:uid LIMIT 1');
$stmt->execute(['smid' => $scn_msg_id, 'uid' => $user_id]);
$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);

if (count($datas)<=0) die(json_encode(['success' => false, 'errid'=>301, 'message' => 'Message not found']));

$stmt = $pdo->prepare('UPDATE messages SET ack=1 WHERE scn_message_id=:smid AND sender_user_id=:uid');
$stmt->execute(['smid' => $scn_msg_id, 'uid' => $user_id]);

echo json_encode(
[
	'success'        => true,
	'prev_ack'       => $datas[0]['ack'],
	'new_ack'        => 1,
	'message'        => 'ok'
]);
return 0;