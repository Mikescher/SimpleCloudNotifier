import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/pages/channel_message_view/channel_message_view.dart';
import 'package:simplecloudnotifier/pages/channel_view/channel_view.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/scn_data_cache.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

enum ChannelListItemMode {
  Messages,
  Extended,
}

class ChannelListItem extends StatefulWidget {
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  const ChannelListItem({
    required this.channel,
    required this.onChannelListReloadTrigger,
    required this.onSubscriptionChanged,
    required this.subscription,
    required this.mode,
    super.key,
  });

  final Channel channel;
  final Subscription? subscription;
  final void Function() onChannelListReloadTrigger;
  final ChannelListItemMode mode;
  final void Function(String, Subscription?) onSubscriptionChanged;

  @override
  State<ChannelListItem> createState() => _ChannelListItemState();
}

class _ChannelListItemState extends State<ChannelListItem> {
  SCNMessage? lastMessage;

  @override
  void initState() {
    super.initState();

    final acc = Provider.of<AppAuth>(context, listen: false);

    if (acc.isAuth() && widget.mode == ChannelListItemMode.Messages) {
      lastMessage = SCNDataCache().getMessagesSorted().where((p) => p.channelID == widget.channel.channelID).firstOrNull;

      () async {
        final (_, channelMessages) = await APIClient.getChannelMessageList(acc, widget.channel.channelID, '@start', pageSize: 1);
        setState(() {
          lastMessage = channelMessages.firstOrNull;
        });
      }();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Card.filled(
      margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
      shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
      color: Theme.of(context).cardTheme.color,
      child: InkWell(
        onTap: () {
          if (widget.mode == ChannelListItemMode.Messages) {
            Navi.push(context, () => ChannelMessageViewPage(channel: widget.channel));
          } else {
            Navi.push(context, () => ChannelViewPage(channelID: widget.channel.channelID, preloadedData: (widget.channel, widget.subscription), needsReload: widget.onChannelListReloadTrigger));
          }
        },
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
                        Expanded(child: (widget.mode == ChannelListItemMode.Messages) ? Text(_preformatTitle(lastMessage), style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160))) : _buildSubscriptionStateText(context)),
                        (widget.mode == ChannelListItemMode.Messages) ? Text(widget.channel.messagesSent.toString(), style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold)) : Text("", style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold)),
                      ],
                    ),
                  ],
                ),
              ),
              SizedBox(width: 4),
              GestureDetector(
                onTap: () {
                  if (widget.mode == ChannelListItemMode.Messages) {
                    Navi.push(context, () => ChannelViewPage(channelID: widget.channel.channelID, preloadedData: (widget.channel, widget.subscription), needsReload: widget.onChannelListReloadTrigger));
                  } else {
                    Navi.push(context, () => ChannelMessageViewPage(channel: widget.channel));
                  }
                },
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: (widget.mode == ChannelListItemMode.Messages) ? Icon(FontAwesomeIcons.solidSquareInfo, color: Theme.of(context).colorScheme.onPrimaryContainer.withAlpha(128), size: 24) : Icon(FontAwesomeIcons.solidEnvelopes, color: Theme.of(context).colorScheme.onPrimaryContainer.withAlpha(128), size: 24),
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
      Widget result = Icon(FontAwesomeIcons.solidSquareDashed, color: Theme.of(context).colorScheme.outline, size: 32); // not-subscribed
      result = GestureDetector(onTap: () => _subscribe(), child: result);
      return result;
    } else if (widget.subscription!.confirmed && widget.channel.ownerUserID == widget.subscription!.subscriberUserID) {
      Widget result = Icon(FontAwesomeIcons.solidSquareRss, color: Theme.of(context).colorScheme.onPrimaryContainer, size: 32); // subscribed (own channel)
      result = GestureDetector(onTap: () => _unsubscribe(widget.subscription!), child: result);
      return result;
    } else if (widget.subscription!.confirmed) {
      Widget result = Icon(FontAwesomeIcons.solidSquareShareNodes, color: Theme.of(context).colorScheme.onPrimaryContainer, size: 32); // subscribed (foreign channel)
      result = GestureDetector(onTap: () => _unsubscribe(widget.subscription!), child: result);
      return result;
    } else {
      Widget result = Icon(FontAwesomeIcons.solidSquareEnvelope, color: Theme.of(context).colorScheme.tertiary, size: 32); // requested
      result = GestureDetector(onTap: () => _unsubscribe(widget.subscription!), child: result);
      return result;
    }
  }

  Widget _buildSubscriptionStateText(BuildContext context) {
    if (widget.subscription == null) {
      return Text("", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)));
    } else if (widget.subscription!.confirmed && widget.channel.ownerUserID == widget.subscription!.subscriberUserID) {
      return Text("subscribed", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)));
    } else if (widget.subscription!.confirmed) {
      return Text("subscripted (foreign channe)", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)));
    } else {
      return Text("subscription requested", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)));
    }
  }

  void _subscribe() async {
    final acc = AppAuth();

    if (acc.isAuth() && widget.channel.ownerUserID == acc.getUserID()) {
      try {
        var sub = await APIClient.subscribeToChannelbyID(acc, widget.channel.channelID);
        widget.onChannelListReloadTrigger.call();

        widget.onSubscriptionChanged(widget.channel.channelID, sub);

        if (sub.confirmed) {
          Toaster.success("Success", 'Subscribed to channel');
        } else {
          Toaster.success("Success", 'Requested widget.subscription to channel');
        }
      } catch (exc, trace) {
        Toaster.error("Error", 'Failed to subscribe to channel');
        ApplicationLog.error('Failed to subscribe to channel: ' + exc.toString(), trace: trace);
      }
    }
  }

  void _unsubscribe(Subscription sub) async {
    final acc = AppAuth();

    if (acc.isAuth() && widget.channel.ownerUserID == acc.getUserID() && widget.subscription != null) {
      try {
        await APIClient.deleteSubscription(acc, widget.channel.channelID, widget.subscription!.subscriptionID);
        widget.onChannelListReloadTrigger.call();

        widget.onSubscriptionChanged?.call(widget.channel.channelID, null);

        Toaster.success("Success", 'Unsubscribed from channel');
      } catch (exc, trace) {
        Toaster.error("Error", 'Failed to unsubscribe from channel');
        ApplicationLog.error('Failed to unsubscribe from channel: ' + exc.toString(), trace: trace);
      }
    }
  }
}
