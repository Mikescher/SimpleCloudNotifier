import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/components/layout/app_bar_filter_dialog.dart';
import 'package:simplecloudnotifier/components/layout/app_bar_progress_indicator.dart';
import 'package:simplecloudnotifier/pages/debug/debug_main.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class SCNAppBar extends StatefulWidget implements PreferredSizeWidget {
  SCNAppBar({
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
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  State<SCNAppBar> createState() => _SCNAppBarState();
}

class _SCNAppBarState extends State<SCNAppBar> {
  final TextEditingController _ctrlSearchField = TextEditingController();

  @override
  void dispose() {
    _ctrlSearchField.dispose();
    super.dispose();
  }

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

    if (widget.showThemeSwitch) {
      actions.add(Consumer<AppTheme>(
        builder: (context, appTheme, child) => IconButton(
          icon: Icon(appTheme.darkMode ? FontAwesomeIcons.solidSun : FontAwesomeIcons.solidMoon),
          tooltip: appTheme.darkMode ? 'Light mode' : 'Dark mode',
          onPressed: appTheme.switchDarkMode,
        ),
      ));
    } else {
      actions.add(_buildSpacer());
    }

    if (widget.showSearch) {
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidFilter),
        tooltip: 'Filter',
        onPressed: () => _showFilterDialog(context),
      ));
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidMagnifyingGlass),
        tooltip: 'Search',
        onPressed: () => AppBarState().setShowSearchField(true),
      ));
    } else if (widget.showShare) {
      actions.add(_buildSpacer());
      actions.add(IconButton(
        icon: const Icon(FontAwesomeIcons.solidShareNodes),
        tooltip: 'Share',
        onPressed: widget.onShare ?? () {},
      ));
    } else {
      actions.add(_buildSpacer());
    }

    return Consumer<AppBarState>(builder: (context, value, child) {
      if (value.showSearchField) {
        return AppBar(
          leading: IconButton(
            icon: const Icon(FontAwesomeIcons.solidArrowLeft),
            onPressed: () {
              value.setShowSearchField(false);
            },
          ),
          title: _buildSearchTextField(context),
          actions: [
            IconButton(
              icon: const Icon(FontAwesomeIcons.solidMagnifyingGlass),
              onPressed: () {
                value.setShowSearchField(false);
                AppEvents().notifySearchListeners(_ctrlSearchField.text);
                _ctrlSearchField.clear();
              },
            ),
          ],
          backgroundColor: Theme.of(context).secondaryHeaderColor,
          bottom: PreferredSize(
            preferredSize: Size(double.infinity, 1.0),
            child: AppBarProgressIndicator(show: value.loadingIndeterminate),
          ),
        );
      } else {
        return AppBar(
          title: Text(widget.title ?? 'SCN'),
          actions: actions,
          backgroundColor: Theme.of(context).secondaryHeaderColor,
          bottom: PreferredSize(
            preferredSize: Size(double.infinity, 1.0),
            child: AppBarProgressIndicator(show: value.loadingIndeterminate),
          ),
        );
      }
    });
  }

  Visibility _buildSpacer() {
    return Visibility(
      visible: false,
      maintainSize: true,
      maintainAnimation: true,
      maintainState: true,
      child: IconButton(
        icon: const Icon(FontAwesomeIcons.square),
        onPressed: () {/* NO-OP */},
      ),
    );
  }

  Widget _buildSearchTextField(BuildContext context) {
    return TextField(
      controller: _ctrlSearchField,
      autofocus: true,
      style: TextStyle(fontSize: 20),
      textInputAction: TextInputAction.search,
      decoration: InputDecoration(
        hintText: 'Search',
      ),
      onSubmitted: (value) {
        AppBarState().setShowSearchField(false);
        AppEvents().notifySearchListeners(_ctrlSearchField.text);
        _ctrlSearchField.clear();
      },
    );
  }

  void _showFilterDialog(BuildContext context) {
    showDialog<void>(
      context: context,
      barrierDismissible: true,
      barrierColor: Colors.transparent,
      builder: (BuildContext context) {
        return Dialog(
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(0)),
          alignment: Alignment.topCenter,
          insetPadding: EdgeInsets.fromLTRB(0, this.widget.preferredSize.height, 0, 0),
          backgroundColor: Colors.transparent,
          child: AppBarFilterDialog(),
        );
      },
    );
  }
}
