import 'package:flutter/material.dart';

class MessageListPage extends StatelessWidget {
  final String title;

  const MessageListPage({super.key, required this.title});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Text(
          title,
          style: const TextStyle(fontSize: 24),
        ),
      ),
    );
  }
}
