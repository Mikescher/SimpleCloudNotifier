import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/api/api_exception.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/api_error.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';

class MessageViewPage extends StatefulWidget {
  const MessageViewPage({super.key, required this.message});

  final Message message; // Potentially trimmed

  @override
  State<MessageViewPage> createState() => _MessageViewPageState();
}

class _MessageViewPageState extends State<MessageViewPage> {
  late Future<(Message, ChannelWithSubscription?, KeyToken?)>? mainFuture;
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  @override
  void initState() {
    super.initState();
    mainFuture = fetchData();
  }

  Future<(Message, ChannelWithSubscription?, KeyToken?)> fetchData() async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    final msg = await APIClient.getMessage(acc, widget.message.messageID);

    ChannelWithSubscription? chn = null;
    try {
      chn = await APIClient.getChannel(acc, msg.channelID);
    } on APIException catch (e) {
      if (e.error == APIError.USER_AUTH_FAILED) {
        chn = null;
      } else {
        rethrow;
      }
    }

    KeyToken? tok = null;
    try {
      tok = await APIClient.getKeyToken(acc, msg.usedKeyID);
    } on APIException catch (e) {
      if (e.error == APIError.USER_AUTH_FAILED) {
        tok = null;
      } else {
        rethrow;
      }
    }

    return (msg, chn, tok);
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
      //TODO showShare: true
      child: FutureBuilder<(Message, ChannelWithSubscription?, KeyToken?)>(
        future: mainFuture,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            final msg = snapshot.data!.$1;
            final chn = snapshot.data!.$2;
            final tok = snapshot.data!.$3;
            return _buildMessageView(context, msg, chn, tok);
          } else if (snapshot.hasError) {
            return Center(child: Text('${snapshot.error}')); //TODO nice error page
          } else if (!widget.message.trimmed) {
            return _buildLoadingView(context, widget.message);
          } else {
            return const Center(child: CircularProgressIndicator());
          }
        },
      ),
    );
  }

  Widget _buildMessageView(BuildContext context, Message message, ChannelWithSubscription? channel, KeyToken? token) {
    //TODO loading true/false indicator
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ..._buildMessageHeader(context, message, channel, token),
            SizedBox(height: 8),
            if (message.content != null) ..._buildMessageContent(context, message, channel, token),
            SizedBox(height: 8),
            if (message.senderName != null) _buildMetaCard(context, FontAwesomeIcons.solidSignature, 'Sender', [message.senderName!], () => {/*TODO*/}),
            if (token != null) _buildMetaCard(context, FontAwesomeIcons.solidGearCode, 'KeyToken', [token.keytokenID, token.name], () => {/*TODO*/}),
            _buildMetaCard(context, FontAwesomeIcons.solidIdCardClip, 'MessageID', [message.messageID, if (message.userMessageID != null) message.userMessageID!], null),
            if (channel != null) _buildMetaCard(context, FontAwesomeIcons.solidSnake, 'Channel', [message.channelID, channel.channel.displayName], () => {/*TODO*/}),
            _buildMetaCard(context, FontAwesomeIcons.solidTimer, 'Timestamp', [message.timestamp], null),
          ],
        ),
      ),
    );
  }

  Widget _buildLoadingView(BuildContext context, Message message) {
    //TODO loading / skeleton use limitdata
    return SizedBox();
  }

  String _resolveChannelName(ChannelWithSubscription? channel, Message message) {
    return channel?.channel.displayName ?? message.channelInternalName;
  }

  List<Widget> _buildMessageHeader(BuildContext context, Message message, ChannelWithSubscription? channel, KeyToken? token) {
    return [
      Row(
        children: [
          Container(
            padding: const EdgeInsets.fromLTRB(4, 0, 4, 0),
            margin: const EdgeInsets.fromLTRB(0, 0, 4, 0),
            decoration: BoxDecoration(
              color: Theme.of(context).hintColor,
              borderRadius: BorderRadius.all(Radius.circular(4)),
            ),
            child: Text(
              _resolveChannelName(channel, message),
              style: TextStyle(fontWeight: FontWeight.bold, color: Theme.of(context).cardColor, fontSize: 16),
              overflow: TextOverflow.clip,
              maxLines: 1,
            ),
          ),
          Expanded(child: SizedBox()),
          Text(_dateFormat.format(DateTime.parse(message.timestamp)), style: const TextStyle(fontSize: 14)),
        ],
      ),
      SizedBox(height: 8),
      Text(message.title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
    ];
  }

  List<Widget> _buildMessageContent(BuildContext context, Message message, ChannelWithSubscription? channel, KeyToken? token) {
    return [
      Row(
        children: [
          Expanded(child: SizedBox()),
          IconButton(
            icon: FaIcon(FontAwesomeIcons.copy),
            iconSize: 18,
            padding: EdgeInsets.all(4),
            constraints: BoxConstraints(),
            style: ButtonStyle(tapTargetSize: MaterialTapTargetSize.shrinkWrap),
            onPressed: () {/*TODO*/},
          ),
          IconButton(
            icon: FaIcon(FontAwesomeIcons.lineColumns),
            iconSize: 18,
            padding: EdgeInsets.all(4),
            constraints: BoxConstraints(),
            style: ButtonStyle(tapTargetSize: MaterialTapTargetSize.shrinkWrap),
            onPressed: () {/*TODO*/},
          ),
        ],
      ),
      Container(
        decoration: BoxDecoration(
          border: Border.all(color: Theme.of(context).hintColor),
          borderRadius: BorderRadius.circular(4),
        ),
        padding: const EdgeInsets.all(4),
        child: Text(message.content ?? ''),
      ),
    ];
  }

  Widget _buildMetaCard(BuildContext context, IconData icn, String title, List<String> values, void Function()? action) {
    final container = Container(
      padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
      decoration: BoxDecoration(
        border: Border.all(color: Theme.of(context).hintColor),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Row(
        children: [
          FaIcon(icn, size: 18),
          SizedBox(width: 16),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
              for (final val in values) Text(val, style: const TextStyle(fontSize: 14)),
            ],
          ),
        ],
      ),
    );

    if (action == null) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: container,
      );
    } else {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: action,
          child: container,
        ),
      );
    }
  }
}
