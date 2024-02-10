import 'package:flutter/material.dart';

import 'nav_layout.dart';

void main() {
  runApp(const SCNApp());
}

class SCNApp extends StatelessWidget {
  const SCNApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'SimpleCloudNotifier',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      home: const SCNNavLayout(),
    );
  }
}
