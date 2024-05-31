import 'dart:io';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';

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
  String clientType = '';
  String deviceModel = '';

  late SharedPreferences sharedPrefs;

  Future<void> init() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();

    this.appName = packageInfo.appName;
    this.packageName = packageInfo.packageName;
    this.version = packageInfo.version;
    this.buildNumber = packageInfo.buildNumber;
    this.platform = Platform.operatingSystem;
    this.hostname = Platform.localHostname;

    if (Platform.isAndroid) {
      this.clientType = 'ANDROID';
      this.deviceModel = (await DeviceInfoPlugin().androidInfo).model;
    } else if (Platform.isIOS) {
      this.clientType = 'IOS';
      this.deviceModel = (await DeviceInfoPlugin().iosInfo).model;
    } else if (Platform.isLinux) {
      this.clientType = 'LINUX';
      this.deviceModel = (await DeviceInfoPlugin().linuxInfo).prettyName;
    } else if (Platform.isWindows) {
      this.clientType = 'WINDOWS';
      this.deviceModel = (await DeviceInfoPlugin().windowsInfo).productName;
    } else if (Platform.isMacOS) {
      this.clientType = 'MACOS';
      this.deviceModel = (await DeviceInfoPlugin().macOsInfo).model;
    } else {
      this.clientType = '?';
    }

    this.sharedPrefs = await SharedPreferences.getInstance();
  }
}
