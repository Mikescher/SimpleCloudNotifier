class SenderNameStatistics {
  final String name;
  final String ts_last;
  final String ts_first;
  final int count;

  const SenderNameStatistics({
    required this.name,
    required this.ts_last,
    required this.ts_first,
    required this.count,
  });

  factory SenderNameStatistics.fromJson(Map<String, dynamic> json) {
    return SenderNameStatistics(
      name: json['name'] as String,
      ts_last: json['ts_last'] as String,
      ts_first: json['ts_first'] as String,
      count: json['count'] as int,
    );
  }

  static List<SenderNameStatistics> fromJsonArray(List<dynamic> jsonArr) {
    return jsonArr.map<SenderNameStatistics>((e) => SenderNameStatistics.fromJson(e as Map<String, dynamic>)).toList();
  }
}
