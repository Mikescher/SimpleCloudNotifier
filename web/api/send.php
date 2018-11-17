<?php

include_once 'model.php';

try
{

//------------------------------------------------------------------
//sleep(1);
//------------------------------------------------------------------

	if ($_SERVER['REQUEST_METHOD'] !== 'POST') api_return(400, json_encode(['success' => false, 'error' =>  ERR::REQ_METHOD, 'errhighlight' => -1, 'message' => 'Invalid request method (must be POST)']));

	$INPUT = array_merge($_GET, $_POST);

	if (!isset($INPUT['user_id']))  api_return(400, json_encode(['success' => false, 'error' => ERR::MISSING_UID,   'errhighlight' => 101, 'message' => 'Missing parameter [[user_id]]']));
	if (!isset($INPUT['user_key'])) api_return(400, json_encode(['success' => false, 'error' => ERR::MISSING_TOK,   'errhighlight' => 102, 'message' => 'Missing parameter [[user_token]]']));
	if (!isset($INPUT['title']))    api_return(400, json_encode(['success' => false, 'error' => ERR::MISSING_TITLE, 'errhighlight' => 103, 'message' => 'Missing parameter [[title]]']));

//------------------------------------------------------------------


	$user_id  = $INPUT['user_id'];
	$user_key = $INPUT['user_key'];
	$message  = $INPUT['title'];
	$content  = isset($INPUT['content'])  ? $INPUT['content']  : '';
	$priority = isset($INPUT['priority']) ? $INPUT['priority'] : '1';
	$usrmsgid = isset($INPUT['msg_id'])   ? $INPUT['msg_id'] : null;

//------------------------------------------------------------------

	if ($priority !== '0' && $priority !== '1' && $priority !== '2') api_return(400, json_encode(['success' => false, 'error' => ERR::INVALID_PRIO, 'errhighlight' => 105, 'message' => 'Invalid priority']));

	if (strlen(trim($message)) == 0)                 api_return(400, json_encode(['success' => false, 'error' => ERR::NO_TITLE,            'errhighlight' => 103, 'message' => 'No title specified']));
	if (strlen($message) > 120)                      api_return(400, json_encode(['success' => false, 'error' => ERR::TITLE_TOO_LONG,      'errhighlight' => 103, 'message' => 'Title too long (120 characters)']));
	if (strlen($content) > 10000)                    api_return(400, json_encode(['success' => false, 'error' => ERR::CONTENT_TOO_LONG,    'errhighlight' => 104, 'message' => 'Content too long (10000 characters)']));
	if ($usrmsgid != null && strlen($usrmsgid) > 64) api_return(400, json_encode(['success' => false, 'error' => ERR::USR_MSG_ID_TOO_LONG, 'errhighlight' => -1,  'message' => 'MessageID too long (64 characters)']));

//------------------------------------------------------------------

	$pdo = getDatabase();

	$stmt = $pdo->prepare('SELECT user_id, user_key, fcm_token, messages_sent, quota_today, is_pro, quota_day FROM users WHERE user_id = :uid LIMIT 1');
	$stmt->execute(['uid' => $user_id]);

	$datas = $stmt->fetchAll(PDO::FETCH_ASSOC);
	if (count($datas)<=0) die(json_encode(['success' => false, 'error' => ERR::USER_NOT_FOUND, 'errhighlight' => 101, 'message' => 'User not found']));
	$data = $datas[0];

	if ($data === null)                     api_return(401, json_encode(['success' => false, 'error' => ERR::USER_NOT_FOUND,   'errhighlight' => 101, 'message' => 'User not found']));
	if ($data['user_id'] !== (int)$user_id) api_return(401, json_encode(['success' => false, 'error' => ERR::USER_NOT_FOUND,   'errhighlight' => 101, 'message' => 'UserID not found']));
	if ($data['user_key'] !== $user_key)    api_return(401, json_encode(['success' => false, 'error' => ERR::USER_AUTH_FAILED, 'errhighlight' => 102, 'message' => 'Authentification failed']));

	$fcm = $data['fcm_token'];

	$new_quota = $data['quota_today'] + 1;
	if ($data['quota_day'] === null || $data['quota_day'] !== date("Y-m-d")) $new_quota=1;
	if ($new_quota > Statics::quota_max($data['is_pro'])) api_return(403, json_encode(['success' => false, 'error' => ERR::QUOTA_REACHED, 'errhighlight' => -1, 'message' => 'Daily quota reached ('.Statics::quota_max($data['is_pro']).')']));

	if ($fcm == null || $fcm == '' ||  $fcm == false)
	{
		api_return(412, json_encode(['success' => false, 'error' => ERR::NO_DEVICE_LINKED, 'errhighlight' => -1, 'message' => 'No device linked with this account']));
	}

//------------------------------------------------------------------

	if ($usrmsgid != null)
	{
		$stmt = $pdo->prepare('SELECT scn_message_id FROM messages WHERE sender_user_id=:uid AND usr_message_id IS NOT NULL AND usr_message_id=:umid LIMIT 1');
		$stmt->execute(['uid' => $user_id, 'umid' => $usrmsgid]);

		if (count($stmt->fetchAll(PDO::FETCH_ASSOC))>0)
		{
			api_return(200, json_encode(
			[
				'success'       => true,
				'message'       => 'Message already sent',
				'suppress_send' => true,
				'response'      => '',
				'messagecount'  => $data['messages_sent']+1,
				'quota'         => $data['quota_today'],
				'is_pro'        => $data['is_pro'],
				'quota_max'     => Statics::quota_max($data['is_pro']),
			]));
		}
	}

//------------------------------------------------------------------

	$pdo->beginTransaction();


	$stmt = $pdo->prepare('INSERT INTO messages (sender_user_id, title, content, priority, fcm_message_id, usr_message_id) VALUES (:suid, :t, :c, :p, :fmid, :umid)');
	$stmt->execute(
	[
		'suid' => $user_id,
		't'    => $message,
		'c'    => $content,
		'p'    => $priority,
		'fmid' => null,
		'umid' => $usrmsgid,
	]);

	$scn_msg_id = $pdo->lastInsertId();

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
				'title'      => $message,
				'body'       => $content,
				'priority'   => $priority,
				'timestamp'  => time(),
				'usr_msg_id' => $usrmsgid,
				'scn_msg_id' => $scn_msg_id,
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

		if (try_json($httpresult, ['success']) != 1)
		{
			reportError("FCM communication failed (success_1 <> true)\n\n".$httpresult);
			$pdo->rollBack();
			api_return(500, json_encode(['success' => false, 'error' => ERR::FIREBASE_COM_ERRORED, 'errhighlight' => -1, 'message' => 'Communication with firebase service failed.']));
		}
	}
	catch (Exception $e)
	{
		reportError("FCM communication failed", $e);
		$pdo->rollBack();
		api_return(500, json_encode(['success' => false, 'error' => ERR::FIREBASE_COM_FAILED, 'errhighlight' => -1, 'message' => 'Communication with firebase service failed.'."\n\n".'Exception: ' . $e->getMessage()]));
	}

	$stmt = $pdo->prepare('UPDATE users SET timestamp_accessed=NOW(), messages_sent=messages_sent+1, quota_today=:q, quota_day=NOW() WHERE user_id = :uid');
	$stmt->execute(['uid' => $user_id, 'q' => $new_quota]);

	$stmt = $pdo->prepare('UPDATE messages SET fcm_message_id=:fmid WHERE scn_message_id=:smid');
	$stmt->execute([ 'fmid' => try_json($httpresult, ['results', 0, 'message_id']), 'smid' => $scn_msg_id ]);

	$pdo->commit();

	api_return(200, json_encode(
	[
		'success'       => true,
		'error'         => ERR::NO_ERROR,
		'errhighlight'  => -1,
		'message'       => 'Message sent',
		'suppress_send' => false,
		'response'      => $httpresult,
		'messagecount'  => $data['messages_sent']+1,
		'quota'         => $new_quota,
		'is_pro'        => $data['is_pro'],
		'quota_max'     => Statics::quota_max($data['is_pro']),
		'scn_msg_id'    => $scn_msg_id,
	]));
}
catch (Exception $mex)
{
	reportError("Root try-catch triggered", $mex);
	if ($pdo->inTransaction()) $pdo->rollBack();
	api_return(500, json_encode(['success' => false, 'error' => ERR::INTERNAL_EXCEPTION, 'errhighlight' => -1, 'message' => 'PHP script threw exception.'."\n\n".'Exception: ' . $e->getMessage()]));
}
