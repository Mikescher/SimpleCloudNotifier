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
    return switch (json) {
      {
        'subscription_id': String subscriptionID,
        'subscriber_user_id': String subscriberUserID,
        'channel_owner_user_id': String channelOwnerUserID,
        'channel_id': String channelID,
        'channel_internal_name': String channelInternalName,
        'timestamp_created': String timestampCreated,
        'confirmed': bool confirmed,
      } =>
        Subscription(
          subscriptionID: subscriptionID,
          subscriberUserID: subscriberUserID,
          channelOwnerUserID: channelOwnerUserID,
          channelID: channelID,
          channelInternalName: channelInternalName,
          timestampCreated: timestampCreated,
          confirmed: confirmed,
        ),
      _ => throw const FormatException('Failed to decode Subscription.'),
    };
  }
}
