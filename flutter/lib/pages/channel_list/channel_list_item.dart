import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/pages/channel_message_view/channel_message_view.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/scn_data_cache.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

class ChannelListItem extends StatefulWidget {
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  const ChannelListItem({
    required this.channel,
    required this.onPressed,
    required this.subscription,
    super.key,
  });

  final Channel channel;
  final Subscription? subscription;
  final Null Function() onPressed;

  @override
  State<ChannelListItem> createState() => _ChannelListItemState();
}

class _ChannelListItemState extends State<ChannelListItem> {
  SCNMessage? lastMessage;

  @override
  void initState() {
    super.initState();

    final acc = Provider.of<AppAuth>(context, listen: false);

    if (acc.isAuth()) {
      lastMessage = SCNDataCache().getMessagesSorted().where((p) => p.channelID == widget.channel.channelID).firstOrNull;

      () async {
        final (_, channelMessages) = await APIClient.getMessageList(acc, '@start', pageSize: 1, channelIDs: [widget.channel.channelID]);
        setState(() {
          lastMessage = channelMessages.firstOrNull;
        });
      }();
    }
  }

  @override
  Widget build(BuildContext context) {
    //TODO subscription status
    return Card.filled(
      margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
      shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
      color: Theme.of(context).cardTheme.color,
      child: InkWell(
        onTap: widget.onPressed,
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: Row(
            children: [
              _buildIcon(context),
              SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            widget.channel.displayName,
                            style: const TextStyle(fontWeight: FontWeight.bold),
                          ),
                        ),
                        Text(
                          (widget.channel.timestampLastSent == null) ? '' : ChannelListItem._dateFormat.format(DateTime.parse(widget.channel.timestampLastSent!).toLocal()),
                          style: const TextStyle(fontSize: 14),
                        ),
                      ],
                    ),
                    SizedBox(height: 4),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Expanded(
                          child: Text(
                            _preformatTitle(lastMessage),
                            style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)),
                          ),
                        ),
                        Text(widget.channel.messagesSent.toString(), style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold)),
                      ],
                    ),
                  ],
                ),
              ),
              SizedBox(width: 4),
              GestureDetector(
                onTap: () {
                  Navi.push(context, () => ChannelMessageViewPage(channel: this.widget.channel));
                },
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: Icon(FontAwesomeIcons.solidEnvelopes, color: Theme.of(context).colorScheme.onPrimaryContainer.withAlpha(128), size: 24),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _preformatTitle(SCNMessage? message) {
    if (message == null) return '...';
    return message.title.replaceAll('\n', '').replaceAll('\r', '').replaceAll('\t', ' ');
  }

  Widget _buildIcon(BuildContext context) {
    if (widget.subscription == null) {
      return Icon(FontAwesomeIcons.solidSquareDashed, color: Theme.of(context).colorScheme.outline, size: 32); // not-subscribed
    } else if (widget.subscription!.confirmed && widget.channel.ownerUserID == widget.subscription!.subscriberUserID) {
      return Icon(FontAwesomeIcons.solidSquareRss, color: Theme.of(context).colorScheme.onPrimaryContainer, size: 32); // subscribed (own channel)
    } else if (widget.subscription!.confirmed) {
      return Icon(FontAwesomeIcons.solidSquareShareNodes, color: Theme.of(context).colorScheme.onPrimaryContainer, size: 32); // subscribed (foreign channel)
    } else {
      return Icon(FontAwesomeIcons.solidSquareEnvelope, color: Theme.of(context).colorScheme.tertiary, size: 32); // requested
    }
  }
}
