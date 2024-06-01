import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/globals.dart';
import 'package:simplecloudnotifier/state/user_account.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';

class AccountLoginPage extends StatefulWidget {
  const AccountLoginPage({super.key});

  @override
  State<AccountLoginPage> createState() => _AccountLoginPageState();
}

class _AccountLoginPageState extends State<AccountLoginPage> {
  final TextEditingController _ctrlUserID = TextEditingController();
  final TextEditingController _ctrlTokenAdmin = TextEditingController();
  final TextEditingController _ctrlTokenSend = TextEditingController();

  bool loading = false;

  @override
  void dispose() {
    _ctrlUserID.dispose();
    _ctrlTokenAdmin.dispose();
    _ctrlTokenSend.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: 'Login',
      showSearch: false,
      child: SingleChildScrollView(
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
                    child: Center(child: FaIcon(FontAwesomeIcons.solidRightToBracket, size: 96, color: Theme.of(context).colorScheme.onSecondary)),
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
              FractionallySizedBox(
                widthFactor: 1.0,
                child: TextField(
                  controller: _ctrlUserID,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'UserID',
                  ),
                ),
              ),
              const SizedBox(height: 16),
              FractionallySizedBox(
                widthFactor: 1.0,
                child: TextField(
                  controller: _ctrlTokenAdmin,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Admin Token',
                  ),
                ),
              ),
              const SizedBox(height: 16),
              FractionallySizedBox(
                widthFactor: 1.0,
                child: TextField(
                  controller: _ctrlTokenSend,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Send Token (optional)',
                  ),
                ),
              ),
              const SizedBox(height: 16),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 24), padding: const EdgeInsets.fromLTRB(8, 12, 8, 12)),
                onPressed: _login,
                child: const Text('Login'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _login() async {
    final acc = Provider.of<UserAccount>(context, listen: false);

    try {
      setState(() => loading = true);

      final uid = _ctrlUserID.text;
      final atokv = _ctrlTokenAdmin.text;
      final stokv = _ctrlTokenSend.text;

      final fcmToken = await FirebaseMessaging.instance.getToken();

      if (fcmToken == null) {
        Toaster.warn("Missing Token", 'No FCM Token found, please allow notifications, ensure you have a network connection and restart the app');
        return;
      }

      final toka = await APIClient.getKeyTokenByToken(uid, atokv);

      if (!toka.allChannels || toka.permissions != 'A') {
        Toaster.error("Error", 'Admin token does not have required permissions');
        return;
      }

      final toks = await APIClient.getKeyTokenByToken(uid, stokv);

      if (!toks.allChannels || toks.permissions != 'CS') {
        Toaster.error("Error", 'Send token does not have required permissions');
        return;
      }

      final kta = KeyTokenAuth(userId: uid, tokenAdmin: atokv, tokenSend: stokv);

      final user = await APIClient.getUser(kta, uid);

      final client = await APIClient.addClient(acc.auth, fcmToken, Globals().deviceModel, Globals().version, Globals().hostname, Globals().clientType);

      acc.set(user, client, kta);
      await acc.save();

      Toaster.success("Login", "Successfully logged in");
    } catch (exc, trace) {
      ApplicationLog.error('Failed to verify token: ' + exc.toString(), trace: trace);
      Toaster.error("Error", 'Failed to verify token');
    } finally {
      setState(() => loading = false);
    }
  }
}
