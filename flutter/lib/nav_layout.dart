import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

import 'bottom_fab/fab_bottom_app_bar.dart';
import 'pages/message_list/message_list.dart';

class SCNNavLayout extends StatefulWidget {
  const SCNNavLayout({super.key});

  @override
  State<SCNNavLayout> createState() => _SCNNavLayoutState();
}

class _SCNNavLayoutState extends State<SCNNavLayout> {
  int _selectedIndex = 0;

  static const TextStyle optionStyle = TextStyle(fontSize: 30, fontWeight: FontWeight.bold);

  static const List<Widget> _subPages = <Widget>[
    MessageListPage(title: 'Messages 1'),
    MessageListPage(title: 'Messages 2'),
    MessageListPage(title: 'Messages 3'),
    MessageListPage(title: 'Messages 4'),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  void _onFABTapped(int index) {
    setState(() {
      //TODO
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(context),
      body: Center(
        child: _subPages.elementAt(_selectedIndex),
      ),
      bottomNavigationBar: _buildNavBar(context),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      floatingActionButton: _buildFAB(context),
    );
  }

  Widget _buildFAB(BuildContext context) {
    return FloatingActionButton(
      onPressed: () {},
      tooltip: 'Increment',
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(17))),
      elevation: 2.0,
      child: const Icon(FontAwesomeIcons.solidPaperPlaneTop),
    );
  }

  Widget _buildNavBar(BuildContext context) {
    return FABBottomAppBar(
      onTabSelected: _onItemTapped,
      color: Colors.grey,
      selectedColor: Colors.black,
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

  AppBar _buildAppBar(BuildContext context) {
    return AppBar(
      title: const Text('Simple Cloud Notifier 2.0'),
      actions: <Widget>[
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
        IconButton(
          icon: const Icon(FontAwesomeIcons.solidQrcode),
          tooltip: 'Show Account QR Code',
          onPressed: () {},
        ),
      ],
      backgroundColor: Theme.of(context).secondaryHeaderColor,
    );
  }
}
