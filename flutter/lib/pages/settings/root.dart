import 'package:flutter/material.dart';

class SettingsRootPage extends StatefulWidget {
  const SettingsRootPage({super.key});

  @override
  State<SettingsRootPage> createState() => _SettingsRootPageState();
}

class _SettingsRootPageState extends State<SettingsRootPage> {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Text('Settings'),
    );
  }
}
