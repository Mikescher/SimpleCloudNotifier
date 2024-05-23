import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/pages/debug/debug_colors.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence.dart';
import 'package:simplecloudnotifier/pages/debug/debug_requests.dart';

class DebugMainPage extends StatefulWidget {
  @override
  _DebugMainPageState createState() => _DebugMainPageState();
}

enum DebugMainPageSubPage { colors, requests, persistence }

class _DebugMainPageState extends State<DebugMainPage> {
  final Map<DebugMainPageSubPage, Widget> _subpages = {
    DebugMainPageSubPage.colors: DebugColorsPage(),
    DebugMainPageSubPage.requests: DebugRequestsPage(),
    DebugMainPageSubPage.persistence: DebugPersistencePage(),
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
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.colors, label: Text('Theme')),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.requests, label: Text('Requests')),
        ButtonSegment<DebugMainPageSubPage>(value: DebugMainPageSubPage.persistence, label: Text('Persistence')),
      ],
      selected: <DebugMainPageSubPage>{_subPage},
      onSelectionChanged: (Set<DebugMainPageSubPage> v) {
        setState(() {
          _subPage = v.first;
        });
      },
    );
  }
}
