import 'package:flutter/foundation.dart';

class AppBarState extends ChangeNotifier {
  static AppBarState? _singleton = AppBarState._internal();

  factory AppBarState() {
    return _singleton ?? (_singleton = AppBarState._internal());
  }

  AppBarState._internal() {}

  bool _loadingIndeterminate = false;
  bool get loadingIndeterminate => _loadingIndeterminate;

  void setLoadingIndeterminate(bool v) {
    if (_loadingIndeterminate == v) return;
    _loadingIndeterminate = v;
    notifyListeners();
  }
}
