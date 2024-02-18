import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:simplecloudnotifier/models/user.dart';

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
}
