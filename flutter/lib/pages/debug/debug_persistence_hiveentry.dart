import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';

class DebugHiveEntryPage extends StatelessWidget {
  final FieldDebuggable value;
  final List<(String, String)> fields;

  DebugHiveEntryPage({required this.value}) : fields = value.debugFieldList();

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'HiveEntry',
      showSearch: false,
      child: ListView.separated(
        itemCount: fields.length,
        itemBuilder: (context, listIndex) {
          return ListTile(
            dense: true,
            contentPadding: EdgeInsets.fromLTRB(8, 0, 8, 0),
            visualDensity: VisualDensity(horizontal: 0, vertical: -4),
            title: Text(fields[listIndex].$1, style: TextStyle(fontWeight: FontWeight.bold)),
            subtitle: Text(fields[listIndex].$2, style: TextStyle(fontFamily: "monospace")),
          );
        },
        separatorBuilder: (context, index) => Divider(),
      ),
    );
  }
}
