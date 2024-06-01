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
      subscriptionID: json['subscription_id'] as String,
      subscriberUserID: json['subscriber_user_id'] as String,
      channelOwnerUserID: json['channel_owner_user_id'] as String,
      channelID: json['channel_id'] as String,
      channelInternalName: json['channel_internal_name'] as String,
      timestampCreated: json['timestamp_created'] as String,
      confirmed: json['confirmed'] as bool,
    );
  }

  static List<Subscription> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<Subscription>((e) => Subscription.fromJson(e as Map<String, dynamic>)).toList();
  }
}
