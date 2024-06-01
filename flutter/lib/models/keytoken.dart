class KeyToken {
  final String keytokenID;
  final String name;
  final String timestampCreated;
  final String? timestampLastused;
  final String ownerUserID;
  final bool allChannels;
  final List<String> channels;
  final String permissions;
  final int messagesSent;

  const KeyToken({
    required this.keytokenID,
    required this.name,
    required this.timestampCreated,
    required this.timestampLastused,
    required this.ownerUserID,
    required this.allChannels,
    required this.channels,
    required this.permissions,
    required this.messagesSent,
  });

  factory KeyToken.fromJson(Map<String, dynamic> json) {
    return KeyToken(
      keytokenID: json['keytoken_id'] as String,
      name: json['name'] as String,
      timestampCreated: json['timestamp_created'] as String,
      timestampLastused: json['timestamp_lastused'] as String?,
      ownerUserID: json['owner_user_id'] as String,
      allChannels: json['all_channels'] as bool,
      channels: (json['channels'] as List<dynamic>).map((e) => e as String).toList(),
      permissions: json['permissions'] as String,
      messagesSent: json['messages_sent'] as int,
    );
  }

  static List<KeyToken> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<KeyToken>((e) => KeyToken.fromJson(e as Map<String, dynamic>)).toList();
  }
}
