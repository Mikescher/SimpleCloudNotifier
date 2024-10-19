import 'dart:io';

import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/pages/debug/debug_persistence_failurelogfile.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:path/path.dart' as path;
import 'package:simplecloudnotifier/utils/toaster.dart';

class DebugFailureLogsPage extends StatefulWidget {
  final String dir;

  DebugFailureLogsPage({required this.dir});

  @override
  State<DebugFailureLogsPage> createState() => _DebugFailureLogsPageState();
}

class _DebugFailureLogsPageState extends State<DebugFailureLogsPage> {
  List<String> files = [];

  _DebugFailureLogsPageState() {
    files = _listFilesInRawLogFolder();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'F-Logs',
      showSearch: false,
      child: ListView.separated(
        itemCount: files.length,
        itemBuilder: (context, listIndex) {
          return GestureDetector(
            onTap: () {
              Navi.push(context, () => DebugFailureLogFilePage(path: files[listIndex]));
            },
            child: Container(
              padding: EdgeInsets.fromLTRB(8, 0, 8, 0),
              child: Row(
                children: [
                  Expanded(child: Text(path.basename(files[listIndex]), style: TextStyle(fontWeight: FontWeight.bold, fontSize: 12))),
                  IconButton(
                    icon: const Icon(FontAwesomeIcons.trash),
                    tooltip: 'Delete',
                    iconSize: 16,
                    color: Colors.red,
                    onPressed: () => _deleteFile(context, files[listIndex]),
                  )
                ],
              ),
            ),
          );
        },
        separatorBuilder: (context, index) => Divider(),
      ),
    );
  }

  List<String> _listFilesInRawLogFolder() {
    final fse = Globals().rawFailureLogsDir.listSync();

    ApplicationLog.debug("Found ${fse.length} files in raw log folder '${Globals().rawFailureLogsDir.path}'");

    var paths = fse.where((element) => element is File).map((e) => e.path).toList();

    paths.sort((a, b) => -1 * a.compareTo(b));

    return paths;
  }

  void _deleteFile(BuildContext context, String fil) {
    final file = File(fil);

    file.deleteSync();

    setState(() {
      files = _listFilesInRawLogFolder();
    });

    Toaster.info("Okay", "File deleted");
  }
}
