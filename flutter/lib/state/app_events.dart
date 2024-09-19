import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/application_log.dart';

class AppEvents {
  static AppEvents? _singleton = AppEvents._internal();

  factory AppEvents() {
    return _singleton ?? (_singleton = AppEvents._internal());
  }

  AppEvents._internal() {}

  // --------------------------------------------------------------------------

  List<void Function(List<MessageFilterChipletType> types, List<MessageFilterChiplet>)> _filterListeners = [];

  void subscribeFilterListener(void Function(List<MessageFilterChipletType> types, List<MessageFilterChiplet>) listener) {
    _filterListeners.add(listener);
  }

  void unsubscribeFilterListener(void Function(List<MessageFilterChipletType> types, List<MessageFilterChiplet>) listener) {
    _filterListeners.remove(listener);
  }

  void notifyFilterListeners(List<MessageFilterChipletType> types, List<MessageFilterChiplet> query) {
    ApplicationLog.debug('[AppEvents] onFilter: [${types.join(" ; ")}], [${query.map((e) => e.label).join('|')}]');

    for (var listener in _filterListeners) {
      listener(types, query);
    }
  }

  // --------------------------------------------------------------------------

  List<void Function(SCNMessage)> _messageReceivedListeners = [];

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
