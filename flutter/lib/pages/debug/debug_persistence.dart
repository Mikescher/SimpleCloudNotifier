import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_hive.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_sharedprefs.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class DebugPersistencePage extends StatefulWidget {
  @override
  _DebugPersistencePageState createState() => _DebugPersistencePageState();
}

class _DebugPersistencePageState extends State<DebugPersistencePage> {
  SharedPreferences? prefs = null;

  @override
  void initState() {
    super.initState();

    SharedPreferences.getInstance().then((value) => setState(() => prefs = value));
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Card.outlined(
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: GestureDetector(
                onTap: () {
                  Navi.push(context, () => DebugSharedPrefPage(sharedPref: prefs!));
                },
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    SizedBox(width: 30, child: Text('')),
                    Expanded(child: Text('Shared Preferences', style: TextStyle(fontWeight: FontWeight.bold), textAlign: TextAlign.center)),
                    SizedBox(width: 30, child: Text('${prefs?.getKeys().length.toString()}', textAlign: TextAlign.end)),
                  ],
                ),
              ),
            ),
          ),
          Card.outlined(
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: GestureDetector(
                onTap: () {
                  Navi.push(context, () => DebugHiveBoxPage(boxName: 'scn-requests', box: Hive.box<SCNRequest>('scn-requests')));
                },
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    SizedBox(width: 30, child: Text('')),
                    Expanded(child: Text('Hive [scn-requests]', style: TextStyle(fontWeight: FontWeight.bold), textAlign: TextAlign.center)),
                    SizedBox(width: 30, child: Text('${Hive.box<SCNRequest>('scn-requests').length.toString()}', textAlign: TextAlign.end)),
                  ],
                ),
              ),
            ),
          ),
          Card.outlined(
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: GestureDetector(
                onTap: () {
                  Navi.push(context, () => DebugHiveBoxPage(boxName: 'scn-requests', box: Hive.box<SCNLog>('scn-logs')));
                },
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    SizedBox(width: 30, child: Text('')),
                    Expanded(child: Text('Hive [scn-logs]', style: TextStyle(fontWeight: FontWeight.bold), textAlign: TextAlign.center)),
                    SizedBox(width: 30, child: Text('${Hive.box<SCNLog>('scn-logs').length.toString()}', textAlign: TextAlign.end)),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
