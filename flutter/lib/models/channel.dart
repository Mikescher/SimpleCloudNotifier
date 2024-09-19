import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';

part 'channel.g.dart';

@HiveType(typeId: 104)
class Channel extends HiveObject implements FieldDebuggable {
  @HiveField(0)
  final String channelID;

  @HiveField(10)
  final String ownerUserID;
  @HiveField(11)
  final String internalName; // = InternalName, used for sending, normalized, cannot be changed
  @HiveField(12)
  final String displayName; // = DisplayName, used for display purposes, can be changed, initially equals InternalName
  @HiveField(13)
  final String? descriptionName; // = DescriptionName, (optional), longer description text, initally nil
  @HiveField(14)
  final String? subscribeKey;
  @HiveField(15)
  final String timestampCreated;
  @HiveField(16)
  final String? timestampLastSent;
  @HiveField(17)
  final int messagesSent;

  Channel({
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

  @override
  String toString() {
    return 'Channel[${this.channelID}]';
  }

  List<(String, String)> debugFieldList() {
    return [
      ('channelID', this.channelID),
      ('ownerUserID', this.ownerUserID),
      ('internalName', this.internalName),
      ('displayName', this.displayName),
      ('descriptionName', this.descriptionName ?? ''),
      ('subscribeKey', this.subscribeKey ?? ''),
      ('timestampCreated', this.timestampCreated),
      ('timestampLastSent', this.timestampLastSent ?? ''),
      ('messagesSent', '${this.messagesSent}'),
    ];
  }

  ChannelPreview toPreview() {
    return ChannelPreview(
      channelID: this.channelID,
      ownerUserID: this.ownerUserID,
      internalName: this.internalName,
      displayName: this.displayName,
      descriptionName: this.descriptionName,
    );
  }
}

class ChannelWithSubscription {
  final Channel channel;
  final Subscription? subscription;

  ChannelWithSubscription({
    required this.channel,
    required this.subscription,
  });

  factory ChannelWithSubscription.fromJson(Map<String, dynamic> json) {
    return ChannelWithSubscription(
      channel: Channel.fromJson(json),
      subscription: json['subscription'] == null ? null : Subscription.fromJson(json['subscription'] as Map<String, dynamic>),
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
