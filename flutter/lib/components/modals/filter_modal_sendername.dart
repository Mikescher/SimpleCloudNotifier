import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';

class FilterModalSendername extends StatefulWidget {
  @override
  _FilterModalSendernameState createState() => _FilterModalSendernameState();
}

class _FilterModalSendernameState extends State<FilterModalSendername> {
  Set<String> _selectedEntries = {};

  late ImmediateFuture<List<String>>? _futureSenders;

  @override
  void initState() {
    super.initState();

    _futureSenders = null;
    _futureSenders = ImmediateFuture.ofFuture(() async {
      final userAcc = Provider.of<AppAuth>(context, listen: false);
      if (!userAcc.isAuth()) throw new Exception('not logged in');

      final senders = await APIClient.getSenderNameList(userAcc);

      return senders;
    }());
  }

  void toggleEntry(String senderID) {
    setState(() {
      if (_selectedEntries.contains(senderID)) {
        _selectedEntries.remove(senderID);
      } else {
        _selectedEntries.add(senderID);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Senders'),
      content: Container(
        width: 9000,
        height: 9000,
        child: () {
          if (_futureSenders == null) {
            return Center(child: CircularProgressIndicator());
          }

          return FutureBuilder(
            future: _futureSenders!.future,
            builder: ((context, snapshot) {
              if (_futureSenders?.value != null) {
                return _buildList(context, _futureSenders!.value!);
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
              label: e,
              value: e,
              type: MessageFilterChipletType.sender,
            ))
        .toList();

    AppEvents().notifyFilterListeners([MessageFilterChipletType.sender], chiplets);
  }

  Widget _buildList(BuildContext context, List<String> list) {
    return ListView.builder(
      shrinkWrap: true,
      itemBuilder: (builder, index) {
        final sender = list[index];
        return ListTile(
          title: Text(sender),
          leading: Icon(_selectedEntries.contains(sender) ? Icons.check_box : Icons.check_box_outline_blank, color: Theme.of(context).primaryColor),
          onTap: () => toggleEntry(sender),
          visualDensity: VisualDensity(vertical: -4),
        );
      },
      itemCount: list.length,
    );
  }
}
