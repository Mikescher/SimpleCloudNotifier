class APIException implements Exception {
  final int httpStatus;
  final String error;
  final String errHighlight;
  final String message;

  APIException(this.httpStatus, this.error, this.errHighlight, this.message);

  @override
  String toString() {
    return '[$error] $message';
  }
}
