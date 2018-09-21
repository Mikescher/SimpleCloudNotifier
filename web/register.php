<?php

include_once 'model.php';

if (!isset($_GET['fcm_token'])) die(json_encode(['success' => false, 'message' => 'Missing parameter [[fcm_token]]']));

$fcmtoken  = $_GET['fcm_token'];
$user_key = generateRandomAuthKey();

$pdo = getDatabase();

$stmt = $pdo->prepare('INSERT INTO users (user_key, fcm_token, timestamp_accessed) VALUES (:key, :token, NOW())');
$stmt->execute(['key' => $user_key, 'token' => $fcmtoken]);
$user_id = $pdo->lastInsertId('user_id');

echo json_encode(['success' => true, 'user_id' => $user_id, 'user_key' => $user_key, 'message' => 'new user registered']);
return 0;