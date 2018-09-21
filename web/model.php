<?php

include('lib/httpful.phar');

class Statics
{
	public static $DB = NULL;
	public static $CFG = NULL;
}

function getConfig()
{
	if (Statics::$CFG !== NULL) return Statics::$CFG;

	return Statics::$CFG = require "config.php";
}

function getDatabase()
{
	if (Statics::$DB !== NULL) return Statics::$DB;

	$_config = getConfig()['database'];

	$dsn = "mysql:host=" . $_config['host'] . ";dbname=" . $_config['database'] . ";charset=utf8";
	$opt = [
		PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
		PDO::ATTR_EMULATE_PREPARES   => false,
	];

	return Statics::$DB = new PDO($dsn, $_config['user'], $_config['password'], $opt);
}

function generateRandomAuthKey()
{
	$random = '';
	for ($i = 0; $i < 64; $i++)

		try {
			switch (random_int(1, 3)) {
				case 1:
					$random .= chr(random_int(ord('0'), ord('9')));
					break;
				case 2:
					$random .= chr(random_int(ord('A'), ord('Z')));
					break;
				case 3:
					$random .= chr(random_int(ord('a'), ord('z')));
					break;

			}
		}
		catch (Exception $e)
		{
			die(json_encode(['success' => false, 'message' => 'Internal error - no randomness']));
		}
	return $random;
}

function sendPOST($url, $body, $header)
{
	$builder = \Httpful\Request::post($url);

	$builder->body($body);

	foreach ($header as $k => $v) $builder->addHeader($k, $v);

	$response = $builder->send();

	if ($response->code != 200) throw new Exception("Repsponse code: " . $response->code);

	return $response->body;
}
