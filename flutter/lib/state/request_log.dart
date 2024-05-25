import 'package:hive/hive.dart';
import 'package:simplecloudnotifier/models/api_error.dart';

part 'request_log.g.dart';

class RequestLog {
  static void addRequestException(String name, DateTime tStart, String method, Uri uri, String reqbody, Map<String, String> reqheaders, dynamic e, StackTrace trace) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      timestampStart: tStart,
      timestampEnd: DateTime.now(),
      name: name,
      method: method,
      url: uri.toString(),
      requestHeaders: reqheaders,
      requestBody: reqbody,
      responseStatusCode: 0,
      responseHeaders: {},
      responseBody: '',
      type: 'EXCEPTION',
      error: (e is Exception) ? e.toString() : '$e',
      stackTrace: trace.toString(),
    ));
  }

  static void addRequestAPIError(String name, DateTime t0, String method, Uri uri, String reqbody, Map<String, String> reqheaders, int responseStatusCode, String responseBody, Map<String, String> responseHeaders, APIError apierr) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      timestampStart: t0,
      timestampEnd: DateTime.now(),
      name: name,
      method: method,
      url: uri.toString(),
      requestHeaders: reqheaders,
      requestBody: reqbody,
      responseStatusCode: responseStatusCode,
      responseHeaders: responseHeaders,
      responseBody: responseBody,
      type: 'API_ERROR',
      error: apierr.message,
      stackTrace: '',
    ));
  }

  static void addRequestErrorStatuscode(String name, DateTime t0, String method, Uri uri, String reqbody, Map<String, String> reqheaders, int responseStatusCode, String responseBody, Map<String, String> responseHeaders) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      timestampStart: t0,
      timestampEnd: DateTime.now(),
      name: name,
      method: method,
      url: uri.toString(),
      requestHeaders: reqheaders,
      requestBody: reqbody,
      responseStatusCode: responseStatusCode,
      responseHeaders: responseHeaders,
      responseBody: responseBody,
      type: 'ERROR_STATUSCODE',
      error: 'API request failed with status code $responseStatusCode',
      stackTrace: '',
    ));
  }

  static void addRequestSuccess(String name, DateTime t0, String method, Uri uri, String reqbody, Map<String, String> reqheaders, int responseStatusCode, String responseBody, Map<String, String> responseHeaders) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      timestampStart: t0,
      timestampEnd: DateTime.now(),
      name: name,
      method: method,
      url: uri.toString(),
      requestHeaders: reqheaders,
      requestBody: reqbody,
      responseStatusCode: responseStatusCode,
      responseHeaders: responseHeaders,
      responseBody: responseBody,
      type: 'SUCCESS',
      error: '',
      stackTrace: '',
    ));
  }

  static void addRequestDecodeError(String name, DateTime t0, String method, Uri uri, String reqbody, Map<String, String> reqheaders, int responseStatusCode, String responseBody, Map<String, String> responseHeaders, Object exc, StackTrace trace) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      timestampStart: t0,
      timestampEnd: DateTime.now(),
      name: name,
      method: method,
      url: uri.toString(),
      requestHeaders: reqheaders,
      requestBody: reqbody,
      responseStatusCode: responseStatusCode,
      responseHeaders: responseHeaders,
      responseBody: responseBody,
      type: 'DECODE_ERROR',
      error: (exc is Exception) ? exc.toString() : '$exc',
      stackTrace: trace.toString(),
    ));
  }
}

@HiveType(typeId: 100)
class SCNRequest extends HiveObject {
  @HiveField(0)
  final DateTime timestampStart;
  @HiveField(1)
  final DateTime timestampEnd;
  @HiveField(2)
  final String name;
  @HiveField(3)
  final String type;
  @HiveField(4)
  final String error;
  @HiveField(5)
  final String stackTrace;

  @HiveField(6)
  final String method;
  @HiveField(7)
  final String url;
  @HiveField(8)
  final Map<String, String> requestHeaders;
  @HiveField(12)
  final String requestBody;

  @HiveField(9)
  final int responseStatusCode;
  @HiveField(10)
  final Map<String, String> responseHeaders;
  @HiveField(11)
  final String responseBody;

  SCNRequest({
    required this.timestampStart,
    required this.timestampEnd,
    required this.name,
    required this.method,
    required this.url,
    required this.requestHeaders,
    required this.requestBody,
    required this.responseStatusCode,
    required this.responseHeaders,
    required this.responseBody,
    required this.type,
    required this.error,
    required this.stackTrace,
  });
}
