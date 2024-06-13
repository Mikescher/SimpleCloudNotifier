import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/pages/debug/debug_main.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class SCNAppBar extends StatelessWidget implements PreferredSizeWidget {
  const SCNAppBar({
    Key? key,
    required this.title,
    required this.showThemeSwitch,
    required this.showDebug,
    required this.showSearch,
    required this.showShare,
    this.onShare = null,
  }) : super(key: key);

  final String? title;
  final bool showThemeSwitch;
  final bool showDebug;
  final bool showSearch;
  final bool showShare;
  final void Function()? onShare;

  @override
  Widget build(BuildContext context) {
    var actions = <Widget>[];

    if (showThemeSwitch) {
      actions.add(Consumer<AppTheme>(
        builder: (context, appTheme, child) => IconButton(
          icon: Icon(appTheme.darkMode ? FontAwesomeIcons.solidSun : FontAwesomeIcons.solidMoon),
          tooltip: appTheme.darkMode ? 'Light mode' : 'Dark mode',
          onPressed: appTheme.switchDarkMode,
        ),
      ));
    } else {
      actions.add(SizedBox.square(dimension: 40));
    }

    if (showDebug) {
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidSpiderBlackWidow),
        tooltip: 'Debug',
        onPressed: () {
          Navi.push(context, () => DebugMainPage());
        },
      ));
    } else {
      actions.add(SizedBox.square(dimension: 40));
    }

    if (showSearch) {
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidMagnifyingGlass),
        tooltip: 'Search',
        onPressed: () {/*TODO*/},
      ));
    } else if (showShare) {
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidShareNodes),
        tooltip: 'Share',
        onPressed: onShare ?? () {},
      ));
    } else {
      actions.add(SizedBox.square(dimension: 40));
    }

    return AppBar(
      title: Text(title ?? 'Simple Cloud Notifier 2.0'),
      actions: actions,
      backgroundColor: Theme.of(context).secondaryHeaderColor,
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
