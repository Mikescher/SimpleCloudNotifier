abstract class TokenSource {
  String getToken();
  String getUserID();
}

class DirectTokenSource implements TokenSource {
  final String _userID;
  final String _token;

  DirectTokenSource(this._userID, this._token);

  @override
  String getUserID() {
    return _userID;
  }

  @override
  String getToken() {
    return _token;
  }
}
