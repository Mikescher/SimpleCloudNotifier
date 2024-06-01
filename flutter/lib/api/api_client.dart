import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:simplecloudnotifier/models/api_error.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

enum ChannelSelector {
  owned(apiKey: 'owned'), // Return all channels of the user
  subscribedAny(apiKey: 'subscribed_any'), // Return all channels that the user is subscribing to (even unconfirmed)
  allAny(apiKey: 'all_any'), // Return channels that the user owns or is subscribing (even unconfirmed)
  subscribed(apiKey: 'subscribed'), // Return all channels that the user is subscribing to
  all(apiKey: 'all'); // Return channels that the user owns or is subscribing

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

    print('[REQUEST|RUN] [${method}] ${name}');

    if (jsonBody != null) {
      req.body = jsonEncode(jsonBody);
      req.headers['Content-Type'] = 'application/json';
    }

    if (auth != null) {
      req.headers['Authorization'] = 'SCN ${auth.tokenAdmin}';
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
      Toaster.error("Error", 'Request "${name}" failed');
      ApplicationLog.error('Request "${name}" failed: ' + exc.toString(), trace: trace);
      rethrow;
    }

    if (responseStatusCode != 200) {
      try {
        final apierr = APIError.fromJson(jsonDecode(responseBody) as Map<String, dynamic>);

        RequestLog.addRequestAPIError(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders, apierr);
        Toaster.error("Error", 'Request "${name}" failed');
        throw Exception(apierr.message);
      } catch (exc, trace) {
        ApplicationLog.warn('Failed to decode api response as error-object', additional: exc.toString() + "\nBody:\n" + responseBody, trace: trace);
      }

      RequestLog.addRequestErrorStatuscode(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
      Toaster.error("Error", 'Request "${name}" failed');
      throw Exception('API request failed with status code ${responseStatusCode}');
    }

    try {
      final data = jsonDecode(responseBody);

      if (fn != null) {
        final result = fn(data as Map<String, dynamic>);
        RequestLog.addRequestSuccess(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
        print('[REQUEST|FIN] [${method}] ${name}');
        return result;
      } else {
        RequestLog.addRequestSuccess(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders);
        print('[REQUEST|FIN] [${method}] ${name}');
        return null as T;
      }
    } catch (exc, trace) {
      RequestLog.addRequestDecodeError(name, t0, method, uri, req.body, req.headers, responseStatusCode, responseBody, responseHeaders, exc, trace);
      Toaster.error("Error", 'Request "${name}" failed');
      ApplicationLog.error('Failed to decode response: ' + exc.toString(), additional: "\nBody:\n" + responseBody, trace: trace);
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
        auth: KeyTokenAuth(userId: uid, tokenAdmin: tok, tokenSend: ''),
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

  static Future<Client> addClient(KeyTokenAuth? auth, String fcmToken, String agentModel, String agentVersion, String? descriptionName, String clientType) async {
    return await _request(
      name: 'addClient',
      method: 'POST',
      relURL: 'users/${auth!.userId}/clients',
      jsonBody: {
        'fcm_token': fcmToken,
        'agent_model': agentModel,
        'agent_version': agentVersion,
        'client_type': clientType,
        'description_name': descriptionName,
      },
      fn: Client.fromJson,
      auth: auth,
    );
  }

  static Future<Client> updateClient(KeyTokenAuth? auth, String clientID, String fcmToken, String agentModel, String? descriptionName, String agentVersion) async {
    return await _request(
      name: 'updateClient',
      method: 'PUT',
      relURL: 'users/${auth!.userId}/clients/$clientID',
      jsonBody: {
        'fcm_token': fcmToken,
        'agent_model': agentModel,
        'agent_version': agentVersion,
        'description_name': descriptionName,
      },
      fn: Client.fromJson,
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

  static Future<(String, List<Message>)> getMessageList(KeyTokenAuth auth, String pageToken, {int? pageSize, List<String>? channelIDs}) async {
    return await _request(
      name: 'getMessageList',
      method: 'GET',
      relURL: 'messages',
      query: {
        'next_page_token': pageToken,
        if (pageSize != null) 'page_size': pageSize.toString(),
        if (channelIDs != null) 'channel_id': channelIDs.join(","),
      },
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

  static Future<List<Subscription>> getSubscriptionList(KeyTokenAuth auth) async {
    return await _request(
      name: 'getSubscriptionList',
      method: 'GET',
      relURL: 'users/${auth.userId}/subscriptions',
      fn: (json) => Subscription.fromJsonArray(json['subscriptions'] as List<dynamic>),
      auth: auth,
    );
  }

  static Future<List<Client>> getClientList(KeyTokenAuth auth) async {
    return await _request(
      name: 'getClientList',
      method: 'GET',
      relURL: 'users/${auth.userId}/clients',
      fn: (json) => Client.fromJsonArray(json['clients'] as List<dynamic>),
      auth: auth,
    );
  }

  static Future<List<KeyToken>> getKeyTokenList(KeyTokenAuth auth) async {
    return await _request(
      name: 'getKeyTokenList',
      method: 'GET',
      relURL: 'users/${auth.userId}/keys',
      fn: (json) => KeyToken.fromJsonArray(json['keys'] as List<dynamic>),
      auth: auth,
    );
  }

  static Future<UserWithClientsAndKeys> createUserWithClient(String? username, String clientFcmToken, String clientAgentModel, String clientAgentVersion, String? clientDescriptionName, String clientType) async {
    return await _request(
      name: 'createUserWithClient',
      method: 'POST',
      relURL: 'users',
      jsonBody: {
        'username': username,
        'fcm_token': clientFcmToken,
        'agent_model': clientAgentModel,
        'agent_version': clientAgentVersion,
        'description_name': clientDescriptionName,
        'client_type': clientType,
        'no_client': false,
      },
      fn: UserWithClientsAndKeys.fromJson,
    );
  }
}
