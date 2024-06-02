import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/api/api_exception.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/token_source.dart';

class AppAuth extends ChangeNotifier implements TokenSource {
  String? _clientID;
  String? _userID;
  String? _tokenAdmin;
  String? _tokenSend;

  User? _user;
  Client? _client;

  String? get userID => _userID;
  String? get tokenAdmin => _tokenAdmin;
  String? get tokenSend => _tokenSend;

  static AppAuth? _singleton = AppAuth._internal();

  factory AppAuth() {
    return _singleton ?? (_singleton = AppAuth._internal());
  }

  AppAuth._internal() {
    load();
  }

  bool isAuth() {
    return _userID != null && _tokenAdmin != null;
  }

  void set(User user, Client client, String tokenAdmin, String tokenSend) {
    _client = client;
    _user = user;
    _userID = user.userID;
    _clientID = client.clientID;
    _tokenAdmin = tokenAdmin;
    _tokenSend = tokenSend;
    notifyListeners();
  }

  void setClientAndClientID(Client client) {
    _client = client;
    _clientID = client.clientID;
    notifyListeners();
  }

  void clear() {
    _clientID = null;
    _userID = null;
    _tokenAdmin = null;
    _tokenSend = null;

    _client = null;
    _user = null;

    notifyListeners();
  }

  void load() {
    final uid = Globals().sharedPrefs.getString('auth.userid');
    final cid = Globals().sharedPrefs.getString('auth.clientid');
    final toka = Globals().sharedPrefs.getString('auth.tokenadmin');
    final toks = Globals().sharedPrefs.getString('auth.tokensend');

    if (uid == null || toka == null || toks == null || cid == null) {
      clear();
      return;
    }

    _clientID = cid;
    _userID = uid;
    _tokenAdmin = toka;
    _tokenSend = toks;

    _client = null;
    _user = null;

    notifyListeners();
  }

  Future<void> save() async {
    final prefs = await SharedPreferences.getInstance();
    if (_clientID == null || _userID == null || _tokenAdmin == null || _tokenSend == null) {
      await prefs.remove('auth.userid');
      await prefs.remove('auth.tokenadmin');
      await prefs.remove('auth.tokensend');
    } else {
      await prefs.setString('auth.userid', _userID!);
      await prefs.setString('auth.clientid', _clientID!);
      await prefs.setString('auth.tokenadmin', _tokenAdmin!);
      await prefs.setString('auth.tokensend', _tokenSend!);
    }
  }

  Future<User> loadUser({bool force = false}) async {
    if (!force && _user != null && _user!.userID == _userID) {
      return _user!;
    }

    if (_userID == null || _tokenAdmin == null) {
      throw Exception('Not authenticated');
    }

    final user = await APIClient.getUser(this, _userID!);

    _user = user;
    notifyListeners();

    await save();

    return user;
  }

  Future<Client?> loadClient({bool force = false}) async {
    if (!force && _client != null && _client!.clientID == _clientID) {
      return _client!;
    }

    if (_clientID == null || _tokenAdmin == null) {
      throw Exception('Not authenticated');
    }

    try {
      final client = await APIClient.getClient(this, _clientID!);

      _client = client;
      notifyListeners();

      await save();

      return client;
    } on APIException catch (_) {
      _client = null;
      notifyListeners();
      return null;
    } catch (exc) {
      _client = null;
      rethrow;
    }
  }

  @override
  String getToken() {
    return _tokenAdmin!;
  }

  @override
  String getUserID() {
    return _userID!;
  }
}
