<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);


if (!isset($INPUT['user_id']))   die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key']))  die(json_encode(['success' => false, 'message' => 'Missing parameter [[user_key]]']));
if (!isset($INPUT['pro']))       die(json_encode(['success' => false, 'message' => 'Missing parameter [[pro]]']));
if (!isset($INPUT['pro_token'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[pro_token]]']));

$user_id   = $INPUT['user_id'];
$user_key  = $INPUT['user_key'];
$ispro     = $INPUT['pro'] == 'true';
$pro_token = $INPUT['pro_token'];

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, quota_today, quota_day, is_pro, pro_token FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'message' => 'Authentification failed']));

if ($ispro)
{
	// set pro=true

	if ($data['pro_token'] != $pro_token)
	{
		if (!verifyOrderToken($pro_token)) die(json_encode(['success' => false, 'message' => 'Purchase token could not be verified']));
	}

	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), is_pro=1, pro_token=:ptk WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'ptk' => $pro_token]);

	$stmt = $pdo->prepare('UPDATE users SET is_pro=0, pro_token=NULL WHERE user_id <> :uid AND pro_token = :ptk');
	$stmt->execute(['uid' => $user_id, 'ptk' => $pro_token]);

	api_return(200,
	[
		'success'  => true,
		'user_id'  => $user_id,
		'quota'    => $data['quota_today'],
		'quota_max'=> Statics::quota_max(true),
		'is_pro'   => true,
		'message'  => 'user updated'
	]);
}
else
{
	// set pro=false

	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), is_pro=0, pro_token=NULL WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id]);

	api_return(200,
	[
		'success'  => true,
		'user_id'  => $user_id,
		'quota'    => $data['quota_today'],
		'quota_max'=> Statics::quota_max(false),
		'is_pro'   => false,
		'message'  => 'user updated'
	]);
}
