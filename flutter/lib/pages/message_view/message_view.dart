import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

class MessageViewPage extends StatefulWidget {
  const MessageViewPage({super.key, required this.messageID});

  final String messageID;

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
    final acc = Provider.of<UserAccount>(context, listen: false);

    return await APIClient.getMessage(acc.auth!, widget.messageID);
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Message',
      child: FutureBuilder<Message>(
        future: futureMessage,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return Center(child: Text(snapshot.data!.title));
          } else if (snapshot.hasError) {
            return Center(child: Text('${snapshot.error}')); //TODO nice error page
          }

          return const Center(child: CircularProgressIndicator());
        },
      ),
    );
  }
}
