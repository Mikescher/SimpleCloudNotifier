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
    return switch (json) {
      {
        'channel_id': String channelID,
        'owner_user_id': String ownerUserID,
        'internal_name': String internalName,
        'display_name': String displayName,
        'description_name': String? descriptionName,
        'subscribe_key': String? subscribeKey,
        'timestamp_created': String timestampCreated,
        'timestamp_lastsent': String? timestampLastSent,
        'messages_sent': int messagesSent,
      } =>
        Channel(
          channelID: channelID,
          ownerUserID: ownerUserID,
          internalName: internalName,
          displayName: displayName,
          descriptionName: descriptionName,
          subscribeKey: subscribeKey,
          timestampCreated: timestampCreated,
          timestampLastSent: timestampLastSent,
          messagesSent: messagesSent,
        ),
      _ => throw const FormatException('Failed to decode Channel.'),
    };
  }
}

class ChannelWithSubscription extends Channel {
  final Subscription subscription;

  ChannelWithSubscription({
    required super.channelID,
    required super.ownerUserID,
    required super.internalName,
    required super.displayName,
    required super.descriptionName,
    required super.subscribeKey,
    required super.timestampCreated,
    required super.timestampLastSent,
    required super.messagesSent,
    required this.subscription,
  });

  factory ChannelWithSubscription.fromJson(Map<String, dynamic> json) {
    return switch (json) {
      {
        'channel_id': String channelID,
        'owner_user_id': String ownerUserID,
        'internal_name': String internalName,
        'display_name': String displayName,
        'description_name': String? descriptionName,
        'subscribe_key': String? subscribeKey,
        'timestamp_created': String timestampCreated,
        'timestamp_lastsent': String? timestampLastSent,
        'messages_sent': int messagesSent,
        'subscription': dynamic subscription,
      } =>
        ChannelWithSubscription(
          channelID: channelID,
          ownerUserID: ownerUserID,
          internalName: internalName,
          displayName: displayName,
          descriptionName: descriptionName,
          subscribeKey: subscribeKey,
          timestampCreated: timestampCreated,
          timestampLastSent: timestampLastSent,
          messagesSent: messagesSent,
          subscription: Subscription.fromJson(subscription),
        ),
      _ => throw const FormatException('Failed to decode Channel.'),
    };
  }
}
