class User {
  final String userID;
  final String? username;
  final String timestampCreated;
  final String? timestampLastRead;
  final String? timestampLastSent;
  final int messagesSent;
  final int quotaUsed;
  final int quotaRemaining;
  final int quotaPerDay;
  final bool isPro;
  final String defaultChannel;
  final int maxBodySize;
  final int maxTitleLength;
  final int defaultPriority;
  final int maxChannelNameLength;
  final int maxChannelDescriptionLength;
  final int maxSenderNameLength;
  final int maxUserMessageIDLength;

  const User({
    required this.userID,
    required this.username,
    required this.timestampCreated,
    required this.timestampLastRead,
    required this.timestampLastSent,
    required this.messagesSent,
    required this.quotaUsed,
    required this.quotaRemaining,
    required this.quotaPerDay,
    required this.isPro,
    required this.defaultChannel,
    required this.maxBodySize,
    required this.maxTitleLength,
    required this.defaultPriority,
    required this.maxChannelNameLength,
    required this.maxChannelDescriptionLength,
    required this.maxSenderNameLength,
    required this.maxUserMessageIDLength,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      userID: json['user_id'],
      username: json['username'],
      timestampCreated: json['timestamp_created'],
      timestampLastRead: json['timestamp_lastread'],
      timestampLastSent: json['timestamp_lastsent'],
      messagesSent: json['messages_sent'],
      quotaUsed: json['quota_used'],
      quotaRemaining: json['quota_remaining'],
      quotaPerDay: json['quota_max'],
      isPro: json['is_pro'],
      defaultChannel: json['default_channel'],
      maxBodySize: json['max_body_size'],
      maxTitleLength: json['max_title_length'],
      defaultPriority: json['default_priority'],
      maxChannelNameLength: json['max_channel_name_length'],
      maxChannelDescriptionLength: json['max_channel_description_length'],
      maxSenderNameLength: json['max_sender_name_length'],
      maxUserMessageIDLength: json['max_user_message_id_length'],
    );
  }
}
