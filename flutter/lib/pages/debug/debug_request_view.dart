import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

class DebugRequestViewPage extends StatelessWidget {
  final SCNRequest request;

  DebugRequestViewPage({required this.request});

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Request',
      showSearch: false,
      showDebug: false,
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              ...buildRow(context, "Name", request.name),
              ...buildRow(context, "Timestamp (Start)", request.timestampStart.toString()),
              ...buildRow(context, "Timestamp (End)", request.timestampEnd.toString()),
              ...buildRow(context, "Duration", request.timestampEnd.difference(request.timestampStart).toString()),
              Divider(),
              ...buildRow(context, "Method", request.method),
              ...buildRow(context, "URL", request.url),
              if (request.requestHeaders.isNotEmpty) ...buildRow(context, "Request->Headers", request.requestHeaders.entries.map((v) => '${v.key} = ${v.value}').join('\n')),
              if (request.requestBody != '') ...buildRow(context, "Request->Body", request.requestBody),
              Divider(),
              if (request.responseStatusCode != 0) ...buildRow(context, "Response->Statuscode", request.responseStatusCode.toString()),
              if (request.responseBody != '') ...buildRow(context, "Reponse->Body", request.responseBody),
              if (request.responseHeaders.isNotEmpty) ...buildRow(context, "Reponse->Headers", request.responseHeaders.entries.map((v) => '${v.key} = ${v.value}').join('\n')),
              Divider(),
              if (request.error != '') ...buildRow(context, "Error", request.error),
              if (request.stackTrace != '') ...buildRow(context, "Stacktrace", request.stackTrace),
            ],
          ),
        ),
      ),
    );
  }

  List<Widget> buildRow(BuildContext context, String title, String value) {
    return [
      Padding(
        padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 8.0),
        child: Row(
          children: [
            Expanded(
              child: Text(title, style: TextStyle(fontWeight: FontWeight.bold)),
            ),
            IconButton(
              icon: FaIcon(
                FontAwesomeIcons.copy,
              ),
              iconSize: 14,
              padding: EdgeInsets.fromLTRB(0, 0, 4, 0),
              constraints: BoxConstraints(),
              onPressed: () {
                Clipboard.setData(new ClipboardData(text: value));
                Toaster.info("Clipboard", 'Copied text to Clipboard');
              },
            ),
          ],
        ),
      ),
      Card.filled(
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: request.type == 'SUCCESS' ? null : Theme.of(context).colorScheme.errorContainer,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 2.0, horizontal: 6.0),
          child: SelectableText(
            value,
            minLines: 1,
            maxLines: 10,
          ),
        ),
      ),
    ];
  }
}
