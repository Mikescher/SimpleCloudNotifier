import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/api/api_client.dart';

import '../models/key_token_auth.dart';
import '../models/user.dart';

class UserAccount extends ChangeNotifier {
  User? _user;
  User? get user => _user;

  KeyTokenAuth? _auth;
  KeyTokenAuth? get auth => _auth;

  UserAccount() {
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

  load() async {
    final prefs = await SharedPreferences.getInstance();

    final uid = prefs.getString('auth.userid');
    final tok = prefs.getString('auth.token');

    if (uid != null && tok != null) {
      setToken(KeyTokenAuth(userId: uid, token: tok));
    } else {
      clearToken();
    }
  }

  save() async {
    final prefs = await SharedPreferences.getInstance();
    if (_auth == null) {
      await prefs.remove('auth.userid');
      await prefs.remove('auth.token');
    } else {
      await prefs.setString('auth.userid', _auth!.userId);
      await prefs.setString('auth.token', _auth!.token);
    }
  }

  loadUser(bool force) async {
    if (!force && _user != null) {
      return _user;
    }

    if (_auth == null) {
      throw Exception('Not authenticated');
    }

    final user = await APIClient.getUser(_auth!.userId, _auth!.token);

    setUser(user);

    await save();

    return user;
  }
}
