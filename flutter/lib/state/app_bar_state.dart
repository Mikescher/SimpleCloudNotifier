import 'package:flutter/foundation.dart';

class AppBarState extends ChangeNotifier {
  static AppBarState? _singleton = AppBarState._internal();

  factory AppBarState() {
    return _singleton ?? (_singleton = AppBarState._internal());
  }

  AppBarState._internal() {}

  List<void Function(String)> _searchListeners = [];

  bool _loadingIndeterminate = false;
  bool get loadingIndeterminate => _loadingIndeterminate;

  bool _showSearchField = false;
  bool get showSearchField => _showSearchField;

  void setLoadingIndeterminate(bool v) {
    if (_loadingIndeterminate == v) return;
    _loadingIndeterminate = v;
    notifyListeners();
  }

  void setShowSearchField(bool v) {
    if (_showSearchField == v) return;
    _showSearchField = v;
    notifyListeners();
  }

  void subscribeSearchListener(void Function(String) listener) {
    _searchListeners.add(listener);
  }

  void unsubscribeSearchListener(void Function(String) listener) {
    _searchListeners.remove(listener);
  }

  void notifySearchListeners(String query) {
    for (var listener in _searchListeners) {
      listener(query);
    }
  }
}
