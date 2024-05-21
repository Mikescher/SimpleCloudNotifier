import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/models/user.dart';

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

  static Future<bool> verifyToken(String uid, String tok) async {
    final uri = Uri.parse('$_base/users/$uid');
    final response = await http.get(uri, headers: {'Authorization': 'SCN $tok'});

    return (response.statusCode == 200);
  }

  static Future<User> getUser(String uid, String tok) async {
    final uri = Uri.parse('$_base/users/$uid');
    final response = await http.get(uri, headers: {'Authorization': 'SCN $tok'});

    if (response.statusCode != 200) {
      throw Exception('API request failed');
    }

    return User.fromJson(jsonDecode(response.body));
  }

  static Future<List<ChannelWithSubscription>> getChannelList(KeyTokenAuth auth, ChannelSelector sel) async {
    var url = '$_base/users/${auth.userId}/channels?selector=${sel.apiKey}';
    final uri = Uri.parse(url);

    final response = await http.get(uri, headers: {'Authorization': 'SCN ${auth.token}'});

    if (response.statusCode != 200) {
      throw Exception('API request failed');
    }

    final data = jsonDecode(response.body);

    return data['channels'].map<ChannelWithSubscription>((e) => ChannelWithSubscription.fromJson(e)).toList() as List<ChannelWithSubscription>;
  }

  static Future<(String, List<Message>)> getMessageList(KeyTokenAuth auth, String pageToken, int? pageSize) async {
    var url = '$_base/messages?next_page_token=$pageToken';
    if (pageSize != null) {
      url += '&page_size=$pageSize';
    }
    final uri = Uri.parse(url);

    final response = await http.get(uri, headers: {'Authorization': 'SCN ${auth.token}'});

    if (response.statusCode != 200) {
      throw Exception('API request failed');
    }

    final data = jsonDecode(response.body);

    final npt = data['next_page_token'] as String;

    final messages = data['messages'].map<Message>((e) => Message.fromJson(e)).toList() as List<Message>;

    return (npt, messages);
  }

  static Future<Message> getMessage(KeyTokenAuth auth, String msgid) async {
    final uri = Uri.parse('$_base/messages/$msgid');
    final response = await http.get(uri, headers: {'Authorization': 'SCN ${auth.token}'});

    if (response.statusCode != 200) {
      throw Exception('API request failed');
    }

    return Message.fromJson(jsonDecode(response.body));
  }
}
