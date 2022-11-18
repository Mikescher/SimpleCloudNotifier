package apierr

type APIError int

//goland:noinspection GoSnakeCaseUsage
const (
	NO_ERROR APIError = 0000

	MISSING_UID          APIError = 1101
	MISSING_TOK          APIError = 1102
	MISSING_TITLE        APIError = 1103
	INVALID_PRIO         APIError = 1104
	REQ_METHOD           APIError = 1105
	INVALID_CLIENTTYPE   APIError = 1106
	BINDFAIL_QUERY_PARAM APIError = 1151
	BINDFAIL_BODY_PARAM  APIError = 1152
	BINDFAIL_URI_PARAM   APIError = 1153

	NO_TITLE               APIError = 1201
	TITLE_TOO_LONG         APIError = 1202
	CONTENT_TOO_LONG       APIError = 1203
	USR_MSG_ID_TOO_LONG    APIError = 1204
	TIMESTAMP_OUT_OF_RANGE APIError = 1205

	USER_NOT_FOUND   APIError = 1301
	USER_AUTH_FAILED APIError = 1302

	NO_DEVICE_LINKED APIError = 1401

	QUOTA_REACHED APIError = 2101

	FAILED_VERIFY_PRO_TOKEN APIError = 3001
	INVALID_PRO_TOKEN       APIError = 3002

	COMMIT_FAILED   = 9001
	DATABASE_ERROR  = 9002
	PERM_QUERY_FAIL = 9003

	FIREBASE_COM_FAILED  APIError = 9901
	FIREBASE_COM_ERRORED APIError = 9902
	INTERNAL_EXCEPTION   APIError = 9903
)
