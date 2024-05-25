import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

class AccountLoginPage extends StatefulWidget {
  final void Function()? onLogin;

  const AccountLoginPage({super.key, this.onLogin});

  @override
  State<AccountLoginPage> createState() => _AccountLoginPageState();
}

class _AccountLoginPageState extends State<AccountLoginPage> {
  final TextEditingController _ctrlUserID = TextEditingController();
  final TextEditingController _ctrlToken = TextEditingController();

  @override
  void dispose() {
    _ctrlUserID.dispose();
    _ctrlToken.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
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
              controller: _ctrlToken,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Token',
              ),
            ),
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            style: ElevatedButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
            onPressed: _login,
            child: const Text('Login'),
          ),
        ],
      ),
    );
  }

  void _login() async {
    final msgr = ScaffoldMessenger.of(context);
    final prov = Provider.of<UserAccount>(context, listen: false);

    try {
      final uid = _ctrlUserID.text;
      final tok = _ctrlToken.text;

      final verified = await APIClient.verifyToken(uid, tok); //TODO verify that this is an perm=ADMIN key
      if (verified) {
        msgr.showSnackBar(
          const SnackBar(
            content: Text('Data ok'), //TODO use toast?
          ),
        );
        prov.setToken(KeyTokenAuth(userId: uid, token: tok));
        await prov.save();
        widget.onLogin?.call();
      } else {
        msgr.showSnackBar(
          const SnackBar(
            content: Text('Failed to verify token'), //TODO use toast?
          ),
        );
      }
    } catch (exc, trace) {
      ApplicationLog.error('Failed to verify token: ' + exc.toString(), trace: trace);
    }
  }
}
