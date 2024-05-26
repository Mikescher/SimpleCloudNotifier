import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

class ChannelListItem extends StatefulWidget {
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  const ChannelListItem({
    required this.channel,
    required this.onPressed,
    super.key,
  });

  final Channel channel;
  final Null Function() onPressed;

  @override
  State<ChannelListItem> createState() => _ChannelListItemState();
}

class _ChannelListItemState extends State<ChannelListItem> {
  Message? lastMessage;

  @override
  void initState() {
    super.initState();

    final acc = Provider.of<UserAccount>(context, listen: false);

    if (acc.auth != null) {
      () async {
        final (_, channelMessages) = await APIClient.getMessageList(acc.auth!, '@start', pageSize: 1, channelIDs: [widget.channel.channelID]);
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
        splashColor: Theme.of(context).splashColor,
        onTap: widget.onPressed,
        child: Padding(
          padding: const EdgeInsets.all(8),
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
                      lastMessage?.title ?? '...',
                      style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)),
                    ),
                  ),
                  Text(widget.channel.messagesSent.toString(), style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold)),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
