import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/state/application_log.dart';

class AppEvents {
  static AppEvents? _singleton = AppEvents._internal();

  factory AppEvents() {
    return _singleton ?? (_singleton = AppEvents._internal());
  }

  AppEvents._internal() {}

  List<void Function(String)> _searchListeners = [];
  List<void Function(SCNMessage)> _messageReceivedListeners = [];

  void subscribeSearchListener(void Function(String) listener) {
    _searchListeners.add(listener);
  }

  void unsubscribeSearchListener(void Function(String) listener) {
    _searchListeners.remove(listener);
  }

  void notifySearchListeners(String query) {
    ApplicationLog.debug('[AppEvents] onSearch: $query');

    for (var listener in _searchListeners) {
      listener(query);
    }
  }

  void subscribeMessageReceivedListener(void Function(SCNMessage) listener) {
    _messageReceivedListeners.add(listener);
  }

  void unsubscribeMessageReceivedListener(void Function(SCNMessage) listener) {
    _messageReceivedListeners.remove(listener);
  }

  void notifyMessageReceivedListeners(SCNMessage msg) {
    ApplicationLog.debug('[AppEvents] onMessageReceived: ${msg.messageID}');

    for (var listener in _messageReceivedListeners) {
      listener(msg);
    }
  }
}
