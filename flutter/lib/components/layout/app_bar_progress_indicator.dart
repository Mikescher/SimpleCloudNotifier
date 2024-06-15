import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';

class AppBarProgressIndicator extends StatelessWidget implements PreferredSizeWidget {
  @override
  Size get preferredSize => Size(double.infinity, 1.0);

  @override
  Widget build(BuildContext context) {
    return Consumer<AppBarState>(
      builder: (context, value, child) {
        if (value.loadingIndeterminate) {
          return LinearProgressIndicator(value: null);
        } else {
          return SizedBox.square(dimension: 4); // 4 height is the same as the LinearProgressIndicator
        }
      },
    );
  }
}
