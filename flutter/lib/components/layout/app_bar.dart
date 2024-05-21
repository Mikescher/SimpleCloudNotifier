import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';

class SCNAppBar extends StatelessWidget implements PreferredSizeWidget {
  const SCNAppBar({Key? key, this.title}) : super(key: key);

  final String? title;

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(title ?? 'Simple Cloud Notifier 2.0'),
      actions: <Widget>[
        Consumer<AppTheme>(
          builder: (context, appTheme, child) => IconButton(
            icon: Icon(appTheme.darkMode ? FontAwesomeIcons.solidSun : FontAwesomeIcons.solidMoon),
            tooltip: 'Debug',
            onPressed: () {
              appTheme.switchDarkMode();
            },
          ),
        ),
        IconButton(
          icon: const Icon(FontAwesomeIcons.solidSpiderBlackWidow),
          tooltip: 'Debug',
          onPressed: () {},
        ),
        IconButton(
          icon: const Icon(FontAwesomeIcons.solidMagnifyingGlass),
          tooltip: 'Search',
          onPressed: () {},
        ),
      ],
      backgroundColor: Theme.of(context).secondaryHeaderColor,
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
