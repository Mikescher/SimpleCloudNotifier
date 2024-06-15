// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'channel.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class ChannelAdapter extends TypeAdapter<Channel> {
  @override
  final int typeId = 104;

  @override
  Channel read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Channel(
      channelID: fields[0] as String,
      ownerUserID: fields[10] as String,
      internalName: fields[11] as String,
      displayName: fields[12] as String,
      descriptionName: fields[13] as String?,
      subscribeKey: fields[14] as String?,
      timestampCreated: fields[15] as String,
      timestampLastSent: fields[16] as String?,
      messagesSent: fields[17] as int,
    );
  }

  @override
  void write(BinaryWriter writer, Channel obj) {
    writer
      ..writeByte(9)
      ..writeByte(0)
      ..write(obj.channelID)
      ..writeByte(10)
      ..write(obj.ownerUserID)
      ..writeByte(11)
      ..write(obj.internalName)
      ..writeByte(12)
      ..write(obj.displayName)
      ..writeByte(13)
      ..write(obj.descriptionName)
      ..writeByte(14)
      ..write(obj.subscribeKey)
      ..writeByte(15)
      ..write(obj.timestampCreated)
      ..writeByte(16)
      ..write(obj.timestampLastSent)
      ..writeByte(17)
      ..write(obj.messagesSent);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ChannelAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
