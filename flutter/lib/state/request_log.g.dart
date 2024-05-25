// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'request_log.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class SCNRequestAdapter extends TypeAdapter<SCNRequest> {
  @override
  final int typeId = 100;

  @override
  SCNRequest read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return SCNRequest(
      timestampStart: fields[0] as DateTime,
      timestampEnd: fields[1] as DateTime,
      name: fields[2] as String,
      method: fields[6] as String,
      url: fields[7] as String,
      requestHeaders: (fields[8] as Map).cast<String, String>(),
      requestBody: fields[12] as String,
      responseStatusCode: fields[9] as int,
      responseHeaders: (fields[10] as Map).cast<String, String>(),
      responseBody: fields[11] as String,
      type: fields[3] as String,
      error: fields[4] as String,
      stackTrace: fields[5] as String,
    );
  }

  @override
  void write(BinaryWriter writer, SCNRequest obj) {
    writer
      ..writeByte(13)
      ..writeByte(0)
      ..write(obj.timestampStart)
      ..writeByte(1)
      ..write(obj.timestampEnd)
      ..writeByte(2)
      ..write(obj.name)
      ..writeByte(3)
      ..write(obj.type)
      ..writeByte(4)
      ..write(obj.error)
      ..writeByte(5)
      ..write(obj.stackTrace)
      ..writeByte(6)
      ..write(obj.method)
      ..writeByte(7)
      ..write(obj.url)
      ..writeByte(8)
      ..write(obj.requestHeaders)
      ..writeByte(12)
      ..write(obj.requestBody)
      ..writeByte(9)
      ..write(obj.responseStatusCode)
      ..writeByte(10)
      ..write(obj.responseHeaders)
      ..writeByte(11)
      ..write(obj.responseBody);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SCNRequestAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
