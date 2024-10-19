class APIException implements Exception {
  final int httpStatus;
  final int error;
  final int errHighlight;
  final String message;

  APIException(this.httpStatus, this.error, this.errHighlight, this.message);

  @override
  String toString() {
    return '[$error] $message';
  }
}
