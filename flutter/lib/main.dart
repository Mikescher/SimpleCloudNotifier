import 'package:fl_toast/fl_toast.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/nav_layout.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Hive.initFlutter();
  await Globals().init();

  Hive.registerAdapter(SCNRequestAdapter());
  Hive.registerAdapter(SCNLogAdapter());

  try {
    await Hive.openBox<SCNRequest>('scn-requests');
    await Hive.openBox<SCNLog>('scn-logs');
  } catch (e) {
    print(e);
    Hive.deleteBoxFromDisk('scn-requests');
    Hive.deleteBoxFromDisk('scn-logs');
    await Hive.openBox<SCNRequest>('scn-requests');
    await Hive.openBox<SCNLog>('scn-logs');
  }

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (context) => UserAccount(),
          lazy: false,
        ),
        ChangeNotifierProvider(
          create: (context) => AppTheme(),
          lazy: false,
        ),
      ],
      child: const SCNApp(),
    ),
  );
}

class SCNApp extends StatelessWidget {
  const SCNApp({super.key});

  @override
  Widget build(BuildContext context) {
    Provider.of<UserAccount>(context); // ensure UserAccount is loaded (unneccessary if lazy: false is set in MultiProvider ??)

    return Consumer<AppTheme>(
      builder: (context, appTheme, child) => MaterialApp(
        title: 'SimpleCloudNotifier',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue, brightness: appTheme.darkMode ? Brightness.dark : Brightness.light),
          useMaterial3: true,
        ),
        home: const ToastProvider(child: SCNNavLayout()),
      ),
    );
  }
}
