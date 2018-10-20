<?php

include_once 'model.php';

$INPUT = array_merge($_GET, $_POST);


if (!isset($INPUT['user_id']))  die(json_encode(['success' => false, 'errid'=>101, 'message' => 'Missing parameter [[user_id]]']));
if (!isset($INPUT['user_key'])) die(json_encode(['success' => false, 'errid'=>102, 'message' => 'Missing parameter [[user_key]]']));

$user_id   = $INPUT['user_id'];
$user_key  = $INPUT['user_key'];

//----------------------

$pdo = getDatabase();

$stmt = $pdo->prepare('SELECT user_id, user_key, quota_today, quota_max, quota_day FROM users WHERE user_id = :uid LIMIT 1');
$stmt->execute(['uid' => $user_id]);

$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
if (count($datas)<=0) die(json_encode(['success' => false, 'errid'=>201, 'message' => 'User not found']));
$data = $datas[0];

if ($data === null) die(json_encode(['success' => false, 'errid'=>202, 'message' => 'User not found']));
if ($data['user_id'] !== (int)$user_id) die(json_encode(['success' => false, 'errid'=>203, 'message' => 'UserID not found']));
if ($data['user_key'] !== $user_key) die(json_encode(['success' => false, 'errid'=>204, 'message' => 'Authentification failed']));

$quota = $data['quota_today'];
$quota_max = $data['quota_max'];

if ($data['quota_day'] === null || $data['quota_day'] !== date("Y-m-d")) $quota=0;

echo json_encode(
[
	'success'  => true,
	'user_id'  => $user_id,
	'quota'    => $quota,
	'quota_max'=> $quota_max,
	'message'  => 'ok'
]);
return 0;