import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_hiveentry.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';

class DebugHiveBoxPage extends StatelessWidget {
  final String boxName;
  final Box<FieldDebuggable> box;

  DebugHiveBoxPage({required this.boxName, required this.box});

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Hive: ' + boxName,
      showSearch: false,
      showDebug: false,
      child: ListView.separated(
        itemCount: box.length,
        itemBuilder: (context, listIndex) {
          return GestureDetector(
            onTap: () {
              Navigator.push(context, MaterialPageRoute<DebugHiveEntryPage>(builder: (context) => DebugHiveEntryPage(value: box.getAt(listIndex)!)));
            },
            child: ListTile(
              title: Text(box.getAt(listIndex).toString(), style: TextStyle(fontWeight: FontWeight.bold)),
            ),
          );
        },
        separatorBuilder: (context, index) => Divider(),
      ),
    );
  }
}
