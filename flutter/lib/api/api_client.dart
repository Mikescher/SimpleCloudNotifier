import 'package:http/http.dart' as http;

class APIClient {
  static const String _base = 'https://simplecloudnotifier.de/api/v2';

  static Future<bool> verifyToken(String uid, String tok) async {
    final uri = Uri.parse('$_base/users/$uid');
    final response = await http.get(uri, headers: {'Authorization': 'SCN $tok'});

    return (response.statusCode == 200);
  }
}
