import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';

part 'scn_message.g.dart';

@HiveType(typeId: 105)
class SCNMessage extends HiveObject implements FieldDebuggable {
  @HiveField(0)
  final String messageID;

  @HiveField(10)
  final String senderUserID;
  @HiveField(11)
  final String channelInternalName;
  @HiveField(12)
  final String channelID;
  @HiveField(13)
  final String? senderName;
  @HiveField(14)
  final String senderIP;
  @HiveField(15)
  final String timestamp;
  @HiveField(16)
  final String title;
  @HiveField(17)
  final String? content;
  @HiveField(18)
  final int priority;
  @HiveField(19)
  final String? userMessageID;
  @HiveField(20)
  final String usedKeyID;
  @HiveField(21)
  final bool trimmed;

  SCNMessage({
    required this.messageID,
    required this.senderUserID,
    required this.channelInternalName,
    required this.channelID,
    required this.senderName,
    required this.senderIP,
    required this.timestamp,
    required this.title,
    required this.content,
    required this.priority,
    required this.userMessageID,
    required this.usedKeyID,
    required this.trimmed,
  });

  factory SCNMessage.fromJson(Map<String, dynamic> json) {
    return SCNMessage(
      messageID: json['message_id'] as String,
      senderUserID: json['sender_user_id'] as String,
      channelInternalName: json['channel_internal_name'] as String,
      channelID: json['channel_id'] as String,
      senderName: json['sender_name'] as String?,
      senderIP: json['sender_ip'] as String,
      timestamp: json['timestamp'] as String,
      title: json['title'] as String,
      content: json['content'] as String?,
      priority: json['priority'] as int,
      userMessageID: json['usr_message_id'] as String?,
      usedKeyID: json['used_key_id'] as String,
      trimmed: json['trimmed'] as bool,
    );
  }

  static (String, List<SCNMessage>) fromPaginatedJsonArray(Map<String, dynamic> data, String keyMessages, String keyToken) {
    final npt = data[keyToken] as String;

    final messages = (data[keyMessages] as List<dynamic>).map<SCNMessage>((e) => SCNMessage.fromJson(e as Map<String, dynamic>)).toList();

    return (npt, messages);
  }

  @override
  String toString() {
    return 'Message[${this.messageID}]';
  }

  List<(String, String)> debugFieldList() {
    return [
      ('messageID', this.messageID),
      ('senderUserID', this.senderUserID),
      ('channelInternalName', this.channelInternalName),
      ('channelID', this.channelID),
      ('senderName', this.senderName ?? ''),
      ('senderIP', this.senderIP),
      ('timestamp', this.timestamp),
      ('title', this.title),
      ('content', this.content ?? ''),
      ('priority', '${this.priority}'),
      ('userMessageID', this.userMessageID ?? ''),
      ('usedKeyID', this.usedKeyID),
      ('trimmed', '${this.trimmed}'),
    ];
  }
}
