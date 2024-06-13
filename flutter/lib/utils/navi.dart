import 'package:flutter/material.dart';

class Navi {
  static void push<T extends Widget>(BuildContext context, T Function() builder) {
    Navigator.push(context, MaterialPageRoute<T>(builder: (context) => builder()));
  }

  static void popToRoot(BuildContext context) {
    Navigator.popUntil(context, (route) => route.isFirst);
  }
}
