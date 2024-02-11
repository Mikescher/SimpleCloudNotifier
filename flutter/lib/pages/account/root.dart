import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/pages/account/login.dart';

import '../../state/user_account.dart';
import 'choose_auth.dart';

class AccountRootPage extends StatefulWidget {
  const AccountRootPage({super.key});

  @override
  State<AccountRootPage> createState() => _AccountRootPageState();
}

enum _SubPage { chooseAuth, login, main }

class _AccountRootPageState extends State<AccountRootPage> {
  late _SubPage _page;

  @override
  void initState() {
    super.initState();

    var prov = Provider.of<UserAccount>(context, listen: false);

    _page = (prov.auth != null) ? _SubPage.main : _SubPage.chooseAuth;

    prov.addListener(_onAuthStateChanged);
  }

  @override
  void dispose() {
    Provider.of<UserAccount>(context, listen: false).removeListener(_onAuthStateChanged);
    super.dispose();
  }

  void _onAuthStateChanged() {
    if (Provider.of<UserAccount>(context, listen: false).auth != null && _page != _SubPage.main) {
      setState(() {
        _page = _SubPage.main;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<UserAccount>(
      builder: (context, acc, child) {
        switch (_page) {
          case _SubPage.main:
            return const Center(
              child: Text(
                'Logged In',
                style: TextStyle(fontSize: 24),
              ),
            );
          case _SubPage.chooseAuth:
            return AccountChoosePage(
              onLogin: () => setState(() {
                _page = _SubPage.login;
              }),
              onCreateAccount: () => setState(() {
                //TODO
              }),
            );
          case _SubPage.login:
            return const AccountLoginPage();
        }
      },
    );
  }
}
