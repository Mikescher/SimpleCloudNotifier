class APIError {
  final bool success;
  final int error;
  final String errhighlight;
  final String message;

  static final MISSING_UID = 1101;
  static final MISSING_TOK = 1102;
  static final MISSING_TITLE = 1103;
  static final INVALID_PRIO = 1104;
  static final REQ_METHOD = 1105;
  static final INVALID_CLIENTTYPE = 1106;
  static final PAGETOKEN_ERROR = 1121;
  static final BINDFAIL_QUERY_PARAM = 1151;
  static final BINDFAIL_BODY_PARAM = 1152;
  static final BINDFAIL_URI_PARAM = 1153;
  static final INVALID_BODY_PARAM = 1161;
  static final INVALID_ENUM_VALUE = 1171;

  static final NO_TITLE = 1201;
  static final TITLE_TOO_LONG = 1202;
  static final CONTENT_TOO_LONG = 1203;
  static final USR_MSG_ID_TOO_LONG = 1204;
  static final TIMESTAMP_OUT_OF_RANGE = 1205;
  static final SENDERNAME_TOO_LONG = 1206;
  static final CHANNEL_TOO_LONG = 1207;
  static final CHANNEL_DESCRIPTION_TOO_LONG = 1208;
  static final CHANNEL_NAME_EMPTY = 1209;

  static final USER_NOT_FOUND = 1301;
  static final CLIENT_NOT_FOUND = 1302;
  static final CHANNEL_NOT_FOUND = 1303;
  static final SUBSCRIPTION_NOT_FOUND = 1304;
  static final MESSAGE_NOT_FOUND = 1305;
  static final SUBSCRIPTION_USER_MISMATCH = 1306;
  static final KEY_NOT_FOUND = 1307;
  static final USER_AUTH_FAILED = 1311;

  static final NO_DEVICE_LINKED = 1401;

  static final CHANNEL_ALREADY_EXISTS = 1501;
  static final CANNOT_SELFDELETE_KEY = 1511;
  static final CANNOT_SELFUPDATE_KEY = 1512;

  static final QUOTA_REACHED = 2101;

  static final FAILED_VERIFY_PRO_TOKEN = 3001;
  static final INVALID_PRO_TOKEN = 3002;

  static final COMMIT_FAILED = 9001;
  static final DATABASE_ERROR = 9002;
  static final PERM_QUERY_FAIL = 9003;
  static final FIREBASE_COM_FAILED = 9901;
  static final FIREBASE_COM_ERRORED = 9902;
  static final INTERNAL_EXCEPTION = 9903;
  static final PANIC = 9904;
  static final NOT_IMPLEMENTED = 9905;

  const APIError({
    required this.success,
    required this.error,
    required this.errhighlight,
    required this.message,
  });

  factory APIError.fromJson(Map<String, dynamic> json) {
    return APIError(
      success: json['success'] as bool,
      error: (json['error'] as num).toInt(),
      errhighlight: json['errhighlight'] as String,
      message: json['message'] as String,
    );
  }
}
