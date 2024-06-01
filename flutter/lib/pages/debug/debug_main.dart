import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/pages/debug/debug_actions.dart';
import 'package:simplecloudnotifier/pages/debug/debug_colors.dart';
import 'package:simplecloudnotifier/pages/debug/debug_logs.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence.dart';
import 'package:simplecloudnotifier/pages/debug/debug_requests.dart';

class DebugMainPage extends StatefulWidget {
  @override
  _DebugMainPageState createState() => _DebugMainPageState();
}

enum DebugMainPageSubPage { colors, requests, persistence, logs, actions }

class _DebugMainPageState extends State<DebugMainPage> {
  final Map<DebugMainPageSubPage, Widget> _subpages = {
    DebugMainPageSubPage.colors: DebugColorsPage(),
    DebugMainPageSubPage.requests: DebugRequestsPage(),
    DebugMainPageSubPage.persistence: DebugPersistencePage(),
    DebugMainPageSubPage.logs: DebugLogsPage(),
    DebugMainPageSubPage.actions: DebugActionsPage(),
  };

  DebugMainPageSubPage _subPage = DebugMainPageSubPage.colors;

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Debug',
      showSearch: false,
      showDebug: false,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                buildSegButton(context),
              ],
            ),
          ),
          Expanded(
            child: _subpages[_subPage]!,
          ),
        ],
      ),
    );
  }

  Widget buildSegButton(BuildContext context) {
    return SegmentedButton<DebugMainPageSubPage>(
      showSelectedIcon: false,
      segments: const <ButtonSegment<DebugMainPageSubPage>>[
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.colors, icon: Icon(FontAwesomeIcons.solidPaintRoller, size: 14)),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.actions, icon: Icon(FontAwesomeIcons.solidHammer, size: 14)),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.requests, icon: Icon(FontAwesomeIcons.solidNetworkWired, size: 14)),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.persistence, icon: Icon(FontAwesomeIcons.solidFloppyDisk, size: 14)),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.logs, icon: Icon(FontAwesomeIcons.solidFileLines, size: 14)),
      ],
      style: ButtonStyle(
        padding: MaterialStateProperty.all<EdgeInsets>(EdgeInsets.fromLTRB(0, 0, 0, 0)),
        visualDensity: VisualDensity(horizontal: -3, vertical: -3),
      ),
      selected: <DebugMainPageSubPage>{_subPage},
      onSelectionChanged: (Set<DebugMainPageSubPage> v) {
        setState(() {
          _subPage = v.first;
        });
      },
    );
  }
}
