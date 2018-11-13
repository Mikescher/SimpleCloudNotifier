<?php

include('lib/httpful.phar');

class ERR
{
	const NO_ERROR             = 0000;

	const MISSING_UID          = 1101;
	const MISSING_TOK          = 1102;
	const MISSING_TITLE        = 1103;
	const INVALID_PRIO         = 1104;
	const REQ_METHOD           = 1105;

	const NO_TITLE             = 1201;
	const TITLE_TOO_LONG       = 1202;
	const CONTENT_TOO_LONG     = 1203;
	const USR_MSG_ID_TOO_LONG  = 1204;

	const USER_NOT_FOUND       = 1301;
	const USER_AUTH_FAILED     = 1302;

	const NO_DEVICE_LINKED     = 1401;

	const QUOTA_REACHED        = 2101;

	const FIREBASE_COM_FAILED  = 9901;
	const FIREBASE_COM_ERRORED = 9902;
	const INTERNAL_EXCEPTION   = 9903;
}

class Statics
{
	public static $DB = NULL;
	public static $CFG = NULL;

	public static function quota_max($is_pro) { return $is_pro ? 1000 : 100; }
}

function getConfig()
{
	if (Statics::$CFG !== NULL) return Statics::$CFG;

	return Statics::$CFG = require "config.php";
}

/**
 * @param String $msg
 * @param Exception $e
 */
function reportError($msg, $e = null)
{
	if ($e != null) $msg = ($msg."\n\n[[EXCEPTION]]\n" . $e . "\n" . $e->getMessage() . "\n" . $e->getTraceAsString());

	$subject = "SCN_Server has encountered an Error at " . date("Y-m-d H:i:s") . "] ";

	$content = "";

	$content .= 'HTTP_HOST: '            . ParamServerOrUndef('HTTP_HOST')            . "\n";
	$content .= 'REQUEST_URI: '          . ParamServerOrUndef('REQUEST_URI')          . "\n";
	$content .= 'TIME: '                 . date('Y-m-d H:i:s')                        . "\n";
	$content .= 'REMOTE_ADDR: '          . ParamServerOrUndef('REMOTE_ADDR')          . "\n";
	$content .= 'HTTP_X_FORWARDED_FOR: ' . ParamServerOrUndef('HTTP_X_FORWARDED_FOR') . "\n";
	$content .= 'HTTP_USER_AGENT: '      . ParamServerOrUndef('HTTP_USER_AGENT')      . "\n";
	$content .= 'MESSAGE:'               . "\n" . $msg                                . "\n";
	$content .= '$_GET:'                 . "\n" . print_r($_GET, true)                . "\n";
	$content .= '$_POST:'                . "\n" . print_r($_POST, true)               . "\n";
	$content .= '$_FILES:'               . "\n" . print_r($_FILES, true)              . "\n";

	if (getConfig()['error_reporting']['send-mail']) sendMail($subject, $content, getConfig()['error_reporting']['email-error-target'], getConfig()['error_reporting']['email-error-sender']);
}

/**
 * @param string $subject
 * @param string $content
 * @param string $to
 * @param string $from
 */
function sendMail($subject, $content, $to, $from) {
	mail($to, $subject, $content, 'From: ' . $from);
}

/**
 * @param string $idx
 * @return string
 */
function ParamServerOrUndef($idx) {
	return isset($_SERVER[$idx]) ? $_SERVER[$idx] : 'NOT_SET';
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

/**
 * @param $url
 * @param $body
 * @param $header
 * @return array|object|string
 * @throws \Httpful\Exception\ConnectionErrorException
 * @throws Exception
 */
function sendPOST($url, $body, $header)
{
	$builder = \Httpful\Request::post($url);

	$builder->body($body);

	foreach ($header as $k => $v) $builder->addHeader($k, $v);

	$response = $builder->send();

	if ($response->code != 200) throw new Exception("Repsponse code: " . $response->code);

	return $response->raw_body;
}

function verifyOrderToken($tok)
{
	// https://developers.google.com/android-publisher/api-ref/purchases/products/get

	try
	{
		$package  = getConfig()['verify_api']['package_name'];
		$product  = getConfig()['verify_api']['product_id'];
		$acctoken = getConfig()['verify_api']['accesstoken'];

		if ($acctoken == '' || $acctoken == null || $acctoken == false) $acctoken = refreshVerifyToken();

		$url = 	'https://www.googleapis.com/androidpublisher/v3/applications/'.$package.'/purchases/products/'.$product.'/tokens/'.$tok.'?access_token='.$acctoken;
		$response = $builder = \Httpful\Request::get($url)->send();
		$obj = json_decode($response->raw_body, true);

		if ($response->code != 401 && ($obj === null || $obj === false))
		{
			reportError('verify-token returned NULL');
			return false;
		}

		if ($response->code == 401 || isset($obj['error']) && isset($obj['error']['code']) && $obj['error']['code'] == 401) // "Invalid Credentials" -- refresh acces_token
		{
			$acctoken = refreshVerifyToken();

			$url = 	'https://www.googleapis.com/androidpublisher/v3/applications/'.$package.'/purchases/products/'.$product.'/tokens/'.$tok.'?access_token='.$acctoken;
			$response = $builder = \Httpful\Request::get($url)->send();
			$obj = json_decode($response->raw_body, true);

			if ($obj === null || $obj === false)
			{
				reportError('verify-token returned NULL');
				return false;
			}
		}

		if (isset($obj['purchaseState']) && $obj['purchaseState'] === 0) return true;

		return false;
	}
	catch (Exception $e)
	{
		reportError("VerifyOrder token threw exception", $e);
		return false;
	}
}

/** @throws Exception */
function refreshVerifyToken()
{
	$url = 	'https://accounts.google.com/o/oauth2/token'.
			'?grant_type=refresh_token'.
			'&refresh_token='.getConfig()['verify_api']['refreshtoken'].
			'&client_id='.getConfig()['verify_api']['clientid'].
			'&client_secret='.getConfig()['verify_api']['clientsecret'];

	$json = sendPOST($url, "", []);
	$obj = json_decode($json, true);
	file_put_contents('.verify_accesstoken', $obj['access_token']);

	return $obj['access_token'];
}

function api_return($http_code, $message)
{
	http_response_code($http_code);
	echo $message;
	die();
}

/**
 * @param String $str
 * @param String[] $path
 * @return mixed|null
 */
function try_json($str, $path)
{
	try
	{
		$o = json_decode($str, true);
		foreach ($path as $p) $o = $o[$p];
		return $o;
	}
	catch (Exception $e)
	{
		return null;
	}
}

//#################################################################################################################

if (getConfig()['global']['prod']) {
	ini_set('display_errors', 0);
	ini_set('log_errors', 1);
} else {
	error_reporting(E_STRICT);
	ini_set('display_errors', 1);
}

//#################################################################################################################
