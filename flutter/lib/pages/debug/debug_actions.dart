import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/utils/notifier.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:simplecloudnotifier/utils/ui.dart';

class DebugActionsPage extends StatefulWidget {
  @override
  _DebugActionsPageState createState() => _DebugActionsPageState();
}

class _DebugActionsPageState extends State<DebugActionsPage> {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              UI.button(
                big: false,
                onPressed: () => Toaster.success("Hello World", "This was a triumph!"),
                text: 'Show Success Notification',
              ),
              SizedBox(height: 4),
              UI.button(
                big: false,
                onPressed: () => Toaster.info("Hello World", "This was a triumph!"),
                text: 'Show Info Notification',
              ),
              SizedBox(height: 4),
              UI.button(
                big: false,
                onPressed: () => Toaster.warn("Hello World", "This was a triumph!"),
                text: 'Show Warn Notification',
              ),
              SizedBox(height: 4),
              UI.button(
                big: false,
                onPressed: () => Toaster.error("Hello World", "This was a triumph!"),
                text: 'Show Info Notification',
              ),
              SizedBox(height: 4),
              UI.button(
                big: false,
                onPressed: () => Toaster.simple("Hello World"),
                text: 'Show Simple Notification',
              ),
              SizedBox(height: 20),
              UI.button(
                big: false,
                onPressed: _sendTokenToServer,
                text: 'Send FCM Token to Server',
              ),
              SizedBox(height: 20),
              UI.button(
                big: false,
                onPressed: () => Notifier.showLocalNotification('TEST_CHANNEL', "Test Channel", "Channel for testing", "Hello World", "Local Notification test", null),
                text: 'Show local notification',
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _sendTokenToServer() {
    //TODO
  }
}
