import 'package:flutter/foundation.dart';
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
    //final cdat = Globals().sharedPrefs.getString('auth.cdate');
    //final mdat = Globals().sharedPrefs.getString('auth.mdate');
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
    if (_clientID == null || _userID == null || _tokenAdmin == null || _tokenSend == null) {
      await Globals().sharedPrefs.remove('auth.userid');
      await Globals().sharedPrefs.remove('auth.clientid');
      await Globals().sharedPrefs.remove('auth.tokenadmin');
      await Globals().sharedPrefs.remove('auth.tokensend');
      await Globals().sharedPrefs.setString('auth.cdate', "");
      await Globals().sharedPrefs.setString('auth.mdate', DateTime.now().toIso8601String());
    } else {
      await Globals().sharedPrefs.setString('auth.userid', _userID!);
      await Globals().sharedPrefs.setString('auth.clientid', _clientID!);
      await Globals().sharedPrefs.setString('auth.tokenadmin', _tokenAdmin!);
      await Globals().sharedPrefs.setString('auth.tokensend', _tokenSend!);
      if (Globals().sharedPrefs.getString('auth.cdate') == null) await Globals().sharedPrefs.setString('auth.cdate', DateTime.now().toIso8601String());
      await Globals().sharedPrefs.setString('auth.mdate', DateTime.now().toIso8601String());
    }

    Globals().sharedPrefs.setString('auth.mdate', DateTime.now().toIso8601String());
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
