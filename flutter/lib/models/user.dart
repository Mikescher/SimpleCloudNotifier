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
    return switch (json) {
      {
        'user_id': String userID,
        'username': String? username,
        'timestamp_created': String timestampCreated,
        'timestamp_lastread': String? timestampLastRead,
        'timestamp_lastsent': String? timestampLastSent,
        'messages_sent': int messagesSent,
        'quota_used': int quotaUsed,
        'quota_remaining': int quotaRemaining,
        'quota_max': int quotaPerDay,
        'is_pro': bool isPro,
        'default_channel': String defaultChannel,
        'max_body_size': int maxBodySize,
        'max_title_length': int maxTitleLength,
        'default_priority': int defaultPriority,
        'max_channel_name_length': int maxChannelNameLength,
        'max_channel_description_length': int maxChannelDescriptionLength,
        'max_sender_name_length': int maxSenderNameLength,
        'max_user_message_id_length': int maxUserMessageIDLength,
      } =>
        User(
          userID: userID,
          username: username,
          timestampCreated: timestampCreated,
          timestampLastRead: timestampLastRead,
          timestampLastSent: timestampLastSent,
          messagesSent: messagesSent,
          quotaUsed: quotaUsed,
          quotaRemaining: quotaRemaining,
          quotaPerDay: quotaPerDay,
          isPro: isPro,
          defaultChannel: defaultChannel,
          maxBodySize: maxBodySize,
          maxTitleLength: maxTitleLength,
          defaultPriority: defaultPriority,
          maxChannelNameLength: maxChannelNameLength,
          maxChannelDescriptionLength: maxChannelDescriptionLength,
          maxSenderNameLength: maxSenderNameLength,
          maxUserMessageIDLength: maxUserMessageIDLength,
        ),
      _ => throw const FormatException('Failed to decode User.'),
    };
  }
}
