import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:toastification/toastification.dart';

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
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: () => Toaster.success("Hello World", "This was a triumph!"),
                child: const Text('Show Success Notification'),
              ),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: () => Toaster.info("Hello World", "This was a triumph!"),
                child: const Text('Show Info Notification'),
              ),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: () => Toaster.warn("Hello World", "This was a triumph!"),
                child: const Text('Show Warn Notification'),
              ),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: () => Toaster.error("Hello World", "This was a triumph!"),
                child: const Text('Show Info Notification'),
              ),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: () => Toaster.simple("Hello World"),
                child: const Text('Show Simple Notification'),
              ),
              SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }
}
