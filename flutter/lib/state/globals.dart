import 'dart:io';

import 'package:package_info_plus/package_info_plus.dart';

class Globals {
  static final Globals _singleton = Globals._internal();

  factory Globals() {
    return _singleton;
  }

  Globals._internal();

  String appName = '';
  String packageName = '';
  String version = '';
  String buildNumber = '';
  String platform = '';
  String hostname = '';

  Future<void> init() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();

    this.appName = packageInfo.appName;
    this.packageName = packageInfo.packageName;
    this.version = packageInfo.version;
    this.buildNumber = packageInfo.buildNumber;
    this.platform = Platform.operatingSystem;
    this.hostname = Platform.localHostname;
  }
}
