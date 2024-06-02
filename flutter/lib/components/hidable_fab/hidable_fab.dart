import 'package:flutter/material.dart';

class HidableFAB extends StatelessWidget {
  final VoidCallback? onPressed;
  final IconData icon;

  const HidableFAB({
    super.key,
    this.onPressed,
    required this.icon,
  });

  Widget build(BuildContext context) {
    return Visibility(
      visible: MediaQuery.viewInsetsOf(context).bottom == 0.0, // hide when keyboard is shown
      child: FloatingActionButton(
        onPressed: onPressed,
        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(17))),
        elevation: 2.0,
        child: Icon(icon),
      ),
    );
  }
}
