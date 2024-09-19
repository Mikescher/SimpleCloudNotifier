import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';

class FilterModalTime extends StatefulWidget {
  @override
  _FilterModalTimeState createState() => _FilterModalTimeState();
}

class _FilterModalTimeState extends State<FilterModalTime> {
  DateTime? _tsBefore = null;
  DateTime? _tsAfter = null;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Timerange'),
      content: Container(
        width: 9000,
        height: 9000,
        child: Placeholder(),
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

    //TODO
  }
}
