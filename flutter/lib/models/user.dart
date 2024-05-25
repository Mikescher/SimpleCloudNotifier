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
      userID: json['user_id'] as String,
      username: json['username'] as String?,
      timestampCreated: json['timestamp_created'] as String,
      timestampLastRead: json['timestamp_lastread'] as String?,
      timestampLastSent: json['timestamp_lastsent'] as String?,
      messagesSent: json['messages_sent'] as int,
      quotaUsed: json['quota_used'] as int,
      quotaRemaining: json['quota_remaining'] as int,
      quotaPerDay: json['quota_max'] as int,
      isPro: json['is_pro'] as bool,
      defaultChannel: json['default_channel'] as String,
      maxBodySize: json['max_body_size'] as int,
      maxTitleLength: json['max_title_length'] as int,
      defaultPriority: json['default_priority'] as int,
      maxChannelNameLength: json['max_channel_name_length'] as int,
      maxChannelDescriptionLength: json['max_channel_description_length'] as int,
      maxSenderNameLength: json['max_sender_name_length'] as int,
      maxUserMessageIDLength: json['max_user_message_id_length'] as int,
    );
  }
}
