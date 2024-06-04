import 'dart:io';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/pages/account/login.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:uuid/uuid.dart';

class AccountRootPage extends StatefulWidget {
  const AccountRootPage({super.key});

  @override
  State<AccountRootPage> createState() => _AccountRootPageState();
}

class _AccountRootPageState extends State<AccountRootPage> {
  late Future<int>? futureSubscriptionCount;
  late Future<int>? futureClientCount;
  late Future<int>? futureKeyCount;
  late Future<int>? futureChannelAllCount;
  late Future<int>? futureChannelSubscribedCount;

  late AppAuth userAcc;

  bool loading = false;

  @override
  void initState() {
    super.initState();

    userAcc = Provider.of<AppAuth>(context, listen: false);
    userAcc.addListener(_onAuthStateChanged);
    _onAuthStateChanged();
  }

  @override
  void dispose() {
    userAcc.removeListener(_onAuthStateChanged);
    super.dispose();
  }

  void _onAuthStateChanged() {
    futureSubscriptionCount = null;
    futureClientCount = null;
    futureKeyCount = null;
    futureChannelAllCount = null;
    futureChannelSubscribedCount = null;

    if (userAcc.isAuth()) {
      futureChannelAllCount = () async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc, ChannelSelector.all);
        return channels.length;
      }();

