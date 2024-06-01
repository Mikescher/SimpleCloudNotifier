import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';
import 'package:xid/xid.dart';

part 'application_log.g.dart';

class ApplicationLog {
  static void debug(String message, {String? additional, StackTrace? trace}) {
    print('[DEBUG] ${message}: ${additional ?? ''}');

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
    print('[INFO] ${message}: ${additional ?? ''}');

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
    print('[WARN] ${message}: ${additional ?? ''}');

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
    print('[ERROR] ${message}: ${additional ?? ''}');

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
    print('[FATAL] ${message}: ${additional ?? ''}');

    Hive.box<SCNLog>('scn-logs').add(SCNLog(
      id: Xid().toString(),
      timestamp: DateTime.now(),
      level: SCNLogLevel.fatal,
      message: message,
      additional: additional ?? '',
      trace: trace?.toString() ?? '',
    ));
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
