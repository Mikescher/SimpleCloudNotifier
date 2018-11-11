<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);

if (!isset($INPUT['fcm_token'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[fcm_token]]']));
if (!isset($INPUT['pro']))       die(json_encode(['success' => false, 'message' => 'Missing parameter [[pro]]']));
if (!isset($INPUT['pro_token'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[pro_token]]']));

$fcmtoken  = $INPUT['fcm_token'];
$ispro     = $INPUT['pro'] == 'true';
$pro_token = $INPUT['pro_token'];
$user_key = generateRandomAuthKey();

$pdo = getDatabase();

if ($ispro)
{
	if (!verifyOrderToken($pro_token)) die(json_encode(['success' => false, 'message' => 'Purchase token could not be verified']));
}

$stmt = $pdo->prepare('INSERT INTO users (user_key, fcm_token, is_pro, pro_token, timestamp_accessed) VALUES (:key, :token, :bpro, :spro, NOW())');
$stmt->execute(['key' => $user_key, 'token' => $fcmtoken, 'bpro' => $ispro, 'spro' => $ispro ? $pro_token : null]);
$user_id = $pdo->lastInsertId('user_id');

echo json_encode(
[
	'success'   => true,
	'user_id'   => $user_id,
	'user_key'  => $user_key,
	'quota'     => 0,
	'quota_max' => Statics::quota_max($ispro),
	'is_pro'    => $ispro,
	'message'   => 'New user registered'
]);

return 0;