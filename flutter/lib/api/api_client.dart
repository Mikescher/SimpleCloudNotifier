import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:simplecloudnotifier/api/api_exception.dart';
import 'package:simplecloudnotifier/models/api_error.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/state/token_source.dart';
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
    String? authToken,
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

    if (authToken != null) {
      req.headers['Authorization'] = 'SCN ${authToken}';
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
        throw APIException(responseStatusCode, apierr.error, apierr.errhighlight, apierr.message);
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

  static Future<User> getUser(TokenSource auth, String uid) async {
    return await _request(
      name: 'getUser',
      method: 'GET',
      relURL: 'users/$uid',
      fn: User.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<UserPreview> getUserPreview(TokenSource auth, String uid) async {
    return await _request(
      name: 'getUserPreview',
      method: 'GET',
      relURL: 'preview/users/$uid',
      fn: UserPreview.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<Client> addClient(TokenSource auth, String fcmToken, String agentModel, String agentVersion, String? name, String clientType) async {
    return await _request(
      name: 'addClient',
      method: 'POST',
      relURL: 'users/${auth.getUserID()}/clients',
      jsonBody: {
        'fcm_token': fcmToken,
        'agent_model': agentModel,
        'agent_version': agentVersion,
        'client_type': clientType,
        'name': name,
      },
      fn: Client.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<Client> updateClient(TokenSource auth, String clientID, String fcmToken, String agentModel, String? name, String agentVersion) async {
    return await _request(
      name: 'updateClient',
      method: 'PUT',
      relURL: 'users/${auth.getUserID()}/clients/$clientID',
      jsonBody: {
        'fcm_token': fcmToken,
        'agent_model': agentModel,
        'agent_version': agentVersion,
        'name': name,
      },
      fn: Client.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<Client> getClient(TokenSource auth, String cid) async {
    return await _request(
      name: 'getClient',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/clients/$cid',
      fn: Client.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<List<ChannelWithSubscription>> getChannelList(TokenSource auth, ChannelSelector sel) async {
    return await _request(
      name: 'getChannelList',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/channels',
      query: {'selector': sel.apiKey},
      fn: (json) => ChannelWithSubscription.fromJsonArray(json['channels'] as List<dynamic>),
      authToken: auth.getToken(),
    );
  }

  static Future<ChannelWithSubscription> getChannel(TokenSource auth, String cid) async {
    return await _request(
      name: 'getChannel',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/channels/${cid}',
      fn: ChannelWithSubscription.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<ChannelPreview> getChannelPreview(TokenSource auth, String cid) async {
    return await _request(
      name: 'getChannelPreview',
      method: 'GET',
      relURL: 'preview/channels/${cid}',
      fn: ChannelPreview.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<(String, List<Message>)> getMessageList(TokenSource auth, String pageToken, {int? pageSize, List<String>? channelIDs}) async {
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
      authToken: auth.getToken(),
    );
  }

  static Future<Message> getMessage(TokenSource auth, String msgid) async {
    return await _request(
      name: 'getMessage',
      method: 'GET',
      relURL: 'messages/$msgid',
      query: {},
      fn: Message.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<List<Subscription>> getSubscriptionList(TokenSource auth) async {
    return await _request(
      name: 'getSubscriptionList',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/subscriptions',
      fn: (json) => Subscription.fromJsonArray(json['subscriptions'] as List<dynamic>),
      authToken: auth.getToken(),
    );
  }

  static Future<List<Client>> getClientList(TokenSource auth) async {
    return await _request(
      name: 'getClientList',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/clients',
      fn: (json) => Client.fromJsonArray(json['clients'] as List<dynamic>),
      authToken: auth.getToken(),
    );
  }

  static Future<List<KeyToken>> getKeyTokenList(TokenSource auth) async {
    return await _request(
      name: 'getKeyTokenList',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/keys',
      fn: (json) => KeyToken.fromJsonArray(json['keys'] as List<dynamic>),
      authToken: auth.getToken(),
    );
  }

  static Future<UserWithClientsAndKeys> createUserWithClient(String? username, String clientFcmToken, String clientAgentModel, String clientAgentVersion, String? clientName, String clientType) async {
    return await _request(
      name: 'createUserWithClient',
      method: 'POST',
      relURL: 'users',
      jsonBody: {
        'username': username,
        'fcm_token': clientFcmToken,
        'agent_model': clientAgentModel,
        'agent_version': clientAgentVersion,
        'client_name': clientName,
        'client_type': clientType,
        'no_client': false,
      },
      fn: UserWithClientsAndKeys.fromJson,
    );
  }

  static Future<KeyToken> getKeyToken(TokenSource auth, String kid) async {
    return await _request(
      name: 'getKeyToken',
      method: 'GET',
      relURL: 'users/${auth.getUserID()}/keys/$kid',
      fn: KeyToken.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<KeyTokenPreview> getKeyTokenPreview(TokenSource auth, String kid) async {
    return await _request(
      name: 'getKeyTokenPreview',
      method: 'GET',
      relURL: 'preview/keys/$kid',
      fn: KeyTokenPreview.fromJson,
      authToken: auth.getToken(),
    );
  }

  static Future<KeyToken> getKeyTokenByToken(String userid, String token) async {
    return await _request(
      name: 'getCurrentKeyToken',
      method: 'GET',
      relURL: 'users/${userid}/keys/current',
      fn: KeyToken.fromJson,
      authToken: token,
    );
  }
}
