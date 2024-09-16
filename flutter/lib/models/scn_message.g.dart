// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'scn_message.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class SCNMessageAdapter extends TypeAdapter<SCNMessage> {
  @override
  final int typeId = 105;

  @override
  SCNMessage read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return SCNMessage(
      messageID: fields[0] as String,
      senderUserID: fields[10] as String,
      channelInternalName: fields[11] as String,
      channelID: fields[12] as String,
      senderName: fields[13] as String?,
      senderIP: fields[14] as String,
      timestamp: fields[15] as String,
      title: fields[16] as String,
      content: fields[17] as String?,
      priority: fields[18] as int,
      userMessageID: fields[19] as String?,
      usedKeyID: fields[20] as String,
      trimmed: fields[21] as bool,
    );
  }

  @override
  void write(BinaryWriter writer, SCNMessage obj) {
    writer
      ..writeByte(13)
      ..writeByte(0)
      ..write(obj.messageID)
      ..writeByte(10)
      ..write(obj.senderUserID)
      ..writeByte(11)
      ..write(obj.channelInternalName)
      ..writeByte(12)
      ..write(obj.channelID)
      ..writeByte(13)
      ..write(obj.senderName)
      ..writeByte(14)
      ..write(obj.senderIP)
      ..writeByte(15)
      ..write(obj.timestamp)
      ..writeByte(16)
      ..write(obj.title)
      ..writeByte(17)
      ..write(obj.content)
      ..writeByte(18)
      ..write(obj.priority)
      ..writeByte(19)
      ..write(obj.userMessageID)
      ..writeByte(20)
      ..write(obj.usedKeyID)
      ..writeByte(21)
      ..write(obj.trimmed);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SCNMessageAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
