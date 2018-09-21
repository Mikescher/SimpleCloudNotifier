<?php

include_once 'model.php';


if (!isset($_GET['user_id']))  die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($_GET['user_key'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_key]]']));
if (!isset($_GET['message']))  die(json_encode(['success' => false, 'message' => 'Missing parameter [[message]]']));

$user_id   = $_GET['user_id'];
$user_key  = $_GET['token'];
$fcm_token = isset($_GET['token']) ? $_GET['token'] : null;

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'message' => 'No User found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'message' => 'Authentification failed']));

$new_userkey = generateRandomAuthKey();

if ($fcm_token === null)
{
	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), user_key=:at WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'at' => $new_userkey]);

	echo json_encode(['success' => true, 'user_id' => $user_id, 'user_key' => $new_userkey, 'message' => 'user updated']);
	return 0;
}
else
{
	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), fcm_token=:ft, user_key=:at WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'ft' => $fcm_token, 'at' => $new_userkey]);

	echo json_encode(['success' => true, 'user_id' => $user_id, 'user_key' => $new_userkey, 'message' => 'user updated']);
	return 0;
}