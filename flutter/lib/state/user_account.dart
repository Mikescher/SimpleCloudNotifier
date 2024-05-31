import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/globals.dart';

class UserAccount extends ChangeNotifier {
  User? _user;
  User? get user => _user;

  Client? _client;
  Client? get client => _client;

  KeyTokenAuth? _auth;
  KeyTokenAuth? get auth => _auth;

  static UserAccount? _singleton = UserAccount._internal();

  factory UserAccount() {
    return _singleton ?? (_singleton = UserAccount._internal());
  }

  UserAccount._internal() {
    load();
  }

  void setToken(KeyTokenAuth auth) {
    _auth = auth;
    _user = null;
    notifyListeners();
  }

  void clearToken() {
    _auth = null;
    _user = null;
    notifyListeners();
  }

  void setUser(User user) {
    _user = user;
    notifyListeners();
  }

  void clearUser() {
    _user = null;
    notifyListeners();
  }

  void setClient(Client client) {
    _client = client;
    notifyListeners();
  }

  void clearClient() {
    _client = null;
    notifyListeners();
  }

  void load() {
    final uid = Globals().sharedPrefs.getString('auth.userid');
    final tok = Globals().sharedPrefs.getString('auth.token');

    if (uid != null && tok != null) {
      setToken(KeyTokenAuth(userId: uid, token: tok));
    } else {
      clearToken();
    }
  }

  Future<void> save() async {
    final prefs = await SharedPreferences.getInstance();
    if (_auth == null) {
      await prefs.remove('auth.userid');
      await prefs.remove('auth.token');
    } else {
      await prefs.setString('auth.userid', _auth!.userId);
      await prefs.setString('auth.token', _auth!.token);
    }
  }

  Future<User> loadUser(bool force) async {
    if (!force && _user != null) {
      return _user!;
    }

    if (_auth == null) {
      throw Exception('Not authenticated');
    }

    final user = await APIClient.getUser(_auth!, _auth!.userId);

    setUser(user);

    await save();

    return user;
  }
}
