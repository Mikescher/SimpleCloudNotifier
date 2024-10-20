class SenderNameStatistics {
  final String name;
  final String lastTimestamp;
  final String firstTimestamp;
  final int count;

  const SenderNameStatistics({
    required this.name,
    required this.lastTimestamp,
    required this.firstTimestamp,
    required this.count,
  });

  factory SenderNameStatistics.fromJson(Map<String, dynamic> json) {
    return SenderNameStatistics(
      name: json['name'] as String,
      lastTimestamp: json['last_timestamp'] as String,
      firstTimestamp: json['first_timestamp'] as String,
      count: json['count'] as int,
    );
  }

  static List<SenderNameStatistics> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<SenderNameStatistics>((e) => SenderNameStatistics.fromJson(e as Map<String, dynamic>)).toList();
  }
}
