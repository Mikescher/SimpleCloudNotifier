import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_hive.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_sharedprefs.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/fb_message.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';
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
          _buildSharedPrefCard(context),
          _buildHiveCard(context, () => Hive.box<SCNRequest>('scn-requests'), 'scn-requests'),
          _buildHiveCard(context, () => Hive.box<SCNLog>('scn-logs'), 'scn-logs'),
          _buildHiveCard(context, () => Hive.box<SCNMessage>('scn-message-cache'), 'scn-message-cache'),
          _buildHiveCard(context, () => Hive.box<Channel>('scn-channel-cache'), 'scn-channel-cache'),
          _buildHiveCard(context, () => Hive.box<FBMessage>('scn-fb-messages'), 'scn-fb-messages'),
        ],
      ),
    );
  }

  Widget _buildSharedPrefCard(BuildContext context) {
    return Card.outlined(
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
              SizedBox(width: 40, child: Text('${prefs?.getKeys().length.toString()}', textAlign: TextAlign.end)),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHiveCard(BuildContext context, Box<FieldDebuggable> Function() boxFunc, String boxname) {
    return Card.outlined(
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: GestureDetector(
          onTap: () {
            Navi.push(context, () => DebugHiveBoxPage(boxName: boxname, box: boxFunc()));
          },
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              SizedBox(width: 30, child: Text('')),
              Expanded(child: Text('Hive [$boxname]', style: TextStyle(fontWeight: FontWeight.bold), textAlign: TextAlign.center)),
              SizedBox(width: 40, child: Text('${boxFunc().length.toString()}', textAlign: TextAlign.end)),
            ],
          ),
        ),
      ),
    );
  }
}
