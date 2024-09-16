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

  static Widget metaCard({required BuildContext context, required IconData icon, required String title, required List<String> values, void Function()? mainAction, List<(IconData, void Function())>? iconActions}) {
    final container = UI.box(
      context: context,
      padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
      child: Row(
        children: [
          FaIcon(icon, size: 18),
          SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                for (final val in values) Text(val, style: const TextStyle(fontSize: 14)),
              ],
            ),
          ),
          if (iconActions != null) ...[
            SizedBox(width: 12),
            for (final iconAction in iconActions) ...[
              SizedBox(width: 4),
              IconButton(icon: FaIcon(iconAction.$1), onPressed: iconAction.$2),
            ],
          ],
        ],
      ),
    );

    if (mainAction == null) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: container,
      );
    } else {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: InkWell(
          splashColor: Theme.of(context).splashColor,
          onTap: mainAction,
          child: container,
        ),
      );
    }
  }
}
