import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/api/api_exception.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/token_source.dart';

class AppAuth extends ChangeNotifier implements TokenSource {
  String? _clientID;
  String? _userID;
  String? _tokenAdmin;
  String? _tokenSend;

  (User, DateTime)? _user;

  (Client, DateTime)? _client;

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
    _client = (client, DateTime.now());

    _user = (user, DateTime.now());

    _userID = user.userID;
    _clientID = client.clientID;

    _tokenAdmin = tokenAdmin;
    _tokenSend = tokenSend;

    notifyListeners();
  }

  void setClientAndClientID(Client client) {
    _client = (client, DateTime.now());
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

    final userjson = Globals().sharedPrefs.getString('auth.user.obj');
    final userqdate = Globals().sharedPrefs.getString('auth.user.qdate');
    final clientjson = Globals().sharedPrefs.getString('auth.client.obj');
    final clientqdate = Globals().sharedPrefs.getString('auth.client.qdate');

    if (userjson != null && userqdate != null) {
      try {
        final ts = DateTime.parse(userqdate);
        final obj = User.fromJson(jsonDecode(userjson) as Map<String, dynamic>);
        _user = (obj, ts);
      } catch (exc, trace) {
        ApplicationLog.error('failed to parse user object from shared-prefs (auth.user.obj): ' + exc.toString(), additional: 'Data:\n${userjson}\nQDate:\n${userqdate}', trace: trace);
        _user = null;
      }
    }

    if (clientjson != null && clientqdate != null) {
      try {
        final ts = DateTime.parse(clientqdate);
        final obj = Client.fromJson(jsonDecode(clientjson) as Map<String, dynamic>);
        _client = (obj, ts);
      } catch (exc, trace) {
        ApplicationLog.error('failed to parse user object from shared-prefs (auth.client.obj): ' + exc.toString(), additional: 'Data:\n${clientjson}\nQDate:\n${clientqdate}', trace: trace);
        _client = null;
      }
    }

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
      await Globals().sharedPrefs.remove('auth.user.obj');
      await Globals().sharedPrefs.remove('auth.user.qdate');
      await Globals().sharedPrefs.remove('auth.client.obj');
      await Globals().sharedPrefs.remove('auth.client.qdate');
    } else {
      await Globals().sharedPrefs.setString('auth.userid', _userID!);
      await Globals().sharedPrefs.setString('auth.clientid', _clientID!);
      await Globals().sharedPrefs.setString('auth.tokenadmin', _tokenAdmin!);
      await Globals().sharedPrefs.setString('auth.tokensend', _tokenSend!);
      if (Globals().sharedPrefs.getString('auth.cdate') == null) await Globals().sharedPrefs.setString('auth.cdate', DateTime.now().toIso8601String());
      await Globals().sharedPrefs.setString('auth.mdate', DateTime.now().toIso8601String());

      if (_user != null) {
        await Globals().sharedPrefs.setString('auth.user.obj', jsonEncode(_user!.$1.toJson()));
        await Globals().sharedPrefs.setString('auth.user.qdate', _user!.$2.toIso8601String());
      } else {
        await Globals().sharedPrefs.remove('auth.user.obj');
        await Globals().sharedPrefs.remove('auth.user.qdate');
      }

      if (_client != null) {
        await Globals().sharedPrefs.setString('auth.client.obj', jsonEncode(_client!.$1.toJson()));
        await Globals().sharedPrefs.setString('auth.client.qdate', _client!.$2.toIso8601String());
      } else {
        await Globals().sharedPrefs.remove('auth.client.obj');
        await Globals().sharedPrefs.remove('auth.client.qdate');
      }
    }

    Globals().sharedPrefs.setString('auth.mdate', DateTime.now().toIso8601String());
  }

  Future<User> loadUser({bool force = false, Duration? forceIfOlder = null}) async {
    if (forceIfOlder != null && _user != null && _user!.$2.difference(DateTime.now()) > forceIfOlder) {
      force = true;
    }

    if (!force && _user != null && _user!.$1.userID == _userID) {
      return _user!.$1;
    }

    if (_userID == null || _tokenAdmin == null) {
      throw Exception('Not authenticated');
    }

    final user = await APIClient.getUser(this, _userID!);

    _user = (user, DateTime.now());

    await save();

    return user;
  }

  User? getUserOrNull() {
    return _user?.$1;
  }

  Future<Client?> loadClient({bool force = false, Duration? forceIfOlder = null}) async {
    if (forceIfOlder != null && _client != null && _client!.$2.difference(DateTime.now()) > forceIfOlder) {
      force = true;
    }

    if (!force && _client != null && _client!.$1.clientID == _clientID) {
      return _client!.$1;
    }

    if (_clientID == null || _tokenAdmin == null) {
      throw Exception('Not authenticated');
    }

    try {
      final client = await APIClient.getClient(this, _clientID!);

      _client = (client, DateTime.now());

      await save();

      return client;
    } on APIException catch (_) {
      _client = null;
      return null;
    } catch (exc) {
      _client = null;
      rethrow;
    }
  }

  Client? getClientOrNull() {
    return _client?.$1;
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
