import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/key_token_auth.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

class AccountLoginPage extends StatefulWidget {
  final void Function()? onLogin;

  const AccountLoginPage({super.key, this.onLogin});

  @override
  State<AccountLoginPage> createState() => _AccountLoginPageState();
}

class _AccountLoginPageState extends State<AccountLoginPage> {
  late TextEditingController _ctrlUserID;
  late TextEditingController _ctrlToken;

  @override
  void initState() {
    super.initState();
    _ctrlUserID = TextEditingController();
    _ctrlToken = TextEditingController();
  }

  @override
  void dispose() {
    _ctrlUserID.dispose();
    _ctrlToken.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SizedBox(
            width: 250,
            child: TextField(
              controller: _ctrlUserID,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'UserID',
              ),
            ),
          ),
          const SizedBox(height: 16),
          SizedBox(
            width: 250,
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

      final verified = await APIClient.verifyToken(uid, tok);
      if (verified) {
        msgr.showSnackBar(
          const SnackBar(
            content: Text('Data ok'), //TODO toast
          ),
        );
        prov.setToken(KeyTokenAuth(userId: uid, token: tok));
        await prov.save();
        widget.onLogin?.call();
      } else {
        msgr.showSnackBar(
          const SnackBar(
            content: Text('Failed to verify token'), //TODO toast
          ),
        );
      }
    } catch (e) {
      //TODO
    }
  }
}
