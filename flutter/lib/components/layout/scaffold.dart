import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';

class SCNScaffold extends StatelessWidget {
  const SCNScaffold({
    Key? key,
    required this.child,
    this.title,
    this.showThemeSwitch = true,
    this.showDebug = true,
    this.showSearch = true,
  }) : super(key: key);

  final Widget child;
  final String? title;
  final bool showThemeSwitch;
  final bool showDebug;
  final bool showSearch;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SCNAppBar(
        title: title,
        showThemeSwitch: showThemeSwitch,
        showDebug: showDebug,
        showSearch: showSearch,
      ),
      body: child,
    );
  }
}