      futureChannelSubscribedCount = () async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc, ChannelSelector.subscribed);
        return channels.length;
      }();

      futureSubscriptionCount = () async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final subs = await APIClient.getSubscriptionList(userAcc);
        return subs.length;
      }();

      futureClientCount = () async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final clients = await APIClient.getClientList(userAcc);
        return clients.length;
      }();

      futureKeyCount = () async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final keys = await APIClient.getKeyTokenList(userAcc);
        return keys.length;
      }();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<AppAuth>(
      builder: (context, acc, child) {
        if (!userAcc.isAuth()) {
          return _buildNoAuth(context);
        } else {
          return FutureBuilder(
            future: acc.loadUser(force: false),
            builder: ((context, snapshot) {
              if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasError) {
                  return Text('Error: ${snapshot.error}'); //TODO better error display
                }
                return _buildShowAccount(context, acc, snapshot.data!);
              }
              return Center(child: CircularProgressIndicator());
            }),
          );
        }
      },
    );
  }

  Widget _buildNoAuth(BuildContext context) {
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 32, 24, 16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            if (!loading)
              Center(
                child: Container(
                  width: 200,
                  height: 200,
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.secondary,
                    borderRadius: BorderRadius.circular(100),
                  ),
                  child: Center(child: FaIcon(FontAwesomeIcons.userSecret, size: 96, color: Theme.of(context).colorScheme.onSecondary)),
                ),
              ),
            if (loading)
              Center(
                child: Container(
                  width: 200,
                  height: 200,
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.secondary,
                    borderRadius: BorderRadius.circular(100),
                  ),
                  child: Center(child: CircularProgressIndicator(color: Theme.of(context).colorScheme.onSecondary)),
                ),
              ),
            const SizedBox(height: 32),
            FilledButton(
              style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 24), padding: const EdgeInsets.fromLTRB(8, 12, 8, 12)),
              onPressed: () {
                if (loading) return;
                _createNewAccount();
              },
              child: const Text('Create new account'),
            ),
            const SizedBox(height: 16),
            FilledButton.tonal(
              style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 24), padding: const EdgeInsets.fromLTRB(8, 12, 8, 12)),
              onPressed: () {
                if (loading) return;
                Navigator.push(context, MaterialPageRoute<AccountLoginPage>(builder: (context) => AccountLoginPage()));
              },
              child: const Text('Use existing account'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildShowAccount(BuildContext context, AppAuth acc, User user) {
    //TODO better layout
    return Column(
      children: [
        SingleChildScrollView(
          scrollDirection: Axis.vertical,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 24.0, 8.0, 8.0),
            child: Column(
              children: [
                _buildHeader(context, user),
                const SizedBox(height: 16),
                Text(user.username ?? user.userID, overflow: TextOverflow.ellipsis, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
                const SizedBox(height: 16),
                ..._buildCards(context, user),
              ],
            ),
          ),
        ),
        const Expanded(child: SizedBox(height: 16)),
        _buildFooter(context, user),
        SizedBox(height: 40)
      ],
    );
  }

  Row _buildHeader(BuildContext context, User user) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 80,
          height: 80,
          child: Stack(
            children: [
              Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: Colors.grey,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Center(child: FaIcon(FontAwesomeIcons.addressCard, size: 55, color: Colors.white))),
              if (user.isPro)
                Align(
                  alignment: Alignment.bottomRight,
                  child: Container(
                    child: Text('PRO', style: TextStyle(fontSize: 14, color: Colors.white, fontWeight: FontWeight.bold)),
                    padding: const EdgeInsets.fromLTRB(4, 1, 4, 1),
                    decoration: BoxDecoration(
                      color: Colors.blue,
                      borderRadius: BorderRadius.only(topLeft: Radius.circular(4)),
                    ),
                  ),
                ),
              if (!user.isPro)
                Align(
                  alignment: Alignment.bottomRight,
                  child: Container(
                    child: Text('FREE', style: TextStyle(fontSize: 14, color: Colors.white, fontWeight: FontWeight.bold)),
                    padding: const EdgeInsets.fromLTRB(4, 1, 4, 1),
                    decoration: BoxDecoration(
                      color: Colors.purple,
                      borderRadius: BorderRadius.only(topLeft: Radius.circular(4)),
                    ),
                  ),
                ),
            ],
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(user.username ?? user.userID, overflow: TextOverflow.ellipsis),
              const SizedBox(height: 4),
              Row(
                children: [
                  SizedBox(width: 80, child: Text("Quota", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)))),
                  Expanded(child: Text('${user.quotaUsed} / ${user.quotaPerDay}')),
                ],
              ),
              Row(
                children: [
                  SizedBox(width: 80, child: Text("Messages", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)))),
                  Expanded(child: Text('${user.messagesSent}')),
                ],
              ),
              Row(
                children: [
                  SizedBox(width: 80, child: Text("Channels", style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)))),
                  FutureBuilder(
                    future: futureChannelAllCount,
                    builder: (context, snapshot) {
                      if (snapshot.connectionState == ConnectionState.done) {
                        return Text('${snapshot.data}');
                      }
                      return const SizedBox(width: 8, height: 8, child: Center(child: CircularProgressIndicator()));
                    },
                  )
                ],
              ),
            ],
          ),
        ),
        Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            IconButton(
              icon: FaIcon(FontAwesomeIcons.pen),
              iconSize: 18,
              padding: EdgeInsets.all(4),
              constraints: BoxConstraints(),
              style: ButtonStyle(tapTargetSize: MaterialTapTargetSize.shrinkWrap),
              onPressed: () {/*TODO*/},
            ),
            const SizedBox(height: 4),
            if (!user.isPro)
              IconButton(
                icon: FaIcon(FontAwesomeIcons.cartCircleArrowUp),
                iconSize: 18,
                padding: EdgeInsets.all(4),
                constraints: BoxConstraints(),
                style: ButtonStyle(tapTargetSize: MaterialTapTargetSize.shrinkWrap),
                onPressed: () {/*TODO*/},
              ),
          ],
        ),
      ],
    );
  }

  List<Widget> _buildCards(BuildContext context, User user) {
    return [
      Card.filled(
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: Theme.of(context).cardTheme.color,
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: () {/*TODO*/},
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                FutureBuilder(
                  future: futureSubscriptionCount,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.done) {
                      return Text('${snapshot.data}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
                    }
                    return const SizedBox(width: 12, height: 12, child: Center(child: CircularProgressIndicator()));
                  },
                ),
                const SizedBox(width: 12),
                Text('Subscriptions', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
              ],
            ),
          ),
        ),
      ),
      Card.filled(
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: Theme.of(context).cardTheme.color,
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: () {/*TODO*/},
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                FutureBuilder(
                  future: futureClientCount,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.done) {
                      return Text('${snapshot.data}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
                    }
                    return const SizedBox(width: 12, height: 12, child: Center(child: CircularProgressIndicator()));
                  },
                ),
                const SizedBox(width: 12),
                Text('Clients', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
              ],
            ),
          ),
        ),
      ),
      Card.filled(
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: Theme.of(context).cardTheme.color,
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: () {/*TODO*/},
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                FutureBuilder(
                  future: futureKeyCount,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.done) {
                      return Text('${snapshot.data}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
                    }
                    return const SizedBox(width: 12, height: 12, child: Center(child: CircularProgressIndicator()));
                  },
                ),
                const SizedBox(width: 12),
                Text('Keys', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
              ],
            ),
          ),
        ),
      ),
      Card.filled(
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: Theme.of(context).cardTheme.color,
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: () {/*TODO*/},
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                FutureBuilder(
                  future: futureChannelSubscribedCount,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.done) {
                      return Text('${snapshot.data}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
                    }
                    return const SizedBox(width: 12, height: 12, child: Center(child: CircularProgressIndicator()));
                  },
                ),
                const SizedBox(width: 12),
                Text('Channels', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
              ],
            ),
          ),
        ),
      ),
      Card.filled(
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
        color: Theme.of(context).cardTheme.color,
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: () {/*TODO*/},
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Text('${user.messagesSent}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
                const SizedBox(width: 12),
                Text('Messages', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
              ],
            ),
          ),
        ),
      ),
    ];
  }

  Widget _buildFooter(BuildContext context, User user) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
      child: Row(
        children: [
          Expanded(child: FilledButton(onPressed: _logout, child: Text('Logout'), style: TextButton.styleFrom(backgroundColor: Colors.orange))),
          const SizedBox(width: 8),
          Expanded(child: FilledButton(onPressed: _deleteAccount, child: Text('Delete Account'), style: TextButton.styleFrom(backgroundColor: Colors.red))),
        ],
      ),
    );
  }

  void _createNewAccount() async {
    setState(() => loading = true);

    final acc = Provider.of<AppAuth>(context, listen: false);

    try {
      final String? fcmToken;
      if (Platform.isLinux) {
        Toaster.warn("Unsupported Platform", 'Your platform is not supported by FCM - notifications will not work');
        fcmToken = '(linux-' + Uuid().v4() + ')';
      } else {
        final notificationSettings = await FirebaseMessaging.instance.requestPermission(provisional: true);

        if (notificationSettings.authorizationStatus == AuthorizationStatus.denied) {
          Toaster.error("Missing Permission", 'Please allow notifications to create an account');
          return;
        }

        fcmToken = await FirebaseMessaging.instance.getToken();
      }

      if (fcmToken == null) {
        Toaster.warn("Missing Token", 'No FCM Token found, please allow notifications, ensure you have a network connection and restart the app');
        return;
      }

      await Globals().setPrefFCMToken(fcmToken);

      final user = await APIClient.createUserWithClient(null, fcmToken, Globals().platform, Globals().version, Globals().hostname, Globals().clientType);

      acc.set(user.user, user.clients[0], user.adminKey, user.sendKey);

      await acc.save();
      Toaster.success("Success", 'Successfully Created a new account');
    } catch (exc, trace) {
      ApplicationLog.error('Failed to create user account: ' + exc.toString(), trace: trace);
      Toaster.error("Error", 'Failed to create user account');
    } finally {
      setState(() => loading = false);
    }
  }

  void _logout() async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    //TODO clear messages/channels/etc in open views
    acc.clear();
    await acc.save();

    Toaster.info('Logout', 'Successfully logged out');
  }

  void _deleteAccount() async {
    //TODO
  }
}
