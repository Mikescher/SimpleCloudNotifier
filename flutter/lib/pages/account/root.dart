import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/pages/account/login.dart';
import 'package:simplecloudnotifier/state/user_account.dart';
import 'package:simplecloudnotifier/pages/account/choose_auth.dart';

class AccountRootPage extends StatefulWidget {
  const AccountRootPage({super.key});

  @override
  State<AccountRootPage> createState() => _AccountRootPageState();
}

enum _SubPage { chooseAuth, login, main }

class _AccountRootPageState extends State<AccountRootPage> {
  late _SubPage _page;
  late UserAccount userAcc;

  @override
  void initState() {
    super.initState();

    userAcc = Provider.of<UserAccount>(context, listen: false);

    _page = (userAcc.auth != null) ? _SubPage.main : _SubPage.chooseAuth;

    userAcc.addListener(_onAuthStateChanged);
  }

  @override
  void dispose() {
    userAcc.removeListener(_onAuthStateChanged);
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
            return Center(
              child: Text(
                'Logged In: ${acc.auth?.userId}',
                style: const TextStyle(fontSize: 24),
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
            return AccountLoginPage(
              onLogin: () => setState(() {
                _page = _SubPage.main;
              }),
            );
        }
      },
    );
  }
}
