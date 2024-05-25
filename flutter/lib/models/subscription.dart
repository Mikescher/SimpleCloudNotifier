class Subscription {
  final String subscriptionID;
  final String subscriberUserID;
  final String channelOwnerUserID;
  final String channelID;
  final String channelInternalName;
  final String timestampCreated;
  final bool confirmed;

  const Subscription({
    required this.subscriptionID,
    required this.subscriberUserID,
    required this.channelOwnerUserID,
    required this.channelID,
    required this.channelInternalName,
    required this.timestampCreated,
    required this.confirmed,
  });

  factory Subscription.fromJson(Map<String, dynamic> json) {
    return Subscription(
      subscriptionID: json['subscription_id'],
      subscriberUserID: json['subscriber_user_id'],
      channelOwnerUserID: json['channel_owner_user_id'],
      channelID: json['channel_id'],
      channelInternalName: json['channel_internal_name'],
      timestampCreated: json['timestamp_created'],
      confirmed: json['confirmed'],
    );
  }
}
