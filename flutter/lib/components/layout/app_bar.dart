import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/components/layout/app_bar_progress_indicator.dart';
import 'package:simplecloudnotifier/pages/debug/debug_main.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class SCNAppBar extends StatelessWidget implements PreferredSizeWidget {
  const SCNAppBar({
    Key? key,
    required this.title,
    required this.showThemeSwitch,
    required this.showSearch,
    required this.showShare,
    this.onShare = null,
  }) : super(key: key);

  final String? title;
  final bool showThemeSwitch;
  final bool showSearch;
  final bool showShare;
  final void Function()? onShare;

  @override
  Widget build(BuildContext context) {
    final cfg = Provider.of<AppSettings>(context);

    var actions = <Widget>[];

    if (cfg.showDebugButton) {
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidSpiderBlackWidow),
        tooltip: 'Debug',
        onPressed: () {
          Navi.push(context, () => DebugMainPage());
        },
      ));
    }

    if (showThemeSwitch) {
      actions.add(Consumer<AppTheme>(
        builder: (context, appTheme, child) => IconButton(
          icon: Icon(appTheme.darkMode ? FontAwesomeIcons.solidSun : FontAwesomeIcons.solidMoon),
          tooltip: appTheme.darkMode ? 'Light mode' : 'Dark mode',
          onPressed: appTheme.switchDarkMode,
        ),
      ));
    } else {
      actions.add(Visibility(
        visible: false,
        maintainSize: true,
        maintainAnimation: true,
        maintainState: true,
        child: IconButton(
          icon: const Icon(FontAwesomeIcons.square),
          onPressed: () {/*TODO*/},
        ),
      ));
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
      actions.add(Visibility(
        visible: false,
        maintainSize: true,
        maintainAnimation: true,
        maintainState: true,
        child: IconButton(
          icon: const Icon(FontAwesomeIcons.square),
          onPressed: () {/*TODO*/},
        ),
      ));
    }

    return AppBar(
      title: Text(title ?? 'Simple Cloud Notifier 2.0'),
      actions: actions,
      backgroundColor: Theme.of(context).secondaryHeaderColor,
      bottom: PreferredSize(
        preferredSize: Size(double.infinity, 1.0),
        child: AppBarProgressIndicator(),
      ),
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
