import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'nav_layout.dart';
import 'state/app_theme.dart';
import 'state/user_account.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (context) => UserAccount()),
        ChangeNotifierProvider(create: (context) => AppTheme()),
      ],
      child: const SCNApp(),
    ),
  );
}

class SCNApp extends StatelessWidget {
  const SCNApp({super.key});

  @override
  Widget build(BuildContext context) {
    Provider.of<UserAccount>(context); // ensure UserAccount is loaded

    return Consumer<AppTheme>(
      builder: (context, appTheme, child) => MaterialApp(
        title: 'SimpleCloudNotifier',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue, brightness: appTheme.darkMode ? Brightness.dark : Brightness.light),
          useMaterial3: true,
        ),
        home: const SCNNavLayout(),
      ),
    );
  }
}
