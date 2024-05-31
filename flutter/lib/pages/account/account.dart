import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

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

  late UserAccount userAcc;

  @override
  void initState() {
    super.initState();

    userAcc = Provider.of<UserAccount>(context, listen: false);
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

    if (userAcc.auth != null) {
      futureChannelAllCount = () async {
        if (userAcc.auth == null) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc.auth!, ChannelSelector.all);
        return channels.length;
      }();

      futureChannelSubscribedCount = () async {
        if (userAcc.auth == null) throw new Exception('not logged in');
        final channels = await APIClient.getChannelList(userAcc.auth!, ChannelSelector.subscribed);
        return channels.length;
      }();

      futureSubscriptionCount = () async {
        if (userAcc.auth == null) throw new Exception('not logged in');
        final subs = await APIClient.getSubscriptionList(userAcc.auth!);
        return subs.length;
      }();

      futureClientCount = () async {
        if (userAcc.auth == null) throw new Exception('not logged in');
        final clients = await APIClient.getClientList(userAcc.auth!);
        return clients.length;
      }();

      futureKeyCount = () async {
        if (userAcc.auth == null) throw new Exception('not logged in');
        final keys = await APIClient.getKeyTokenList(userAcc.auth!);
        return keys.length;
      }();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<UserAccount>(
      builder: (context, acc, child) {
        if (acc.auth == null) {
          return buildNoAuth(context);
        } else {
          return FutureBuilder(
            future: acc.loadUser(false),
            builder: ((context, snapshot) {
              if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasError) {
                  return Text('Error: ${snapshot.error}'); //TODO better error display
                }
                return buildShowAccount(context, acc, snapshot.data!);
              }
              return Center(child: CircularProgressIndicator());
            }),
          );
        }
      },
    );
  }

  Widget buildNoAuth(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          ElevatedButton(
            style: ElevatedButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
            onPressed: () {
              //TODO
            },
            child: const Text('Use existing account'),
          ),
          const SizedBox(height: 32),
          ElevatedButton(
            style: ElevatedButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
            onPressed: () {
              //TODO
            },
            child: const Text('Create new account'),
          ),
        ],
      ),
    );
  }

  Widget buildShowAccount(BuildContext context, UserAccount acc, User user) {
    //TODO better layout
    return Column(
      children: [
        SingleChildScrollView(
          scrollDirection: Axis.vertical,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 24.0, 8.0, 8.0),
            child: Column(
              children: [
                buildHeader(context, user),
                const SizedBox(height: 16),
                Text(user.username ?? user.userID, overflow: TextOverflow.ellipsis, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
                const SizedBox(height: 16),
                ...buildCards(context, user),
              ],
            ),
          ),
        ),
        const Expanded(child: SizedBox(height: 16)),
        buildFooter(context, user),
        SizedBox(height: 40)
      ],
    );
  }

  Row buildHeader(BuildContext context, User user) {
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
                      borderRadius: BorderRadius.only(topLeft: Radius.circular(8)),
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
                      borderRadius: BorderRadius.only(topLeft: Radius.circular(8)),
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
              Row(
                children: [
                  Expanded(child: Text(user.username ?? user.userID, overflow: TextOverflow.ellipsis)),
                  IconButton(
                    icon: FaIcon(FontAwesomeIcons.pen),
                    iconSize: 18,
                    padding: EdgeInsets.fromLTRB(0, 0, 4, 0),
                    constraints: BoxConstraints(),
                    onPressed: () {/*TODO*/},
                  ),
                ],
              ),
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
      ],
    );
  }

  List<Widget> buildCards(BuildContext context, User user) {
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

  Widget buildFooter(BuildContext context, User user) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(8, 0, 8, 0),
      child: Row(
        children: [
          Expanded(child: FilledButton(onPressed: () {/*TODO*/}, child: Text('Logout'), style: TextButton.styleFrom(backgroundColor: Colors.orange))),
          const SizedBox(width: 8),
          Expanded(child: FilledButton(onPressed: () {/*TODO*/}, child: Text('Delete Account'), style: TextButton.styleFrom(backgroundColor: Colors.red))),
        ],
      ),
    );
  }
}
