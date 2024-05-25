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
      fields[0] as DateTime,
      fields[1] as SCNLogLevel,
      fields[2] as String,
      fields[3] as String,
      fields[4] as String,
    );
  }

  @override
  void write(BinaryWriter writer, SCNLog obj) {
    writer
      ..writeByte(5)
      ..writeByte(0)
      ..write(obj.timestamp)
      ..writeByte(1)
      ..write(obj.level)
      ..writeByte(2)
      ..write(obj.message)
      ..writeByte(3)
      ..write(obj.additional)
      ..writeByte(4)
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
