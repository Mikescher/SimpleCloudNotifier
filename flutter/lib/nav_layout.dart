import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/pages/send/root.dart';

import 'bottom_fab/fab_bottom_app_bar.dart';
import 'pages/account/root.dart';
import 'pages/message_list/message_list.dart';
import 'state/app_theme.dart';

class SCNNavLayout extends StatefulWidget {
  const SCNNavLayout({super.key});

  @override
  State<SCNNavLayout> createState() => _SCNNavLayoutState();
}

class _SCNNavLayoutState extends State<SCNNavLayout> {
  int _selectedIndex = 0; // 4 == FAB

  static const List<Widget> _subPages = <Widget>[
    MessageListPage(title: 'Messages'),
    MessageListPage(title: 'Page 2'),
    AccountRootPage(),
    MessageListPage(title: 'Page 4'),
    SendRootPage(),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  void _onFABTapped() {
    setState(() {
      _selectedIndex = 4;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(context),
      body: _subPages.elementAt(_selectedIndex),
      bottomNavigationBar: _buildNavBar(context),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      floatingActionButton: _buildFAB(context),
    );
  }

  Widget _buildFAB(BuildContext context) {
    return FloatingActionButton(
      onPressed: _onFABTapped,
      tooltip: 'Increment',
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(17))),
      elevation: 2.0,
      child: const Icon(FontAwesomeIcons.solidPaperPlaneTop),
    );
  }

  Widget _buildNavBar(BuildContext context) {
    return FABBottomAppBar(
      selectedIndex: _selectedIndex,
      onTabSelected: _onItemTapped,
      color: Theme.of(context).disabledColor,
      selectedColor: Theme.of(context).primaryColorDark,
      notchedShape: const AutomaticNotchedShape(
        RoundedRectangleBorder(
          borderRadius: BorderRadius.only(
            topLeft: Radius.circular(15),
            topRight: Radius.circular(15),
          ),
        ),
        RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(17)),
        ),
      ),
      items: [
        FABBottomAppBarItem(iconData: FontAwesomeIcons.solidEnvelope, text: 'Messages'),
        FABBottomAppBarItem(iconData: FontAwesomeIcons.solidSnake, text: 'Channels'),
        FABBottomAppBarItem(iconData: FontAwesomeIcons.solidFileUser, text: 'Account'),
        FABBottomAppBarItem(iconData: FontAwesomeIcons.solidGear, text: 'Settings'),
      ],
    );
  }

  PreferredSizeWidget _buildAppBar(BuildContext context) {
    return AppBar(
      title: const Text('Simple Cloud Notifier 2.0'),
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
}
