import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';

class MessageViewPage extends StatefulWidget {
  const MessageViewPage({super.key, required this.message});

  final Message message; // Potentially trimmed

  @override
  State<MessageViewPage> createState() => _MessageViewPageState();
}

class _MessageViewPageState extends State<MessageViewPage> {
  late Future<Message>? futureMessage;

  @override
  void initState() {
    super.initState();
    futureMessage = fetchMessage();
  }

  Future<Message> fetchMessage() async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    return await APIClient.getMessage(acc, widget.message.messageID);
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Message',
      showSearch: false,
      child: FutureBuilder<Message>(
        future: futureMessage,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return buildMessageView(snapshot.data!, false);
          } else if (snapshot.hasError) {
            return Center(child: Text('${snapshot.error}')); //TODO nice error page
          } else if (!widget.message.trimmed) {
            return buildMessageView(widget.message, true);
          } else {
            return const Center(child: CircularProgressIndicator());
          }
        },
      ),
    );
  }

  Widget buildMessageView(Message message, bool loading) {
    //TODO loading true/false indicator
    return Center(
      child: Column(
        children: [
          Text(message.title),
          Text(message.content ?? ''),
        ],
      ),
    );
  }
}
