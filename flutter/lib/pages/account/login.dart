import 'package:flutter/material.dart';
import 'package:simplecloudnotifier/api/api_client.dart';

class AccountLoginPage extends StatefulWidget {
  const AccountLoginPage({super.key});

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

    final verified = await APIClient.verifyToken(_ctrlUserID.text, _ctrlToken.text);
    if (verified) {
      msgr.showSnackBar(
        const SnackBar(
          content: Text('Data ok'),
        ),
      );
    } else {
      msgr.showSnackBar(
        const SnackBar(
          content: Text('Failed to verify token'),
        ),
      );
    }
  }
}
