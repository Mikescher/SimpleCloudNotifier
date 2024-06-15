import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';

class SCNScaffold extends StatelessWidget {
  const SCNScaffold({
    Key? key,
    required this.child,
    this.title,
    this.showThemeSwitch = true,
    this.showSearch = true,
    this.showShare = false,
    this.onShare = null,
  }) : super(key: key);

  final Widget child;
  final String? title;
  final bool showThemeSwitch;
  final bool showSearch;
  final bool showShare;
  final void Function()? onShare;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SCNAppBar(
        title: title,
        showThemeSwitch: showThemeSwitch,
        showSearch: showSearch,
        showShare: showShare,
        onShare: onShare ?? () {},
      ),
      body: child,
    );
  }
}
