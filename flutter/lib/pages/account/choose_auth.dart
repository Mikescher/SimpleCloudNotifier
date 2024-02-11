import 'package:flutter/material.dart';

class AccountChoosePage extends StatelessWidget {
  final void Function()? onLogin;
  final void Function()? onCreateAccount;

  const AccountChoosePage({super.key, this.onLogin, this.onCreateAccount});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          ElevatedButton(
            style: ElevatedButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
            onPressed: () {
              onLogin?.call();
            },
            child: const Text('Use existing account'),
          ),
          const SizedBox(height: 32),
          ElevatedButton(
            style: ElevatedButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
            onPressed: () {
              onCreateAccount?.call();
            },
            child: const Text('Create new account'),
          ),
        ],
      ),
    );
  }
}
