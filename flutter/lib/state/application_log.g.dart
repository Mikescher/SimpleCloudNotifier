// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'application_log.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class SCNLogAdapter extends TypeAdapter<SCNLog> {
  @override
  final int typeId = 101;

  @override
  SCNLog read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return SCNLog(
      id: fields[0] as String,
      timestamp: fields[10] as DateTime,
      level: fields[11] as SCNLogLevel,
      message: fields[12] as String,
      additional: fields[13] as String,
      trace: fields[14] as String,
    );
  }

  @override
  void write(BinaryWriter writer, SCNLog obj) {
    writer
      ..writeByte(6)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(10)
      ..write(obj.timestamp)
      ..writeByte(11)
      ..write(obj.level)
      ..writeByte(12)
      ..write(obj.message)
      ..writeByte(13)
      ..write(obj.additional)
      ..writeByte(14)
      ..write(obj.trace);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SCNLogAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}

class SCNLogLevelAdapter extends TypeAdapter<SCNLogLevel> {
  @override
  final int typeId = 103;

  @override
  SCNLogLevel read(BinaryReader reader) {
    switch (reader.readByte()) {
      case 0:
        return SCNLogLevel.debug;
      case 1:
        return SCNLogLevel.info;
      case 2:
        return SCNLogLevel.warning;
      case 3:
        return SCNLogLevel.error;
      case 4:
        return SCNLogLevel.fatal;
      default:
        return SCNLogLevel.debug;
    }
  }

  @override
  void write(BinaryWriter writer, SCNLogLevel obj) {
    switch (obj) {
      case SCNLogLevel.debug:
        writer.writeByte(0);
        break;
      case SCNLogLevel.info:
        writer.writeByte(1);
        break;
      case SCNLogLevel.warning:
        writer.writeByte(2);
        break;
      case SCNLogLevel.error:
        writer.writeByte(3);
        break;
      case SCNLogLevel.fatal:
        writer.writeByte(4);
        break;
    }
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SCNLogLevelAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
