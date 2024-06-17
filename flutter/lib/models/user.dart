import 'package:simplecloudnotifier/models/client.dart';

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

  Map<String, dynamic> toJson() {
    return {
      'user_id': userID,
      'username': username,
      'timestamp_created': timestampCreated,
      'timestamp_lastread': timestampLastRead,
      'timestamp_lastsent': timestampLastSent,
      'messages_sent': messagesSent,
      'quota_used': quotaUsed,
      'quota_remaining': quotaRemaining,
      'quota_max': quotaPerDay,
      'is_pro': isPro,
      'default_channel': defaultChannel,
      'max_body_size': maxBodySize,
      'max_title_length': maxTitleLength,
      'default_priority': defaultPriority,
      'max_channel_name_length': maxChannelNameLength,
      'max_channel_description_length': maxChannelDescriptionLength,
      'max_sender_name_length': maxSenderNameLength,
      'max_user_message_id_length': maxUserMessageIDLength,
    };
  }
}

class UserWithClientsAndKeys {
  final User user;
  final List<Client> clients;
  final String sendKey;
  final String readKey;
  final String adminKey;

  UserWithClientsAndKeys({
    required this.user,
    required this.clients,
    required this.sendKey,
    required this.readKey,
    required this.adminKey,
  });

  factory UserWithClientsAndKeys.fromJson(Map<String, dynamic> json) {
    return UserWithClientsAndKeys(
      user: User.fromJson(json),
      clients: Client.fromJsonArray(json['clients'] as List<dynamic>),
      sendKey: json['send_key'] as String,
      readKey: json['read_key'] as String,
      adminKey: json['admin_key'] as String,
    );
  }
}

class UserPreview {
  final String userID;
  final String? username;

  const UserPreview({
    required this.userID,
    required this.username,
  });

  factory UserPreview.fromJson(Map<String, dynamic> json) {
    return UserPreview(
      userID: json['user_id'] as String,
      username: json['username'] as String?,
    );
  }
}
