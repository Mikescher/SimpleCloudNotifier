import 'package:flutter/material.dart';

class AppSettings extends ChangeNotifier {
  bool groupNotifications = true;
  int messagePageSize = 128;
  bool showDebugButton = true;

  static AppSettings? _singleton = AppSettings._internal();

  factory AppSettings() {
    return _singleton ?? (_singleton = AppSettings._internal());
  }

  AppSettings._internal() {
    load();
  }

  void clear() {
    //TODO

    notifyListeners();
  }

  void load() {
    //TODO

    notifyListeners();
  }

  Future<void> save() async {
    //TODO
  }
}
