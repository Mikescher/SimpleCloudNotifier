class Client {
  final String clientID;
  final String userID;
  final String type;
  final String fcmToken;
  final String timestampCreated;
  final String agentModel;
  final String agentVersion;
  final String? name;

  const Client({
    required this.clientID,
    required this.userID,
    required this.type,
    required this.fcmToken,
    required this.timestampCreated,
    required this.agentModel,
    required this.agentVersion,
    required this.name,
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
      name: json['name'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'client_id': clientID,
      'user_id': userID,
      'type': type,
      'fcm_token': fcmToken,
      'timestamp_created': timestampCreated,
      'agent_model': agentModel,
      'agent_version': agentVersion,
      'name': name,
    };
  }

  static List<Client> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<Client>((e) => Client.fromJson(e as Map<String, dynamic>)).toList();
  }
}
