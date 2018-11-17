<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);


if (!isset($INPUT['user_id']))  die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_key]]']));

$user_id   = $INPUT['user_id'];
$user_key  = $INPUT['user_key'];
$fcm_token = isset($INPUT['fcm_token']) ? $INPUT['fcm_token'] : null;

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, quota_today, quota_day, is_pro FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'message' => 'Authentification failed']));

$quota = $data['quota_today'];
$is_pro = $data['is_pro'];

$new_userkey = generateRandomAuthKey();

if ($fcm_token === null)
{
	// only gen new user_secret

	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), user_key=:at WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'at' => $new_userkey]);

	api_return(200,
	[
		'success'  => true,
		'user_id'  => $user_id,
		'user_key' => $new_userkey,
		'quota'    => $quota,
		'quota_max'=> Statics::quota_max($data['is_pro']),
		'is_pro'   => $is_pro,
		'message'  => 'user updated'
	]);
}
else
{
	// update fcm and gen new user_secret

	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), fcm_token=:ft, user_key=:at WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'ft' => $fcm_token, 'at' => $new_userkey]);

	$stmt = $pdo->prepare('UPDATE users SET fcm_token=NULL WHERE user_id <> :uid AND fcm_token=:ft');
	$stmt->execute(['uid' => $user_id, 'ft' => $fcm_token]);

	api_return(200,
	[
		'success' => true,
		'user_id' => $user_id,
		'user_key' => $new_userkey,
		'quota'    => $quota,
		'quota_max'=> Statics::quota_max($data['is_pro']),
		'is_pro'   => $is_pro,
		'message' => 'user updated'
	]);
}