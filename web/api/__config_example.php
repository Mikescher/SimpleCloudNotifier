<?php

// insert your values here and rename to config.php

return
[
	'global' =>
	[
		'prod' => true,
	],

	'database' =>
	[
		'host'     => '?',
		'database' => '?',
		'user'     => '?',
		'password' => '?',
	],

	'firebase' =>
	[
		'type'                        => 'service_account',
		'project_id'                  => '?',
		'private_key_id'              => '???',
		'client_email'                => '???.iam.gserviceaccount.com',
		'client_id'                   => '???',
		'auth_uri'                    => 'https://accounts.google.com/o/oauth2/auth',
		'token_uri'                   => 'https://oauth2.googleapis.com/token',
		'auth_provider_x509_cert_url' => 'https://www.googleapis.com/oauth2/v1/certs',
		'client_x509_cert_url'        => 'https://www.googleapis.com/robot/v1/metadata/x509/???f.iam.gserviceaccount.com',
		'private_key'                 => "-----BEGIN PRIVATE KEY-----\n"
		                               . "??????????\n"
		                               . "-----END PRIVATE KEY-----\n",
		'server_key'                  => '????',
	],

	'verify_api' =>
	[
		'package_name'   => 'com.blackforestbytes.simplecloudnotifier',
		'product_id'     => '???',

		'clientid'       => '???.apps.googleusercontent.com',
		'clientsecret'   => '???',
		'accesstoken'    => file_exists('.verify_accesstoken') ? file_get_contents('.verify_accesstoken') : '',
		'refreshtoken'   => '???',
		'scope'          => 'https://www.googleapis.com/auth/androidpublisher',
	],

	'error_reporting' =>
	[
		'send-mail' => true,
		'email-error-target' => '???@???.com',
		'email-error-sender' => '???@???.com',
	],

];