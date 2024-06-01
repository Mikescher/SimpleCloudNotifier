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
      id: fields[0] as String,
      timestampStart: fields[10] as DateTime,
      timestampEnd: fields[11] as DateTime,
      name: fields[12] as String,
      method: fields[21] as String,
      url: fields[22] as String,
      requestHeaders: (fields[23] as Map).cast<String, String>(),
      requestBody: fields[24] as String,
      responseStatusCode: fields[31] as int,
      responseHeaders: (fields[32] as Map).cast<String, String>(),
      responseBody: fields[33] as String,
      type: fields[13] as String,
      error: fields[14] as String,
      stackTrace: fields[15] as String,
    );
  }

  @override
  void write(BinaryWriter writer, SCNRequest obj) {
    writer
      ..writeByte(14)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(10)
      ..write(obj.timestampStart)
      ..writeByte(11)
      ..write(obj.timestampEnd)
      ..writeByte(12)
      ..write(obj.name)
      ..writeByte(13)
      ..write(obj.type)
      ..writeByte(14)
      ..write(obj.error)
      ..writeByte(15)
      ..write(obj.stackTrace)
      ..writeByte(21)
      ..write(obj.method)
      ..writeByte(22)
      ..write(obj.url)
      ..writeByte(23)
      ..write(obj.requestHeaders)
      ..writeByte(24)
      ..write(obj.requestBody)
      ..writeByte(31)
      ..write(obj.responseStatusCode)
      ..writeByte(32)
      ..write(obj.responseHeaders)
      ..writeByte(33)
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
