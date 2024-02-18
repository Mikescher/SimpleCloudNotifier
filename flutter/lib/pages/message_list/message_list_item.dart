import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/models/message.dart';

class MessageListItem extends StatelessWidget {
  const MessageListItem({
    required this.message,
    super.key,
  });

  final Message message;

  @override
  Widget build(BuildContext context) => ListTile(
        leading: const SizedBox(width: 40, height: 40, child: const Placeholder()),
        title: Text(message.messageID),
      );
}
