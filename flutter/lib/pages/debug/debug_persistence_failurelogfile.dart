import 'dart:io';

import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';

class DebugFailureLogFilePage extends StatefulWidget {
  final String path;

  DebugFailureLogFilePage({required this.path}) {}

  @override
  State<DebugFailureLogFilePage> createState() => _DebugFailureLogFilePageState();
}

class _DebugFailureLogFilePageState extends State<DebugFailureLogFilePage> {
  ImmediateFuture<String>? _futureContent;

  @override
  void initState() {
    super.initState();

    _futureContent = ImmediateFuture.ofFuture(new File(this.widget.path).readAsString());
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'FailureLog',
      showSearch: false,
      child: () {
        if (_futureContent == null) {
          return Center(child: CircularProgressIndicator());
        }

        return FutureBuilder(
          future: _futureContent!.future,
          builder: ((context, snapshot) {
            if (_futureContent?.value != null) {
              return _buildContent(context, _futureContent!.value!);
            } else if (snapshot.connectionState == ConnectionState.done && snapshot.hasError) {
              return Text('Error: ${snapshot.error}'); //TODO better error display
            } else if (snapshot.connectionState == ConnectionState.done) {
              return _buildContent(context, snapshot.data!);
            } else {
              return Center(child: CircularProgressIndicator());
            }
          }),
        );
      }(),
    );
  }

  Widget _buildContent(BuildContext context, String value) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Text(value, style: TextStyle(fontFamily: "monospace")),
      ),
    );
  }
}
