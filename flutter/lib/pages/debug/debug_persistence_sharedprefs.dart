import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';

class DebugSharedPrefPage extends StatelessWidget {
  final SharedPreferences sharedPref;
  final List<String> keys;

  DebugSharedPrefPage({required this.sharedPref}) : keys = sharedPref.getKeys().toList();

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'SharedPreferences',
      showSearch: false,
      child: ListView.separated(
        itemCount: sharedPref.getKeys().length,
        itemBuilder: (context, listIndex) {
          return ListTile(
            title: Text(keys[listIndex], style: TextStyle(fontWeight: FontWeight.bold)),
            subtitle: Text(sharedPref.get(keys[listIndex]).toString()),
          );
        },
        separatorBuilder: (context, index) => Divider(),
      ),
    );
  }
}
