import 'dart:io';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/pages/account/login.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:simplecloudnotifier/utils/ui.dart';
import 'package:uuid/uuid.dart';

class AccountRootPage extends StatefulWidget {
  const AccountRootPage({super.key, required this.isVisiblePage});

  final bool isVisiblePage;

  @override
  State<AccountRootPage> createState() => _AccountRootPageState();
}

class _AccountRootPageState extends State<AccountRootPage> {
  late ImmediateFuture<int>? futureSubscriptionCount;
  late ImmediateFuture<int>? futureClientCount;
  late ImmediateFuture<int>? futureKeyCount;
  late ImmediateFuture<int>? futureChannelAllCount;
  late ImmediateFuture<int>? futureChannelSubscribedCount;
  late ImmediateFuture<User>? futureUser;

  late AppAuth userAcc;

  bool loading = false;

  bool _isInitialized = false;

  @override
  void initState() {
    super.initState();

    userAcc = Provider.of<AppAuth>(context, listen: false);
    userAcc.addListener(_onAuthStateChanged);

    if (widget.isVisiblePage && !_isInitialized) _realInitState();
  }

  @override
  void didUpdateWidget(AccountRootPage oldWidget) {
    super.didUpdateWidget(oldWidget);

    if (oldWidget.isVisiblePage != widget.isVisiblePage && widget.isVisiblePage) {
      if (!_isInitialized) {
        _realInitState();
      } else {
        _backgroundRefresh();
      }
    }
  }

  void _realInitState() {
    ApplicationLog.debug('AccountRootPage::_realInitState');
    _onAuthStateChanged();
    _isInitialized = true;
  }

  @override
  void dispose() {
    ApplicationLog.debug('AccountRootPage::dispose');
    userAcc.removeListener(_onAuthStateChanged);
    super.dispose();
  }

  void _onAuthStateChanged() {
    ApplicationLog.debug('AccountRootPage::_onAuthStateChanged');
    _createFutures();
  }

  void _createFutures() {
    futureSubscriptionCount = null;
    futureClientCount = null;
    futureKeyCount = null;
    futureChannelAllCount = null;
    futureChannelSubscribedCount = null;

    if (userAcc.isAuth()) {
      futureChannelAllCount = ImmediateFuture.ofFuture(() async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc, ChannelSelector.all);
        return channels.length;
      }());

