import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/models/api_error.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';
import 'package:xid/xid.dart';

part 'request_log.g.dart';

class RequestLog {
  //TODO max size, auto clear old

  static void addRequestException(String name, DateTime tStart, String method, Uri uri, String reqbody, Map<String, String> reqheaders, dynamic e, StackTrace trace) {
    Hive.box<SCNRequest>('scn-requests').add(SCNRequest(
      id: Xid().toString(),
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
      id: Xid().toString(),
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
      id: Xid().toString(),
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
      id: Xid().toString(),
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
      id: Xid().toString(),
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
class SCNRequest extends HiveObject implements FieldDebuggable {
  @HiveField(0)
  final String id;

  @HiveField(10)
  final DateTime timestampStart;
  @HiveField(11)
  final DateTime timestampEnd;
  @HiveField(12)
  final String name;
  @HiveField(13)
  final String type; // SUCCESS | EXCEPTION | API_ERROR | ERROR_STATUSCODE | DECODE_ERROR
  @HiveField(14)
  final String error;
  @HiveField(15)
  final String stackTrace;

  @HiveField(21)
  final String method;
  @HiveField(22)
  final String url;
  @HiveField(23)
  final Map<String, String> requestHeaders;
  @HiveField(24)
  final String requestBody;

  @HiveField(31)
  final int responseStatusCode;
  @HiveField(32)
  final Map<String, String> responseHeaders;
  @HiveField(33)
  final String responseBody;

  SCNRequest({
    required this.id,
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

  @override
  String toString() {
    return 'SCNRequest[${this.id}]';
  }

  List<(String, String)> debugFieldList() {
    return [
      ('id', this.id),
      ('timestampStart', this.timestampStart.toIso8601String()),
      ('timestampEnd', this.timestampEnd.toIso8601String()),
      ('name', this.name),
      ('method', this.method),
      ('url', this.url),
      for (var (idx, item) in this.requestHeaders.entries.indexed) ('requestHeaders[$idx]', '${item.key}=${item.value}'),
      ('requestBody', this.requestBody),
      ('responseStatusCode', this.responseStatusCode.toString()),
      for (var (idx, item) in this.responseHeaders.entries.indexed) ('responseHeaders[$idx]', '${item.key}=${item.value}'),
      ('responseBody', this.responseBody),
      ('type', this.type),
      ('error', this.error),
      ('stackTrace', this.stackTrace),
    ];
  }
}
