import 'dart:convert';

import 'package:fl_toast/fl_toast.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:simplecloudnotifier/models/api_error.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';

import '../models/channel.dart';
import '../models/message.dart';

enum ChannelSelector {
  owned(apiKey: 'owned'), // Return all channels of the user
  subscribedAny(apiKey: 'subscribed_any'), // Return all channels that the user is subscribing to
  allAny(apiKey: 'all_any'), // Return channels that the user owns or is subscribing
  subscribed(apiKey: 'subscribed'), // Return all channels that the user is subscribing to (even unconfirmed)
  all(apiKey: 'all'); // Return channels that the user owns or is subscribing (even unconfirmed)

  const ChannelSelector({required this.apiKey});
  final String apiKey;
}

class APIClient {
  static const String _base = 'https://simplecloudnotifier.de/api/v2';

  static Future<T> _request<T>({
    required String name,
    required String method,
    required String relURL,
    Map<String, String>? query,
    required T Function(Map<String, dynamic> json)? fn,
    dynamic jsonBody,
    KeyTokenAuth? auth,
    Map<String, String>? header,
  }) async {
    final t0 = DateTime.now();

    final uri = Uri.parse('$_base/$relURL').replace(queryParameters: query ?? {});

    final req = http.Request(method, uri);

    if (jsonBody != null) {
      req.body = jsonEncode(jsonBody);
      req.headers['Content-Type'] = 'application/json';
    }

    if (auth != null) {
      req.headers['Authorization'] = 'SCN ${auth.token}';
    }

    req.headers['User-Agent'] = 'simplecloudnotifier/flutter/${Globals().platform.replaceAll(' ', '_')} ${Globals().version}+${Globals().buildNumber}';

    if (header != null && !header.isEmpty) {
      req.headers.addAll(header);
    }

    int responseStatusCode = 0;
    String responseBody = '';
    Map<String, String> responseHeaders = {};

    try {
      final response = await req.send();
      responseBody = await response.stream.bytesToString();
      responseStatusCode = response.statusCode;
      responseHeaders = response.headers;
    } catch (exc, trace) {
      RequestLog.addRequestException(name, t0, method, uri, req.body, req.headers, exc, trace);
      showPlatformToast(child: Text('Request "${name}" is fehlgeschlagen'), context: ToastProvider.context);
      rethrow;
    }

    if (responseStatusCode != 200) {
      try {
        final apierr = APIError.fromJson(jsonDecode(responseBody) as Map<String, dynamic>);

        RequestLog.addRequestAPIError(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders, apierr);
        showPlatformToast(child: Text('Request "${name}" is fehlgeschlagen'), context: ToastProvider.context);
        throw Exception(apierr.message);
      } catch (_) {}

      RequestLog.addRequestErrorStatuscode(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
      showPlatformToast(child: Text('Request "${name}" is fehlgeschlagen'), context: ToastProvider.context);
      throw Exception('API request failed with status code ${responseStatusCode}');
    }

    try {
      final data = jsonDecode(responseBody);

      if (fn != null) {
        final result = fn(data as Map<String, dynamic>);
        RequestLog.addRequestSuccess(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
        return result;
      } else {
        RequestLog.addRequestSuccess(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
        return null as T;
      }
    } catch (exc, trace) {
      RequestLog.addRequestDecodeError(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders, exc, trace);
      showPlatformToast(child: Text('Request "${name}" is fehlgeschlagen'), context: ToastProvider.context);
      rethrow;
    }
  }

  // ==========================================================================================================================================================

  static Future<bool> verifyToken(String uid, String tok) async {
    try {
      await _request<void>(
        name: 'verifyToken',
        method: 'GET',
        relURL: '/users/$uid',
        fn: null,
        auth: KeyTokenAuth(userId: uid, token: tok),
      );
      return true;
    } catch (e) {
      return false;
    }
  }

  static Future<User> getUser(KeyTokenAuth auth, String uid) async {
    return await _request(
      name: 'getUser',
      method: 'GET',
      relURL: 'users/$uid',
      fn: User.fromJson,
      auth: auth,
    );
  }

  static Future<List<ChannelWithSubscription>> getChannelList(KeyTokenAuth auth, ChannelSelector sel) async {
    return await _request(
      name: 'getChannelList',
      method: 'GET',
      relURL: 'users/${auth.userId}/channels',
      query: {'selector': sel.apiKey},
      fn: (json) => ChannelWithSubscription.fromJsonArray(json['channels'] as List<dynamic>),
      auth: auth,
    );
  }

  static Future<(String, List<Message>)> getMessageList(KeyTokenAuth auth, String pageToken, int? pageSize) async {
    return await _request(
      name: 'getMessageList',
      method: 'GET',
      relURL: 'messages',
      query: {'next_page_token': pageToken, if (pageSize != null) 'page_size': pageSize.toString()},
      fn: (json) => Message.fromPaginatedJsonArray(json, 'messages', 'next_page_token'),
      auth: auth,
    );
  }

  static Future<Message> getMessage(KeyTokenAuth auth, String msgid) async {
    return await _request(
      name: 'getMessage',
      method: 'GET',
      relURL: 'messages/$msgid',
      query: {},
      fn: Message.fromJson,
      auth: auth,
    );
  }
}
