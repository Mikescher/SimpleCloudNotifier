import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/keytoken.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';

class FilterModalKeytoken extends StatefulWidget {
  @override
  _FilterModalKeytokenState createState() => _FilterModalKeytokenState();
}

class _FilterModalKeytokenState extends State<FilterModalKeytoken> {
  Set<String> _selectedEntries = {};

  late ImmediateFuture<List<KeyToken>>? _futureKeyTokens;

  @override
  void initState() {
    super.initState();

    _futureKeyTokens = null;
    _futureKeyTokens = ImmediateFuture.ofFuture(() async {
      final userAcc = Provider.of<AppAuth>(context, listen: false);
      if (!userAcc.isAuth()) throw new Exception('not logged in');

      final toks = await APIClient.getKeyTokenList(userAcc);

      return toks;
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
          if (_futureKeyTokens == null) {
            return Center(child: CircularProgressIndicator());
          }

          return FutureBuilder(
            future: _futureKeyTokens!.future,
            builder: ((context, snapshot) {
              if (_futureKeyTokens?.value != null) {
                return _buildList(context, _futureKeyTokens!.value!);
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
              label: _futureKeyTokens?.get()?.map((e) => e as KeyToken?).firstWhere((p) => p?.keytokenID == e, orElse: () => null)?.name ?? '???',
              value: e,
              type: MessageFilterChipletType.sender,
            ))
        .toList();

    AppEvents().notifyFilterListeners([MessageFilterChipletType.sender], chiplets);
  }

  Widget _buildList(BuildContext context, List<KeyToken> list) {
    return ListView.builder(
      shrinkWrap: true,
      itemBuilder: (builder, index) {
        final sender = list[index];
        return ListTile(
          title: Text(sender.name),
          leading: Icon(_selectedEntries.contains(sender.keytokenID) ? Icons.check_box : Icons.check_box_outline_blank, color: Theme.of(context).primaryColor),
          onTap: () => toggleEntry(sender.keytokenID),
          visualDensity: VisualDensity(vertical: -4),
        );
      },
      itemCount: list.length,
    );
  }
}
