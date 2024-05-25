import 'package:hive_flutter/hive_flutter.dart';

part 'application_log.g.dart';

class ApplicationLog {}

enum SCNLogLevel { debug, info, warning, error, fatal }

@HiveType(typeId: 101)
class SCNLog extends HiveObject {
  @HiveField(0)
  final DateTime timestamp;
  @HiveField(1)
  final SCNLogLevel level;
  @HiveField(2)
  final String message;
  @HiveField(3)
  final String additional;
  @HiveField(4)
  final String trace;

  SCNLog(
    this.timestamp,
    this.level,
    this.message,
    this.additional,
    this.trace,
  );
}
