import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';

class FilterModalChannel extends StatefulWidget {
  @override
  _FilterModalChannelState createState() => _FilterModalChannelState();
}

class _FilterModalChannelState extends State<FilterModalChannel> {
  Set<String> _selectedEntries = {};

  late ImmediateFuture<List<Channel>>? _futureChannels;

  @override
  void initState() {
    super.initState();

    _futureChannels = null;
    _futureChannels = ImmediateFuture.ofFuture(() async {
      final userAcc = Provider.of<AppAuth>(context, listen: false);
      if (!userAcc.isAuth()) throw new Exception('not logged in');

      final channels = await APIClient.getChannelList(userAcc, ChannelSelector.all);

      return channels.where((p) => p.subscription?.confirmed ?? false).map((e) => e.channel).toList(); // return only subscribed channels
    }());
  }

  void toggleEntry(String channelID) {
    setState(() {
      if (_selectedEntries.contains(channelID)) {
        _selectedEntries.remove(channelID);
      } else {
        _selectedEntries.add(channelID);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Channels'),
      content: Container(
        width: 9000,
        height: 9000,
        child: () {
          if (_futureChannels == null) {
            return Center(child: CircularProgressIndicator());
          }

          return FutureBuilder(
            future: _futureChannels!.future,
            builder: ((context, snapshot) {
              if (_futureChannels?.value != null) {
                return _buildList(context, _futureChannels!.value!);
              } else if (snapshot.connectionState == ConnectionState.done && snapshot.hasError) {
                return Text('Error: ${snapshot.error}'); //TODO better error display
              } else if (snapshot.connectionState == ConnectionState.done) {
                return _buildList(context, snapshot.data!);
              } else {
                return Center(child: CircularProgressIndicator());
              }
            }),
          );
        }(),
      ),
      actions: <Widget>[
        TextButton(
          style: TextButton.styleFrom(textStyle: Theme.of(context).textTheme.labelLarge),
          child: const Text('Apply'),
          onPressed: () {
            onOkay();
          },
        ),
      ],
    );
  }

  void onOkay() {
    Navigator.of(context).pop();

    final chiplets = _selectedEntries
        .map((e) => MessageFilterChiplet(
              label: _futureChannels?.get()?.map((e) => e as Channel?).firstWhere((p) => p?.channelID == e, orElse: () => null)?.displayName ?? '???',
              value: e,
              type: MessageFilterChipletType.channel,
            ))
        .toList();

    AppEvents().notifyFilterListeners([MessageFilterChipletType.channel], chiplets);
  }

  Widget _buildList(BuildContext context, List<Channel> list) {
    return ListView.builder(
      shrinkWrap: true,
      itemBuilder: (builder, index) {
        final channel = list[index];
        return ListTile(
          title: Text(channel.displayName),
          leading: Icon(_selectedEntries.contains(channel.channelID) ? Icons.check_box : Icons.check_box_outline_blank, color: Theme.of(context).primaryColor),
          onTap: () => toggleEntry(channel.channelID),
          visualDensity: VisualDensity(vertical: -4),
        );
      },
      itemCount: list.length,
    );
  }
}
