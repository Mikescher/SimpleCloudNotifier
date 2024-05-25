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
      messageID: json['message_id'] as String,
      senderUserID: json['sender_user_id'] as String,
      channelInternalName: json['channel_internal_name'] as String,
      channelID: json['channel_id'] as String,
      senderName: json['sender_name'] as String,
      senderIP: json['sender_ip'] as String,
      timestamp: json['timestamp'] as String,
      title: json['title'] as String,
      content: json['content'] as String,
      priority: json['priority'] as int,
      userMessageID: json['usr_message_id'] as String,
      usedKeyID: json['used_key_id'] as String,
      trimmed: json['trimmed'] as bool,
    );
  }

  static (String, List<Message>) fromPaginatedJsonArray(Map<String, dynamic> data, String keyMessages, String keyToken) {
    final npt = data[keyToken] as String;

    final messages = (data[keyMessages] as List<dynamic>).map<Message>((e) => Message.fromJson(e as Map<String, dynamic>)).toList();

    return (npt, messages);
  }
}
