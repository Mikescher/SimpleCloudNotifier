import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:share_plus/share_plus.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:simplecloudnotifier/utils/ui.dart';

class ChannelViewPage extends StatefulWidget {
  const ChannelViewPage({
    required this.channel,
    required this.subscription,
    super.key,
  });

  final Channel channel;
  final Subscription? subscription;

  @override
  State<ChannelViewPage> createState() => _ChannelViewPageState();
}

class _ChannelViewPageState extends State<ChannelViewPage> {
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Channel',
      showSearch: false,
      showShare: false,
      child: _buildChannelView(context),
    );
  }

  Widget _buildChannelView(BuildContext context) {
    final userAccUserID = context.select<AppAuth, String?>((v) => v.userID);

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            ..._buildChannelHeader(context),
            SizedBox(height: 8),
            _buildQRCode(context),
            SizedBox(height: 8),
            //TODO icons
            _buildMetaCard(context, FontAwesomeIcons.solidQuestion, 'ChannelID', ['...'], null),
            _buildMetaCard(context, FontAwesomeIcons.solidQuestion, 'InternalName', ['...'], null),
            _buildMetaCard(context, FontAwesomeIcons.solidQuestion, 'DisplayName', ['...'], null), //TODO edit icon on right to edit name
            _buildMetaCard(context, FontAwesomeIcons.solidQuestion, 'Subscription (own)', ['...'], null), //TODO sub/unsub icon on right
            //TODO list foreign subscriptions (with accept/decline/delete button on right)
            _buildMetaCard(context, FontAwesomeIcons.solidQuestion, 'Messages', ['...'], () {/*TODO*/}),
          ],
        ),
      ),
    );
  }

  List<Widget> _buildChannelHeader(BuildContext context) {
    return [
      Text(widget.channel.displayName, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
    ];
  }

  Widget _buildMetaCard(BuildContext context, IconData icn, String title, List<String> values, void Function()? action) {
    final container = UI.box(
      context: context,
      padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
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

  Widget _buildQRCode(BuildContext context) {
    var text = 'TODO' + widget.channel.channelID; //TODO subkey+channelid with deeplink-y
    return GestureDetector(
      onTap: () {
        //TODO share
      },
      child: Center(
        child: QrImageView(
          data: text,
          version: QrVersions.auto,
          size: 300.0,
          eyeStyle: QrEyeStyle(
            eyeShape: QrEyeShape.square,
            color: Theme.of(context).textTheme.bodyLarge?.color,
          ),
          dataModuleStyle: QrDataModuleStyle(
            dataModuleShape: QrDataModuleShape.square,
            color: Theme.of(context).textTheme.bodyLarge?.color,
          ),
        ),
      ),
    );
  }
}
