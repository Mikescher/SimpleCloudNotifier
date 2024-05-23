import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/pages/debug/debug_main.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';

class SCNAppBar extends StatelessWidget implements PreferredSizeWidget {
  const SCNAppBar({
    Key? key,
    required this.title,
    required this.showThemeSwitch,
    required this.showDebug,
    required this.showSearch,
  }) : super(key: key);

  final String? title;
  final bool showThemeSwitch;
  final bool showDebug;
  final bool showSearch;

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(title ?? 'Simple Cloud Notifier 2.0'),
      actions: <Widget>[
        if (showThemeSwitch)
          Consumer<AppTheme>(
            builder: (context, appTheme, child) => IconButton(
              icon: Icon(appTheme.darkMode ? FontAwesomeIcons.solidSun : FontAwesomeIcons.solidMoon),
              tooltip: 'Debug',
              onPressed: () {
                appTheme.switchDarkMode();
              },
            ),
          ),
        if (!showThemeSwitch) SizedBox.square(dimension: 40),
        if (showDebug)
          IconButton(
            icon: const Icon(FontAwesomeIcons.solidSpiderBlackWidow),
            tooltip: 'Debug',
            onPressed: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => DebugMainPage()));
            },
          ),
        if (!showDebug) SizedBox.square(dimension: 40),
        if (showSearch)
          IconButton(
            icon: const Icon(FontAwesomeIcons.solidMagnifyingGlass),
            tooltip: 'Search',
            onPressed: () {},
          ),
        if (!showSearch) SizedBox.square(dimension: 40),
      ],
      backgroundColor: Theme.of(context).secondaryHeaderColor,
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
