class Client {
  final String clientID;
  final String userID;
  final String type;
  final String fcmToken;
  final String timestampCreated;
  final String agentModel;
  final String agentVersion;

  const Client({
    required this.clientID,
    required this.userID,
    required this.type,
    required this.fcmToken,
    required this.timestampCreated,
    required this.agentModel,
    required this.agentVersion,
  });

  factory Client.fromJson(Map<String, dynamic> json) {
    return Client(
      clientID: json['client_id'] as String,
      userID: json['user_id'] as String,
      type: json['type'] as String,
      fcmToken: json['fcm_token'] as String,
      timestampCreated: json['timestamp_created'] as String,
      agentModel: json['agent_model'] as String,
      agentVersion: json['agent_version'] as String,
    );
  }

  static List<Client> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<Client>((e) => Client.fromJson(e as Map<String, dynamic>)).toList();
  }
}
