import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list.dart';
import 'package:simplecloudnotifier/pages/send/root.dart';
import 'package:simplecloudnotifier/components/bottom_fab/fab_bottom_app_bar.dart';
import 'package:simplecloudnotifier/pages/account/root.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list.dart';
import 'package:simplecloudnotifier/pages/settings/root.dart';

class SCNNavLayout extends StatefulWidget {
  const SCNNavLayout({super.key});

  @override
  State<SCNNavLayout> createState() => _SCNNavLayoutState();
}

class _SCNNavLayoutState extends State<SCNNavLayout> {
  int _selectedIndex = 0; // 4 == FAB

  static const List<Widget> _subPages = <Widget>[
    MessageListPage(),
    ChannelRootPage(),
    AccountRootPage(),
    SettingsRootPage(),
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
      appBar: SCNAppBar(
        title: null,
        showDebug: true,
        showSearch: _selectedIndex == 0 || _selectedIndex == 1,
        showThemeSwitch: true,
      ),
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
}
