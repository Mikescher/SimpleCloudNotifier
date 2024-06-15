import 'dart:io';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/client.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/nav_layout.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/app_theme.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/fb_message.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/request_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:toastification/toastification.dart';
import 'firebase_options.dart';

void main() async {
  print('[INIT] Application starting...');

  print('[INIT] Ensure WidgetsFlutterBinding...');

  WidgetsFlutterBinding.ensureInitialized();

  print('[INIT] Init Globals...');

  await Globals().init();

  print('[INIT] Init Hive...');

  await Hive.initFlutter();

  Hive.registerAdapter(SCNRequestAdapter());
  Hive.registerAdapter(SCNLogAdapter());
  Hive.registerAdapter(SCNLogLevelAdapter());
  Hive.registerAdapter(MessageAdapter());
  Hive.registerAdapter(ChannelAdapter());
  Hive.registerAdapter(FBMessageAdapter());

  print('[INIT] Load Hive<scn-requests>...');

  try {
    await Hive.openBox<SCNRequest>('scn-requests');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-requests');
    await Hive.openBox<SCNRequest>('scn-requests');
    ApplicationLog.error('Failed to open Hive-Box: scn-requests: ' + exc.toString(), trace: trace);
  }

  print('[INIT] Load Hive<scn-logs>...');

  try {
    await Hive.openBox<SCNLog>('scn-logs');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-logs');
    await Hive.openBox<SCNLog>('scn-logs');
    ApplicationLog.error('Failed to open Hive-Box: scn-logs: ' + exc.toString(), trace: trace);
  }

  print('[INIT] Load Hive<scn-message-cache>...');

  try {
    await Hive.openBox<Message>('scn-message-cache');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-message-cache');
    await Hive.openBox<Message>('scn-message-cache');
    ApplicationLog.error('Failed to open Hive-Box: scn-message-cache' + exc.toString(), trace: trace);
  }

  print('[INIT] Load Hive<scn-channel-cache>...');

  try {
    await Hive.openBox<Channel>('scn-channel-cache');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-channel-cache');
    await Hive.openBox<Channel>('scn-channel-cache');
    ApplicationLog.error('Failed to open Hive-Box: scn-channel-cache' + exc.toString(), trace: trace);
  }

  print('[INIT] Load Hive<scn-fb-messages>...');

  try {
    await Hive.openBox<FBMessage>('scn-fb-messages');
  } catch (exc, trace) {
    Hive.deleteBoxFromDisk('scn-fb-messages');
    await Hive.openBox<FBMessage>('scn-fb-messages');
    ApplicationLog.error('Failed to open Hive-Box: scn-fb-messages' + exc.toString(), trace: trace);
  }

  print('[INIT] Load AppAuth...');

  final appAuth = AppAuth(); // ensure UserAccount is loaded

  if (appAuth.isAuth()) {
    try {
      print('[INIT] Load User...');
      await appAuth.loadUser();
      //TODO fallback to cached user (perhaps event use cached client (if exists) directly and only update/load in background)
    } catch (exc, trace) {
      ApplicationLog.error('Failed to load user (on startup): ' + exc.toString(), trace: trace);
    }
    try {
      print('[INIT] Load Client...');
      await appAuth.loadClient();
      //TODO fallback to cached client (perhaps event use cached client (if exists) directly and only update/load in background)
    } catch (exc, trace) {
      ApplicationLog.error('Failed to load user (on startup): ' + exc.toString(), trace: trace);
    }
  }

  if (!Platform.isLinux) {
    print('[INIT] Init Firebase...');
    await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

    print('[INIT] Request Notification permissions...');
    await FirebaseMessaging.instance.requestPermission(provisional: true);

    FirebaseMessaging.instance.onTokenRefresh.listen((fcmToken) {
      try {
        setFirebaseToken(fcmToken);
      } catch (exc, trace) {
        ApplicationLog.error('Failed to set firebase token: ' + exc.toString(), trace: trace);
      }
    }).onError((dynamic err) {
      ApplicationLog.error('Failed to listen to token refresh events: ' + (err?.toString() ?? ''));
    });

    try {
      print('[INIT] Query firebase token...');
      final fcmToken = await FirebaseMessaging.instance.getToken();
      if (fcmToken != null) {
        setFirebaseToken(fcmToken);
      }
    } catch (exc, trace) {
      ApplicationLog.error('Failed to get+set firebase token: ' + exc.toString(), trace: trace);
    }

    FirebaseMessaging.onBackgroundMessage(_onBackgroundMessage);
    FirebaseMessaging.onMessage.listen(_onForegroundMessage);
  } else {
    print('[INIT] Skip Firebase init (Platform == Linux)...');
  }

  ApplicationLog.debug('[INIT] Application started');

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (context) => AppAuth(), lazy: false),
        ChangeNotifierProvider(create: (context) => AppTheme(), lazy: false),
        ChangeNotifierProvider(create: (context) => AppBarState(), lazy: false),
      ],
      child: SCNApp(),
    ),
  );
}

class SCNApp extends StatelessWidget {
  SCNApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ToastificationWrapper(
      config: ToastificationConfig(
        itemWidth: 440,
        marginBuilder: (alignment) => EdgeInsets.symmetric(vertical: 64),
        animationDuration: Duration(milliseconds: 200),
      ),
      child: Consumer<AppTheme>(
        builder: (context, appTheme, child) => MaterialApp(
          title: 'SimpleCloudNotifier',
          navigatorObservers: [Navi.routeObserver, Navi.modalRouteObserver],
          theme: ThemeData(
            //TODO color settings
            colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue, brightness: appTheme.darkMode ? Brightness.dark : Brightness.light),
            useMaterial3: true,
          ),
          home: SCNNavLayout(),
        ),
      ),
    );
  }
}

void setFirebaseToken(String fcmToken) async {
  final acc = AppAuth();

  final oldToken = Globals().getPrefFCMToken();

  await Globals().setPrefFCMToken(fcmToken);

  ApplicationLog.info('New firebase token received', additional: 'Token: $fcmToken (old: $oldToken)');

  if (!acc.isAuth()) return;

  Client? client;
  try {
    client = await acc.loadClient(forceIfOlder: Duration(seconds: 60));
  } catch (exc, trace) {
    ApplicationLog.error('Failed to get client: ' + exc.toString(), trace: trace);
    return;
  }

  if (oldToken != null && oldToken == fcmToken && client != null && client.fcmToken == fcmToken) {
    ApplicationLog.info('Firebase token unchanged - do nothing', additional: 'Token: $fcmToken');
    return;
  }

  if (client == null) {
    // should not really happen - perhaps someone externally deleted the client?
    final newClient = await APIClient.addClient(acc, fcmToken, Globals().deviceModel, Globals().version, Globals().hostname, Globals().clientType);
    acc.setClientAndClientID(newClient);
    await acc.save();
  } else {
    final newClient = await APIClient.updateClient(acc, client.clientID, fcmToken, Globals().deviceModel, Globals().hostname, Globals().version);
    acc.setClientAndClientID(newClient);
    await acc.save();
  }
}

Future<void> _onBackgroundMessage(RemoteMessage message) async {
  await _receiveMessage(message, false);
}

void _onForegroundMessage(RemoteMessage message) {
  _receiveMessage(message, true);
}

Future<void> _receiveMessage(RemoteMessage message, bool foreground) async {
  // ensure init
  Hive.openBox<SCNLog>('scn-logs');

  ApplicationLog.info('Received FB message (${foreground ? 'foreground' : 'background'}): ${message.messageId ?? 'NULL'}');
  FBMessageLog.insert(message);
}
