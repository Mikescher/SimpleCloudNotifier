import 'package:simplecloudnotifier/models/subscription.dart';

class Channel {
  final String channelID;
  final String ownerUserID;
  final String internalName;
  final String displayName;
  final String? descriptionName;
  final String? subscribeKey;
  final String timestampCreated;
  final String? timestampLastSent;
  final int messagesSent;

  const Channel({
    required this.channelID,
    required this.ownerUserID,
    required this.internalName,
    required this.displayName,
    required this.descriptionName,
    required this.subscribeKey,
    required this.timestampCreated,
    required this.timestampLastSent,
    required this.messagesSent,
  });

  factory Channel.fromJson(Map<String, dynamic> json) {
    return Channel(
      channelID: json['channel_id'] as String,
      ownerUserID: json['owner_user_id'] as String,
      internalName: json['internal_name'] as String,
      displayName: json['display_name'] as String,
      descriptionName: json['description_name'] as String?,
      subscribeKey: json['subscribe_key'] as String?,
      timestampCreated: json['timestamp_created'] as String,
      timestampLastSent: json['timestamp_lastsent'] as String?,
      messagesSent: json['messages_sent'] as int,
    );
  }
}

class ChannelWithSubscription {
  final Channel channel;
  final Subscription subscription;

  ChannelWithSubscription({
    required this.channel,
    required this.subscription,
  });

  factory ChannelWithSubscription.fromJson(Map<String, dynamic> json) {
    return ChannelWithSubscription(
      channel: Channel.fromJson(json),
      subscription: Subscription.fromJson(json['subscription'] as Map<String, dynamic>),
    );
  }

  static List<ChannelWithSubscription> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<ChannelWithSubscription>((e) => ChannelWithSubscription.fromJson(e as Map<String, dynamic>)).toList();
  }
}

class ChannelPreview {
  final String channelID;
  final String ownerUserID;
  final String internalName;
  final String displayName;
  final String? descriptionName;

  const ChannelPreview({
    required this.channelID,
    required this.ownerUserID,
    required this.internalName,
    required this.displayName,
    required this.descriptionName,
  });

  factory ChannelPreview.fromJson(Map<String, dynamic> json) {
    return ChannelPreview(
      channelID: json['channel_id'] as String,
      ownerUserID: json['owner_user_id'] as String,
      internalName: json['internal_name'] as String,
      displayName: json['display_name'] as String,
      descriptionName: json['description_name'] as String?,
    );
  }
}
