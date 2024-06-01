class APIError {
  final String success;
  final String error;
  final String errhighlight;
  final String message;

  const APIError({
    required this.success,
    required this.error,
    required this.errhighlight,
    required this.message,
  });

  factory APIError.fromJson(Map<String, dynamic> json) {
    return APIError(
      success: json['success'] as String,
      error: json['error'] as String,
      errhighlight: json['errhighlight'] as String,
      message: json['message'] as String,
    );
  }
}
