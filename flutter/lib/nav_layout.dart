import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/components/hidable_fab/hidable_fab.dart';
import 'package:simplecloudnotifier/components/layout/app_bar.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list.dart';
import 'package:simplecloudnotifier/pages/send/send.dart';
import 'package:simplecloudnotifier/components/bottom_fab/fab_bottom_app_bar.dart';
import 'package:simplecloudnotifier/pages/account/account.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list.dart';
import 'package:simplecloudnotifier/pages/settings/root.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
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
    final userAcc = Provider.of<AppAuth>(context, listen: false);
    if (!userAcc.isAuth()) _selectedIndex = 2;

    super.initState();
  }

  void _onItemTapped(int index) {
    final userAcc = Provider.of<AppAuth>(context, listen: false);
    if (!userAcc.isAuth()) {
      Toaster.info("Not logged in", "Please login or create a new account first");
      return;
    }

    setState(() {
      _selectedIndex = index;
    });
  }

  void _onFABTapped() {
    final userAcc = Provider.of<AppAuth>(context, listen: false);
    if (!userAcc.isAuth()) {
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
        showSearch: _selectedIndex == 0,
        showShare: false,
        showThemeSwitch: true,
      ),
      body: IndexedStack(
        children: [
          ExcludeFocus(excluding: _selectedIndex != 0, child: MessageListPage(isVisiblePage: _selectedIndex == 0)),
          ExcludeFocus(excluding: _selectedIndex != 1, child: ChannelRootPage(isVisiblePage: _selectedIndex == 1)),
          ExcludeFocus(excluding: _selectedIndex != 2, child: AccountRootPage(isVisiblePage: _selectedIndex == 2)),
          ExcludeFocus(excluding: _selectedIndex != 3, child: SettingsRootPage(isVisiblePage: _selectedIndex == 3)),
          ExcludeFocus(excluding: _selectedIndex != 4, child: SendRootPage(isVisiblePage: _selectedIndex == 4)),
        ],
        index: _selectedIndex,
      ),
      bottomNavigationBar: _buildNavBar(context),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      floatingActionButton: HidableFAB(
        heroTag: 'fab_main',
        onPressed: _onFABTapped,
        icon: FontAwesomeIcons.solidPaperPlaneTop,
      ),
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
