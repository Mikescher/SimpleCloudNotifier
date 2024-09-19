import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/state/app_events.dart';

class FilterModalPriority extends StatefulWidget {
  @override
  _FilterModalPriorityState createState() => _FilterModalPriorityState();
}

class _FilterModalPriorityState extends State<FilterModalPriority> {
  Set<int> _selectedEntries = {};

  Map<int, (String, String)> _texts = {
    0: ('Low (0)', 'Low'),
    1: ('Normal (1)', 'Normal'),
    2: ('High (2)', 'High'),
  };

  void toggleEntry(int entry) {
    setState(() {
      if (_selectedEntries.contains(entry)) {
        _selectedEntries.remove(entry);
      } else {
        _selectedEntries.add(entry);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Priority'),
      content: Container(
        width: 0,
        height: 200,
        child: ListView.builder(
          shrinkWrap: true,
          itemBuilder: (builder, index) {
            return ListTile(
              title: Text(_texts[index]?.$1 ?? '???'),
              leading: Icon(_selectedEntries.contains(index) ? Icons.check_box : Icons.check_box_outline_blank, color: Theme.of(context).primaryColor),
              onTap: () => toggleEntry(index),
            );
          },
          itemCount: 3,
        ),
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

    final chiplets = _selectedEntries.map((e) => MessageFilterChiplet(label: _texts[e]?.$2 ?? '???', value: e, type: MessageFilterChipletType.priority)).toList();

    AppEvents().notifyFilterListeners([MessageFilterChipletType.priority], chiplets);
  }
}
