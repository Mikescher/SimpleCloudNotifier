import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:fl_toast/fl_toast.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/nav_layout.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/state/user_account.dart';
import 'package:firebase_core/firebase_core.dart';
import 'firebase_options.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Hive.initFlutter();
  await Globals().init();

  Hive.registerAdapter(SCNRequestAdapter());
  Hive.registerAdapter(SCNLogAdapter());
  Hive.registerAdapter(SCNLogLevelAdapter());

  try {
    await Hive.openBox<SCNRequest>('scn-requests');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-requests');
    await Hive.openBox<SCNRequest>('scn-requests');
    ApplicationLog.error('Failed to open Hive-Box: scn-requests: ' + exc.toString(), trace: trace);
  }

  try {
    await Hive.openBox<SCNLog>('scn-logs');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-logs');
    await Hive.openBox<SCNLog>('scn-logs');
    ApplicationLog.error('Failed to open Hive-Box: scn-logs: ' + exc.toString(), trace: trace);
  }

  UserAccount(); // ensure UserAccount is loaded

  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  final notificationSettings = await FirebaseMessaging.instance.requestPermission(provisional: true);

  FirebaseMessaging.instance.onTokenRefresh.listen((fcmToken) {
    try {
      setFirebaseToken(fcmToken);
    } catch (exc, trace) {
      ApplicationLog.error('Failed to set firebase token: ' + exc.toString(), trace: trace);
    }
  }).onError((dynamic err) {
    ApplicationLog.error('Failed to listen to token refresh events: ' + (err?.toString() ?? ''));
  });

  ApplicationLog.debug('Application started');

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (context) => UserAccount(), lazy: false),
        ChangeNotifierProvider(create: (context) => AppTheme(), lazy: false),
      ],
      child: const SCNApp(),
    ),
  );
}

void setFirebaseToken(String fcmToken) async {
  ApplicationLog.info('New firebase token: $fcmToken');
  final acc = UserAccount();
  if (acc.auth != null) {
    if (acc.client == null) {
      final client = await APIClient.addClient(acc.auth, fcmToken, Globals().platform, Globals().version, Globals().clientType);
      acc.setClient(client);
    } else {
      final client = await APIClient.updateClient(acc.auth, acc.client!.clientID, fcmToken, Globals().platform, Globals().version);
      acc.setClient(client);
    }
  }
}

class SCNApp extends StatelessWidget {
  const SCNApp({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<AppTheme>(
      builder: (context, appTheme, child) => MaterialApp(
        title: 'SimpleCloudNotifier',
        theme: ThemeData(
          //TODO color settings
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue, brightness: appTheme.darkMode ? Brightness.dark : Brightness.light),
          useMaterial3: true,
        ),
        home: const ToastProvider(child: SCNNavLayout()),
      ),
    );
  }
}
