import 'dart:io';

import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';

class Notifier {
  static void showLocalNotification(String messageID, String channelID, String channelName, String channelDescr, String title, String body, DateTime? timestamp) async {
    final nid = Globals().sharedPrefs.getInt('notifier.nextid') ?? 1000;
    Globals().sharedPrefs.setInt('notifier.nextid', nid + 7);

    final existingSummaryNID = Globals().sharedPrefs.getInt('notifier.summary.$channelID');

    final flutterLocalNotificationsPlugin = FlutterLocalNotificationsPlugin();

    if (Platform.isAndroid && AppSettings().groupNotifications) {
      final activeNotifications = (await flutterLocalNotificationsPlugin.getActiveNotifications()).where((p) => p.groupKey == channelID).toList();
      final summaryNotification = activeNotifications.where((p) => p.id == existingSummaryNID).toList();

      ApplicationLog.debug('found ${activeNotifications.length} active notifications in this group (${summaryNotification.length} summary notifications for channel ${channelID} with nid [${existingSummaryNID}])');

      if (activeNotifications.isNotEmpty && !activeNotifications.any((p) => p.id == existingSummaryNID)) {
        // ======== SHOW SUMMARY/GROUPING NOTIFICATION ========
        final newSummaryNID = nid + 1;
        ApplicationLog.debug('Create new summary notifications for channel ${channelID} with nid [${newSummaryNID}])');
        Globals().sharedPrefs.setInt('notifier.summary.$channelID', newSummaryNID);

        var payload = '';
        if (messageID != '') {
          payload = ['@SCN_MESSAGE_SUMMARY', channelID, newSummaryNID].join("\n");
        }

        await flutterLocalNotificationsPlugin.show(
          newSummaryNID,
          channelName,
          "(multiple notifications)",
          NotificationDetails(
            android: AndroidNotificationDetails(
              channelID,
              channelName,
              importance: Importance.max,
              priority: Priority.high,
              groupKey: channelID,
              setAsGroupSummary: true,
              subText: (channelName == 'main') ? null : channelName,
            ),
          ),
          payload: payload,
        );
      }
    }

    final newMessageNID = nid + 2;

    ApplicationLog.debug('Create new local notifications for message in channel ${channelID} with nid [${newMessageNID}])');

    var payload = '';
    if (messageID != '') {
      payload = ['@SCN_MESSAGE', messageID, channelID, newMessageNID].join("\n");
    }

    // ======== SHOW NOTIFICATION ========
    await flutterLocalNotificationsPlugin.show(
      newMessageNID,
      title,
      body,
      NotificationDetails(
        android: AndroidNotificationDetails(
          channelID,
          channelName,
          channelDescription: channelDescr,
          importance: Importance.max,
          priority: Priority.high,
          when: timestamp?.millisecondsSinceEpoch,
          groupKey: channelID,
          subText: (channelName == 'main') ? null : channelName,
        ),
      ),
      payload: payload,
    );
  }
}
