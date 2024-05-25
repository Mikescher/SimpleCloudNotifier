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
      channelID: json['channel_id'],
      ownerUserID: json['owner_user_id'],
      internalName: json['internal_name'],
      displayName: json['display_name'],
      descriptionName: json['description_name'],
      subscribeKey: json['subscribe_key'],
      timestampCreated: json['timestamp_created'],
      timestampLastSent: json['timestamp_lastsent'],
      messagesSent: json['messages_sent'],
    );
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
    return ChannelWithSubscription(
      channelID: json['channel_id'],
      ownerUserID: json['owner_user_id'],
      internalName: json['internal_name'],
      displayName: json['display_name'],
      descriptionName: json['description_name'],
      subscribeKey: json['subscribe_key'],
      timestampCreated: json['timestamp_created'],
      timestampLastSent: json['timestamp_lastsent'],
      messagesSent: json['messages_sent'],
      subscription: Subscription.fromJson(json['subscription']),
    );
  }

  static List<ChannelWithSubscription> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<ChannelWithSubscription>((e) => ChannelWithSubscription.fromJson(e)).toList();
  }
}
