import 'dart:convert';
import 'dart:io';

import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';
import 'package:xid/xid.dart';
import 'package:path/path.dart' as path;

part 'application_log.g.dart';

class ApplicationLog {
  //TODO max size, auto clear old

  static void debug(String message, {String? additional, StackTrace? trace}) {
    (additional != null && additional != '') ? print('[DEBUG] ${message}: ${additional}') : print('[DEBUG] ${message}');

    if (!Hive.isBoxOpen('scn-logs')) return;
    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.debug,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
  }

  static void info(String message, {String? additional, StackTrace? trace}) {
    (additional != null && additional != '') ? print('[INFO] ${message}: ${additional}') : print('[INFO] ${message}');

    if (!Hive.isBoxOpen('scn-logs')) return;
    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.info,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
  }

  static void warn(String message, {String? additional, StackTrace? trace}) {
    (additional != null && additional != '') ? print('[WARN] ${message}: ${additional}') : print('[WARN] ${message}');

    if (!Hive.isBoxOpen('scn-logs')) return;
    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.warning,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
  }

  static void error(String message, {String? additional, StackTrace? trace}) {
    (additional != null && additional != '') ? print('[ERROR] ${message}: ${additional}') : print('[ERROR] ${message}');

    if (!Hive.isBoxOpen('scn-logs')) return;
    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.error,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
  }

  static void fatal(String message, {String? additional, StackTrace? trace}) {
    (additional != null && additional != '') ? print('[FATAL] ${message}: ${additional}') : print('[FATAL] ${message}');

    if (!Hive.isBoxOpen('scn-logs')) return;
    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.fatal,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
  }

  static void writeRawFailure(String message, Map<String, dynamic> extraData) async {
    try {
      await Globals().init();

      final fn = path.join(Globals().rawFailureLogsDir.path, 'failure-${DateTime.now().toIso8601String()}.log');

      var txt = "[TEXT]\n${message}\n\n";
      for (var k in extraData.keys) {
        txt += "[${k}]\n${_debugToStr(extraData[k])}\n\n";
      }

      await File(fn).writeAsString(txt);

      ApplicationLog.debug("Wrote raw failure log to '${fn}'  ('${message}')");
    } catch (e) {
      print("Failed to <writeRawFailure>: ${e}");
    }
  }

  static String _debugToStr(dynamic v) {
    if (v is String) {
      return v;
    }
    if (v is StackTrace) {
      return v.toString();
    }

    try {
      final enc = new JsonEncoder.withIndent("  ");
      return enc.convert(v);
    } catch (e) {
      // ignore
    }

    try {
      return jsonEncode(v);
    } catch (e) {
      // ignore
    }

    try {
      return v.toString();
    } catch (e) {
      // ignore
    }

    try {
      return "${v}";
    } catch (e) {
      // ignore
    }

    return "<[!]FAILED_TO_PRINT_OBJECT>";
  }
}

@HiveType(typeId: 103)
enum SCNLogLevel {
  @HiveField(0)
  debug,
  @HiveField(1)
  info,
  @HiveField(2)
  warning,
  @HiveField(3)
  error,
  @HiveField(4)
  fatal
}

@HiveType(typeId: 101)
class SCNLog extends HiveObject implements FieldDebuggable {
  @HiveField(0)
  final String id;

  @HiveField(10)
  final DateTime timestamp;
  @HiveField(11)
  final SCNLogLevel level;
  @HiveField(12)
  final String message;
  @HiveField(13)
  final String additional;
  @HiveField(14)
  final String trace;

  SCNLog({
    required this.id,
    required this.timestamp,
    required this.level,
    required this.message,
    required this.additional,
    required this.trace,
  });

  @override
  String toString() {
    return 'SCNLog[${this.id}]';
  }

  List<(String, String)> debugFieldList() {
    return [
      ('id', this.id),
      ('timestamp', this.timestamp.toIso8601String()),
      ('level', this.level.name),
      ('message', this.message),
      ('additional', this.additional),
      ('trace', this.trace),
    ];
  }
}
