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
    return Message(
      messageID: json['message_id'],
      senderUserID: json['sender_user_id'],
      channelInternalName: json['channel_internal_name'],
      channelID: json['channel_id'],
      senderName: json['sender_name'],
      senderIP: json['sender_ip'],
      timestamp: json['timestamp'],
      title: json['title'],
      content: json['content'],
      priority: json['priority'],
      userMessageID: json['usr_message_id'],
      usedKeyID: json['used_key_id'],
      trimmed: json['trimmed'],
    );
  }

  static fromPaginatedJsonArray(Map<String, dynamic> data, String keyMessages, String keyToken) {
    final npt = data[keyToken] as String;

    final messages = (data[keyMessages] as List<dynamic>).map<Message>((e) => Message.fromJson(e)).toList();

    return (npt, messages);
  }
}
