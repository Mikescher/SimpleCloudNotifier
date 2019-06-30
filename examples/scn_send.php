<?php

/**
 * @param string $title
 * @param string $content
 * @param int $priority
 * @return bool
 */
function sendSCN($title, $content, $priority) {
	global $config;

	$data =
	[
		'user_id'  => '', //TODO set your userid
		'user_key' => '', //TODO set your userkey
		'title'    => $title,
		'content'  => $content,
		'priority' => $priority,
	];

	$ch = curl_init();

	curl_setopt($ch, CURLOPT_URL, "https://simplecloudnotifier.blackforestbytes.com/send.php");
	curl_setopt($ch, CURLOPT_POST, 1);
	curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($data));
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

	$result = curl_exec($ch);

	curl_close($ch);
	if ($result === false) return false;

	$json = json_decode($result, true);
	return $json['success'];
}