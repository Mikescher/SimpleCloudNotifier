import 'package:flutter/foundation.dart';

import '../models/key_token_auth.dart';
import '../models/user.dart';

class UserAccount extends ChangeNotifier {
  User? _user;
  User? get user => _user;

  KeyTokenAuth? _auth;
  KeyTokenAuth? get auth => _auth;

  void setToken(KeyTokenAuth auth) {
    _auth = auth;
    _user = user;
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
}
