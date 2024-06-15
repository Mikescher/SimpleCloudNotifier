import 'package:flutter/material.dart';

class AppBarProgressIndicator extends StatelessWidget implements PreferredSizeWidget {
  AppBarProgressIndicator({required this.show});

  final bool show;

  @override
  Size get preferredSize => Size(double.infinity, 1.0);

  @override
  Widget build(BuildContext context) {
    if (show) {
      return LinearProgressIndicator(value: null);
    } else {
      return SizedBox.square(dimension: 4); // 4 height is the same as the LinearProgressIndicator
    }
  }
}
