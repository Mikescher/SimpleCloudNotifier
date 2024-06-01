import 'package:flutter/material.dart';
import 'package:flutter_lazy_indexed_stack/flutter_lazy_indexed_stack.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list.dart';
import 'package:simplecloudnotifier/pages/send/root.dart';
import 'package:simplecloudnotifier/components/bottom_fab/fab_bottom_app_bar.dart';
import 'package:simplecloudnotifier/pages/account/account.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list.dart';
import 'package:simplecloudnotifier/pages/settings/root.dart';
import 'package:simplecloudnotifier/state/user_account.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

class SCNNavLayout extends StatefulWidget {
  const SCNNavLayout({super.key});

  @override
  State<SCNNavLayout> createState() => _SCNNavLayoutState();
}

class _SCNNavLayoutState extends State<SCNNavLayout> {
  int _selectedIndex = 0; // 4 == FAB

  @override
  initState() {
    final userAcc = Provider.of<UserAccount>(context, listen: false);
    if (userAcc.auth == null) _selectedIndex = 2;

    super.initState();
  }

  void _onItemTapped(int index) {
    final userAcc = Provider.of<UserAccount>(context, listen: false);
    if (userAcc.auth == null) {
      Toaster.info("Not logged in", "Please login or create a new account first");
      return;
    }

    setState(() {
      _selectedIndex = index;
    });
  }

  void _onFABTapped() {
    final userAcc = Provider.of<UserAccount>(context, listen: false);
    if (userAcc.auth == null) {
      Toaster.info("Not logged in", "Please login or create a new account first");
      return;
    }

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
      body: LazyIndexedStack(
        children: [
          MessageListPage(),
          ChannelRootPage(),
          AccountRootPage(),
          SettingsRootPage(),
          SendRootPage(),
        ],
        index: _selectedIndex,
      ),
      bottomNavigationBar: _buildNavBar(context),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      floatingActionButton: _buildFAB(context),
    );
  }

  Widget _buildFAB(BuildContext context) {
    return FloatingActionButton(
      onPressed: _onFABTapped,
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