      futureChannelSubscribedCount = ImmediateFuture.ofFuture(() async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc, ChannelSelector.subscribed);
        return channels.length;
      }());

      futureSubscriptionCount = ImmediateFuture.ofFuture(() async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final subs = await APIClient.getSubscriptionList(userAcc);
        return subs.length;
      }());

      futureClientCount = ImmediateFuture.ofFuture(() async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final clients = await APIClient.getClientList(userAcc);
        return clients.length;
      }());

      futureKeyCount = ImmediateFuture.ofFuture(() async {
        if (!userAcc.isAuth()) throw new Exception('not logged in');
        final keys = await APIClient.getKeyTokenList(userAcc);
        return keys.length;
      }());

      futureUser = ImmediateFuture.ofFuture(userAcc.loadUser(force: false));
    }
  }

  Future<void> _backgroundRefresh() async {
    if (userAcc.isAuth()) {
      try {
        await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

        AppBarState().setLoadingIndeterminate(true);

        // refresh all data and then replace teh futures used in build()

        final channelsAll = await APIClient.getChannelList(userAcc, ChannelSelector.all);
        final channelsSubscribed = await APIClient.getChannelList(userAcc, ChannelSelector.subscribed);
        final subs = await APIClient.getSubscriptionList(userAcc);
        final clients = await APIClient.getClientList(userAcc);
        final keys = await APIClient.getKeyTokenList(userAcc);
        final user = await userAcc.loadUser(force: true);

        setState(() {
          futureChannelAllCount = ImmediateFuture.ofValue(channelsAll.length);
          futureChannelSubscribedCount = ImmediateFuture.ofValue(channelsSubscribed.length);
          futureSubscriptionCount = ImmediateFuture.ofValue(subs.length);
          futureClientCount = ImmediateFuture.ofValue(clients.length);
          futureKeyCount = ImmediateFuture.ofValue(keys.length);
          futureUser = ImmediateFuture.ofValue(user);
        });
      } catch (exc, trace) {
        ApplicationLog.error('Failed to refresh account data: ' + exc.toString(), trace: trace);
        Toaster.error("Error", 'Failed to refresh account data');
      } finally {
        AppBarState().setLoadingIndeterminate(false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<AppAuth>(
      builder: (context, acc, child) {
        if (!_isInitialized) return SizedBox();

        if (!userAcc.isAuth()) {
          return _buildNoAuth(context);
        } else {
          return FutureBuilder(
            future: futureUser!.future,
            builder: ((context, snapshot) {
              if (futureUser?.value != null) {
                return _buildShowAccount(context, acc, futureUser!.value!);
              } else if (snapshot.connectionState == ConnectionState.done && snapshot.hasError) {
                return Text('Error: ${snapshot.error}'); //TODO better error display
              } else if (snapshot.connectionState == ConnectionState.done) {
                return _buildShowAccount(context, acc, snapshot.data!);
              } else {
                return Center(child: CircularProgressIndicator());
              }
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
            UI.button(
              text: 'Create new account',
              onPressed: () {
                if (loading) return;
                _createNewAccount();
              },
              big: true,
            ),
            const SizedBox(height: 16),
            UI.button(
              text: 'Use existing account',
              onPressed: () {
                if (loading) return;
                Navi.push(context, () => AccountLoginPage());
              },
              tonal: true,
              big: true,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildShowAccount(BuildContext context, AppAuth acc, User user) {
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(8.0, 24.0, 8.0, 8.0),
        child: Column(
          children: [
            _buildHeader(context, user),
            const SizedBox(height: 16),
            Text(user.username ?? user.userID, overflow: TextOverflow.ellipsis, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
            const SizedBox(height: 16),
            ..._buildCards(context, user),
            SizedBox(height: 16),
            _buildFooter(context, user),
            SizedBox(height: 40),
          ],
        ),
      ),
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
                    future: futureChannelAllCount!.future,
                    builder: (context, snapshot) {
                      if (futureChannelAllCount?.value != null) {
                        return Text('${futureChannelAllCount!.value}');
                      } else if (snapshot.connectionState == ConnectionState.done) {
                        return Text('${snapshot.data}');
                      } else {
                        return const SizedBox(width: 8, height: 8, child: Center(child: CircularProgressIndicator()));
                      }
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
            UI.buttonIconOnly(
              onPressed: () {/*TODO*/},
              icon: FontAwesomeIcons.pen,
            ),
            const SizedBox(height: 4),
            if (!user.isPro)
              UI.buttonIconOnly(
                onPressed: () {/*TODO*/},
                icon: FontAwesomeIcons.cartCircleArrowUp,
              ),
          ],
        ),
      ],
    );
  }

  List<Widget> _buildCards(BuildContext context, User user) {
    return [
      _buildNumberCard(context, 'Subscriptions', futureSubscriptionCount, () {/*TODO*/}),
      _buildNumberCard(context, 'Clients', futureClientCount, () {/*TODO*/}),
      _buildNumberCard(context, 'Keys', futureKeyCount, () {/*TODO*/}),
      _buildNumberCard(context, 'Channels', futureChannelSubscribedCount, () {/*TODO*/}),
      UI.buttonCard(
        context: context,
        margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
        child: Row(
          children: [
            Text('${user.messagesSent}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
            const SizedBox(width: 12),
            Text('Messages', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
          ],
        ),
        onTap: () {/*TODO*/},
      ),
    ];
  }

  Widget _buildNumberCard(BuildContext context, String txt, ImmediateFuture<int>? future, void Function() action) {
    return UI.buttonCard(
      context: context,
      margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
      child: Row(
        children: [
          FutureBuilder(
            future: future?.future,
            builder: (context, snapshot) {
              if (future?.value != null) {
                return Text('${future?.value}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
              } else if (snapshot.connectionState == ConnectionState.done) {
                return Text('${snapshot.data}', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20));
              } else {
                return const SizedBox(width: 12, height: 12, child: Center(child: CircularProgressIndicator()));
              }
            },
          ),
          const SizedBox(width: 12),
          Text(txt, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
        ],
      ),
      onTap: action,
    );
  }

  Widget _buildFooter(BuildContext context, User user) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
      child: Row(
        children: [
          Expanded(
              child: UI.button(
            text: 'Logout',
            onPressed: _logout,
            color: Colors.orange,
          )),
          const SizedBox(width: 8),
          Expanded(
              child: UI.button(
            text: 'Delete Account',
            onPressed: _deleteAccount,
            color: Colors.red,
          )),
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
