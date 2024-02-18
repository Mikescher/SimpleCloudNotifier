class Message {
  final String messageID;
  final String senderUserID;
  final String channelInternalName;
  final String channelID;
  final String? senderName;
  final String senderIP;
  final String timestamp;
  final String title;
  final String? content;
  final int priority;
  final String? userMessageID;
  final String usedKeyID;
  final bool trimmed;

  const Message({
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

  factory Message.fromJson(Map<String, dynamic> json) {
    return switch (json) {
      {
        'message_id': String messageID,
        'sender_user_id': String senderUserID,
        'channel_internal_name': String channelInternalName,
        'channel_id': String channelID,
        'sender_name': String? senderName,
        'sender_ip': String senderIP,
        'timestamp': String timestamp,
        'title': String title,
        'content': String? content,
        'priority': int priority,
        'usr_message_id': String? userMessageID,
        'used_key_id': String usedKeyID,
        'trimmed': bool trimmed,
      } =>
        Message(
          messageID: messageID,
          senderUserID: senderUserID,
          channelInternalName: channelInternalName,
          channelID: channelID,
          senderName: senderName,
          senderIP: senderIP,
          timestamp: timestamp,
          title: title,
          content: content,
          priority: priority,
          userMessageID: userMessageID,
          usedKeyID: usedKeyID,
          trimmed: trimmed,
        ),
      _ => throw const FormatException('Failed to decode Message.'),
    };
  }
}
