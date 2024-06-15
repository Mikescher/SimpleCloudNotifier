// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'fb_message.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class FBMessageAdapter extends TypeAdapter<FBMessage> {
  @override
  final int typeId = 106;

  @override
  FBMessage read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return FBMessage(
      senderId: fields[0] as String?,
      category: fields[1] as String?,
      collapseKey: fields[2] as String?,
      contentAvailable: fields[3] as bool,
      data: (fields[4] as Map).cast<String, String>(),
      from: fields[5] as String?,
      messageId: fields[6] as String?,
      messageType: fields[7] as String?,
      mutableContent: fields[8] as bool,
      notification: fields[9] as RemoteNotification?,
      sentTime: fields[10] as DateTime?,
      threadId: fields[11] as String?,
      ttl: fields[12] as int?,
      notificationAndroidChannelId: fields[20] as String?,
      notificationAndroidClickAction: fields[21] as String?,
      notificationAndroidColor: fields[22] as String?,
      notificationAndroidCount: fields[23] as int?,
      notificationAndroidImageUrl: fields[24] as String?,
      notificationAndroidLink: fields[25] as String?,
      notificationAndroidPriority: fields[26] as AndroidNotificationPriority?,
      notificationAndroidSmallIcon: fields[27] as String?,
      notificationAndroidSound: fields[28] as String?,
      notificationAndroidTicker: fields[29] as String?,
      notificationAndroidVisibility:
          fields[30] as AndroidNotificationVisibility?,
      notificationAndroidTag: fields[31] as String?,
      notificationAppleBadge: fields[40] as String?,
      notificationAppleSound: fields[41] as AppleNotificationSound?,
      notificationAppleImageUrl: fields[42] as String?,
      notificationAppleSubtitle: fields[43] as String?,
      notificationAppleSubtitleLocArgs: (fields[44] as List?)?.cast<String>(),
      notificationAppleSubtitleLocKey: fields[45] as String?,
      notificationWebAnalyticsLabel: fields[50] as String?,
      notificationWebImage: fields[51] as String?,
      notificationWebLink: fields[52] as String?,
      notificationTitle: fields[60] as String?,
      notificationTitleLocArgs: (fields[61] as List?)?.cast<String>(),
      notificationTitleLocKey: fields[62] as String?,
      notificationBody: fields[63] as String?,
      notificationBodyLocArgs: (fields[64] as List?)?.cast<String>(),
      notificationBodyLocKey: fields[65] as String?,
    );
  }

  @override
  void write(BinaryWriter writer, FBMessage obj) {
    writer
      ..writeByte(40)
      ..writeByte(0)
      ..write(obj.senderId)
      ..writeByte(1)
      ..write(obj.category)
      ..writeByte(2)
      ..write(obj.collapseKey)
      ..writeByte(3)
      ..write(obj.contentAvailable)
      ..writeByte(4)
      ..write(obj.data)
      ..writeByte(5)
      ..write(obj.from)
      ..writeByte(6)
      ..write(obj.messageId)
      ..writeByte(7)
      ..write(obj.messageType)
      ..writeByte(8)
      ..write(obj.mutableContent)
      ..writeByte(9)
      ..write(obj.notification)
      ..writeByte(10)
      ..write(obj.sentTime)
      ..writeByte(11)
      ..write(obj.threadId)
      ..writeByte(12)
      ..write(obj.ttl)
      ..writeByte(20)
      ..write(obj.notificationAndroidChannelId)
      ..writeByte(21)
      ..write(obj.notificationAndroidClickAction)
      ..writeByte(22)
      ..write(obj.notificationAndroidColor)
      ..writeByte(23)
      ..write(obj.notificationAndroidCount)
      ..writeByte(24)
      ..write(obj.notificationAndroidImageUrl)
      ..writeByte(25)
      ..write(obj.notificationAndroidLink)
      ..writeByte(26)
      ..write(obj.notificationAndroidPriority)
      ..writeByte(27)
      ..write(obj.notificationAndroidSmallIcon)
      ..writeByte(28)
      ..write(obj.notificationAndroidSound)
      ..writeByte(29)
      ..write(obj.notificationAndroidTicker)
      ..writeByte(30)
      ..write(obj.notificationAndroidVisibility)
      ..writeByte(31)
      ..write(obj.notificationAndroidTag)
      ..writeByte(40)
      ..write(obj.notificationAppleBadge)
      ..writeByte(41)
      ..write(obj.notificationAppleSound)
      ..writeByte(42)
      ..write(obj.notificationAppleImageUrl)
      ..writeByte(43)
      ..write(obj.notificationAppleSubtitle)
      ..writeByte(44)
      ..write(obj.notificationAppleSubtitleLocArgs)
      ..writeByte(45)
      ..write(obj.notificationAppleSubtitleLocKey)
      ..writeByte(50)
      ..write(obj.notificationWebAnalyticsLabel)
      ..writeByte(51)
      ..write(obj.notificationWebImage)
      ..writeByte(52)
      ..write(obj.notificationWebLink)
      ..writeByte(60)
      ..write(obj.notificationTitle)
      ..writeByte(61)
      ..write(obj.notificationTitleLocArgs)
      ..writeByte(62)
      ..write(obj.notificationTitleLocKey)
      ..writeByte(63)
      ..write(obj.notificationBody)
      ..writeByte(64)
      ..write(obj.notificationBodyLocArgs)
      ..writeByte(65)
      ..write(obj.notificationBodyLocKey);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is FBMessageAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
