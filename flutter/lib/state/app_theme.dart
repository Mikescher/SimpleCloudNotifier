import 'package:flutter/foundation.dart';

class AppTheme extends ChangeNotifier {
  bool _darkmode = false;
  bool get darkMode => _darkmode;

  void setDarkMode(bool v) {
    _darkmode = v;
    notifyListeners();
  }

  void switchDarkMode() {
    _darkmode = !_darkmode;
    notifyListeners();
  }
}
