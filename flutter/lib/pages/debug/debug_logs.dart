import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:intl/intl.dart';
import 'package:simplecloudnotifier/state/application_log.dart';

class DebugLogsPage extends StatefulWidget {
  @override
  _DebugLogsPageState createState() => _DebugLogsPageState();
}

class _DebugLogsPageState extends State<DebugLogsPage> {
  Box<SCNLog> logBox = Hive.box<SCNLog>('scn-logs');

  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  @override
  Widget build(BuildContext context) {
    return Container(
      child: ValueListenableBuilder(
        valueListenable: logBox.listenable(),
        builder: (context, Box<SCNLog> box, _) {
          return ListView.builder(
            itemCount: logBox.length,
            itemBuilder: (context, listIndex) {
              final log = logBox.getAt(logBox.length - listIndex - 1)!;
              switch (log.level) {
                case SCNLogLevel.debug:
                  return buildItem(context, log, Theme.of(context).hintColor);
                case SCNLogLevel.info:
                  return buildItem(context, log, Colors.blueAccent);
                case SCNLogLevel.warning:
                  return buildItem(context, log, Colors.orangeAccent);
                case SCNLogLevel.error:
                  return buildItem(context, log, Colors.redAccent);
                case SCNLogLevel.fatal:
                  return buildItem(context, log, Colors.black);
              }
            },
          );
        },
      ),
    );
  }

  Widget buildItem(BuildContext context, SCNLog log, Color tagColor) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 2.0),
      child: Card.filled(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.fromLTRB(12, 1, 12, 1),
                  decoration: BoxDecoration(
                    color: tagColor,
                    borderRadius: BorderRadius.only(topLeft: Radius.circular(8)),
                  ),
                  child: Text(
                    log.level.name.toUpperCase(),
                    style: TextStyle(fontWeight: FontWeight.bold, color: Theme.of(context).cardColor, fontSize: 14),
                  ),
                ),
                Expanded(child: SizedBox()),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 0, 8, 0),
                  child: Text(_dateFormat.format(log.timestamp), style: TextStyle(fontSize: 12)),
                ),
              ],
            ),
            SizedBox(height: 4),
            if (log.message.isNotEmpty)
              Padding(
                padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
                child: Text(log.message, style: TextStyle(fontWeight: FontWeight.bold)),
              ),
            if (log.additional.isNotEmpty)
              Padding(
                padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
                child: SelectableText(
                  log.additional,
                  style: TextStyle(fontSize: 12),
                  minLines: 1,
                  maxLines: 10,
                ),
              ),
            if (log.trace.isNotEmpty)
              Padding(
                padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
                child: SelectableText(
                  log.trace,
                  style: TextStyle(fontSize: 12),
                  minLines: 1,
                  maxLines: 10,
                ),
              ),
            SizedBox(height: 8),
          ],
        ),
      ),
    );
  }
}
