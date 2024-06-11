import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

class UI {
  static const double DefaultBorderRadius = 4;

  static Widget button({required String text, required void Function() onPressed, bool big = false, Color? color = null, bool tonal = false, IconData? icon = null}) {
    final double fontSize = big ? 24 : 14;
    final padding = big ? EdgeInsets.fromLTRB(8, 12, 8, 12) : null;

    final style = FilledButton.styleFrom(
      textStyle: TextStyle(fontSize: fontSize),
      padding: padding,
      backgroundColor: color,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(DefaultBorderRadius)),
    );

    if (tonal) {
      if (icon != null) {
        return FilledButton.tonalIcon(
          style: style,
          onPressed: onPressed,
          icon: Icon(icon),
          label: Text(text),
        );
      } else {
        return FilledButton.tonal(
          style: style,
          onPressed: onPressed,
          child: Text(text),
        );
      }
    } else {
      if (icon != null) {
        return FilledButton.icon(
          style: style,
          onPressed: onPressed,
          icon: Icon(icon),
          label: Text(text),
        );
      } else {
        return FilledButton(
          style: style,
          onPressed: onPressed,
          child: Text(text),
        );
      }
    }
  }

  static Widget buttonIconOnly({
    required void Function() onPressed,
    required IconData icon,
    double? iconSize = null,
  }) {
    return IconButton(
      icon: FaIcon(icon),
      iconSize: iconSize ?? 18,
      padding: EdgeInsets.all(4),
      constraints: BoxConstraints(),
      style: ButtonStyle(tapTargetSize: MaterialTapTargetSize.shrinkWrap),
      onPressed: onPressed,
    );
  }

  static Widget buttonCard({required BuildContext context, required Widget child, required void Function() onTap, EdgeInsets? margin = null}) {
    return Card.filled(
      margin: margin,
      shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(DefaultBorderRadius)),
      color: Theme.of(context).cardTheme.color,
      child: InkWell(
        splashColor: Theme.of(context).splashColor,
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: child,
        ),
      ),
    );
  }

  static Widget channelChip({required BuildContext context, required String text, EdgeInsets? margin = null, double fontSize = 12}) {
    return Container(
      padding: const EdgeInsets.fromLTRB(4, 0, 4, 0),
      margin: margin,
      decoration: BoxDecoration(
        color: Theme.of(context).hintColor,
        borderRadius: BorderRadius.all(Radius.circular(DefaultBorderRadius)),
      ),
      child: Text(
        text,
        style: TextStyle(fontWeight: FontWeight.bold, color: Theme.of(context).cardColor, fontSize: fontSize),
        overflow: TextOverflow.clip,
        maxLines: 1,
      ),
    );
  }

  static Widget box({required BuildContext context, required Widget child, required EdgeInsets? padding, Color? borderColor = null}) {
    return Container(
      padding: padding ?? EdgeInsets.all(4),
      decoration: BoxDecoration(
        border: Border.all(color: borderColor ?? Theme.of(context).hintColor),
        borderRadius: BorderRadius.circular(DefaultBorderRadius),
      ),
      child: child,
    );
  }
}
