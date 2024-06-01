import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:intl/intl.dart';
import 'package:simplecloudnotifier/pages/debug/debug_request_view.dart';
import 'package:simplecloudnotifier/state/request_log.dart';

class DebugRequestsPage extends StatefulWidget {
  @override
  _DebugRequestsPageState createState() => _DebugRequestsPageState();
}

class _DebugRequestsPageState extends State<DebugRequestsPage> {
  Box<SCNRequest> requestsBox = Hive.box<SCNRequest>('scn-requests');

  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');

  @override
  Widget build(BuildContext context) {
    return Container(
      child: ValueListenableBuilder(
        valueListenable: requestsBox.listenable(),
        builder: (context, Box<SCNRequest> box, _) {
          return ListView.builder(
            itemCount: requestsBox.length,
            itemBuilder: (context, listIndex) {
              final req = requestsBox.getAt(requestsBox.length - listIndex - 1)!;
              if (req.type == 'SUCCESS') {
                return buildItemSuccess(context, req);
              } else {
                return buildItemError(context, req);
              }
            },
          );
        },
      ),
    );
  }

  Padding buildItemError(BuildContext context, SCNRequest req) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 2.0),
      child: GestureDetector(
        onTap: () => Navigator.push(context, MaterialPageRoute<DebugRequestViewPage>(builder: (context) => DebugRequestViewPage(request: req))),
        child: ListTile(
            tileColor: Theme.of(context).colorScheme.errorContainer,
            textColor: Theme.of(context).colorScheme.onErrorContainer,
            title: Row(
              children: [
                SizedBox(
                  width: 120,
                  child: Text(_dateFormat.format(req.timestampStart), style: TextStyle(fontSize: 12)),
                ),
                Expanded(
                  child: Text(req.name, style: TextStyle(fontWeight: FontWeight.bold)),
                ),
                SizedBox(width: 2),
                Text('${req.timestampEnd.difference(req.timestampStart).inMilliseconds}ms', style: TextStyle(fontSize: 12)),
              ],
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(req.type),
                Text(
                  req.error,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            )),
      ),
    );
  }

  Padding buildItemSuccess(BuildContext context, SCNRequest req) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 2.0),
      child: GestureDetector(
        onTap: () => Navigator.push(context, MaterialPageRoute<DebugRequestViewPage>(builder: (context) => DebugRequestViewPage(request: req))),
        child: ListTile(
          title: Row(
            children: [
              SizedBox(
                width: 120,
                child: Text(_dateFormat.format(req.timestampStart), style: TextStyle(fontSize: 12)),
              ),
              Expanded(
                child: Text(req.name, style: TextStyle(fontWeight: FontWeight.bold)),
              ),
              SizedBox(width: 2),
              Text('${req.timestampEnd.difference(req.timestampStart).inMilliseconds}ms', style: TextStyle(fontSize: 12)),
            ],
          ),
          subtitle: Text(req.type),
        ),
      ),
    );
  }
}
