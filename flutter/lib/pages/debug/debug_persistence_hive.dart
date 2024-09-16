import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_hiveentry.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class DebugHiveBoxPage extends StatelessWidget {
  final String boxName;
  final Box<FieldDebuggable> box;

  DebugHiveBoxPage({required this.boxName, required this.box});

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Hive: ' + boxName,
      showSearch: false,
      child: ListView.separated(
        itemCount: box.length,
        itemBuilder: (context, listIndex) {
          return GestureDetector(
            onTap: () {
              Navi.push(context, () => DebugHiveEntryPage(value: box.getAt(listIndex)!));
            },
            child: Container(
              padding: EdgeInsets.fromLTRB(8, 4, 8, 4),
              child: Text(box.getAt(listIndex).toString(), style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12)),
            ),
          );
        },
        separatorBuilder: (context, index) => Divider(),
      ),
    );
  }
}
