import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:share_plus/share_plus.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/pages/channel_view/channel_view.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:simplecloudnotifier/utils/ui.dart';

class MessageViewPage extends StatefulWidget {
  const MessageViewPage({
    super.key,
    required this.messageID,
    required this.preloadedData,
  });

  final String messageID; // Potentially trimmed
  final (SCNMessage,)? preloadedData; // Message is potentially trimmed, whole object is potentially null

  @override
  State<MessageViewPage> createState() => _MessageViewPageState();
}

class _MessageViewPageState extends State<MessageViewPage> {
  late Future<(SCNMessage, ChannelPreview, KeyTokenPreview, UserPreview)>? mainFuture;
  (SCNMessage, ChannelPreview, KeyTokenPreview, UserPreview)? mainFutureSnapshot = null;
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  bool _monospaceMode = false;

  SCNMessage? message = null;

  @override
  void initState() {
    if (widget.preloadedData != null) {
      message = widget.preloadedData!.$1;
    }

    mainFuture = fetchData();
    super.initState();
  }

  Future<(SCNMessage, ChannelPreview, KeyTokenPreview, UserPreview)> fetchData() async {
    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      AppBarState().setLoadingIndeterminate(true);

      final acc = Provider.of<AppAuth>(context, listen: false);

      final msg = await APIClient.getMessage(acc, widget.messageID);

      final fut_chn = APIClient.getChannelPreview(acc, msg.channelID);
      final fut_key = APIClient.getKeyTokenPreview(acc, msg.usedKeyID);
      final fut_usr = APIClient.getUserPreview(acc, msg.senderUserID);

      final chn = await fut_chn;
      final key = await fut_key;
      final usr = await fut_usr;

      //await Future.delayed(const Duration(seconds: 10), () {});

      final r = (msg, chn, key, usr);

      mainFutureSnapshot = r;

      return r;
    } finally {
      AppBarState().setLoadingIndeterminate(false);
    }
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
      showShare: true,
      onShare: _share,
      child: FutureBuilder<(SCNMessage, ChannelPreview, KeyTokenPreview, UserPreview)>(
        future: mainFuture,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            final (msg, chn, tok, usr) = snapshot.data!;
            return _buildMessageView(context, msg, chn, tok, usr);
          } else if (snapshot.hasError) {
            return Center(child: Text('${snapshot.error}')); //TODO nice error page
          } else if (message != null && !this.message!.trimmed) {
            return _buildMessageView(context, this.message!, null, null, null);
          } else {
            return const Center(child: CircularProgressIndicator());
          }
        },
      ),
    );
  }

  void _share() async {
    if (this.message == null) return;

    var msg = this.message!;
    if (mainFutureSnapshot != null) {
      (msg, _, _, _) = mainFutureSnapshot!;
    }

    if (msg.content != null) {
      final result = await Share.share(msg.content!, subject: msg.title);

      if (result.status == ShareResultStatus.unavailable) {
        Toaster.error('Error', "Failed to open share dialog");
      }
    } else {
      final result = await Share.share(msg.title);

      if (result.status == ShareResultStatus.unavailable) {
        Toaster.error('Error', "Failed to open share dialog");
      }
    }
  }

  Widget _buildMessageView(BuildContext context, SCNMessage message, ChannelPreview? channel, KeyTokenPreview? token, UserPreview? user) {
    final userAccUserID = context.select<AppAuth, String?>((v) => v.userID);

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ..._buildMessageHeader(context, message, channel),
            SizedBox(height: 8),
            if (message.content != null) ..._buildMessageContent(context, message),
            SizedBox(height: 8),
            if (message.senderName != null)
              UI.metaCard(
                context: context,
                icon: FontAwesomeIcons.solidSignature,
                title: 'Sender',
                values: [message.senderName!],
                mainAction: () => {/*TODO*/},
              ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidGearCode,
              title: 'KeyToken',
              values: [message.usedKeyID, token?.name ?? '...'],
              mainAction: () => {/*TODO*/},
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidIdCardClip,
              title: 'MessageID',
              values: [message.messageID, message.userMessageID ?? ''],
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidSnake,
              title: 'Channel',
              values: [message.channelID, channel?.displayName ?? message.channelInternalName],
              mainAction: (channel != null)
                  ? () {
                      Navi.push(context, () => ChannelViewPage(channelID: channel.channelID, preloadedData: null, needsReload: null));
                    }
                  : null,
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidTimer,
              title: 'Timestamp',
              values: [message.timestamp],
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidUser,
              title: 'User',
              values: [user?.userID ?? '...', user?.username ?? ''],
              mainAction: () => {/*TODO*/},
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidBolt,
              title: 'Priority',
              values: [_prettyPrintPriority(message.priority)],
              mainAction: () => {/*TODO*/},
            ),
            if (message.senderUserID == userAccUserID) UI.button(text: "Delete Message", onPressed: () {/*TODO*/}, color: Colors.red[900]),
          ],
        ),
      ),
    );
  }

  String _resolveChannelName(ChannelPreview? channel, SCNMessage message) {
    return channel?.displayName ?? message.channelInternalName;
  }

  List<Widget> _buildMessageHeader(BuildContext context, SCNMessage message, ChannelPreview? channel) {
    return [
      Row(
        children: [
          UI.channelChip(
            context: context,
            text: _resolveChannelName(channel, message),
            margin: const EdgeInsets.fromLTRB(0, 0, 4, 0),
            fontSize: 16,
          ),
          Expanded(child: SizedBox()),
          Text(_dateFormat.format(DateTime.parse(message.timestamp)), style: const TextStyle(fontSize: 14)),
        ],
      ),
      SizedBox(height: 8),
      Text(_preformatTitle(message), style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
    ];
  }

  List<Widget> _buildMessageContent(BuildContext context, SCNMessage message) {
    return [
      Row(
        children: [
          if (message.priority == 2) FaIcon(FontAwesomeIcons.solidTriangleExclamation, size: 16, color: Colors.red[900]),
          if (message.priority == 0) FaIcon(FontAwesomeIcons.solidDown, size: 16, color: Colors.lightBlue[900]),
          Expanded(child: SizedBox()),
          UI.buttonIconOnly(
            onPressed: () {
              Clipboard.setData(new ClipboardData(text: message.content ?? ''));
              Toaster.info("Clipboard", 'Copied text to Clipboard');
              print('================= [CLIPBOARD] =================\n${message.content}\n================= [/CLIPBOARD] =================');
            },
            icon: FontAwesomeIcons.copy,
          ),
          UI.buttonIconOnly(
            icon: _monospaceMode ? FontAwesomeIcons.lineColumns : FontAwesomeIcons.alignLeft,
            onPressed: () {
              setState(() {
                _monospaceMode = !_monospaceMode;
              });
            },
          ),
        ],
      ),
      _monospaceMode
          ? UI.box(
              context: context,
              padding: const EdgeInsets.all(4),
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Text(
                  message.content ?? '',
                  style: TextStyle(fontFamily: "monospace", fontFamilyFallback: <String>["Courier"]),
                ),
              ),
              borderColor: (message.priority == 2) ? Colors.red[900] : null,
            )
          : UI.box(
              context: context,
              padding: const EdgeInsets.all(4),
              child: Text(message.content ?? ''),
              borderColor: (message.priority == 2) ? Colors.red[900] : null,
            )
    ];
  }

  String _preformatTitle(SCNMessage message) {
    return message.title.replaceAll('\n', '').replaceAll('\r', '').replaceAll('\t', ' ');
  }

  String _prettyPrintPriority(int priority) {
    switch (priority) {
      case 0:
        return 'Low (0)';
      case 1:
        return 'Normal (1)';
      case 2:
        return 'High (2)';
      default:
        return 'Unknown ($priority)';
    }
  }
}
