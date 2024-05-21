import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';

class SCNScaffold extends StatelessWidget {
  const SCNScaffold({Key? key, required this.child, this.title}) : super(key: key);

  final Widget child;
  final String? title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SCNAppBar(
        title: title,
      ),
      body: child,
    );
  }
}
