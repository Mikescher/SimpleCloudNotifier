import 'package:flutter/material.dart';

class DebugColorsPage extends StatefulWidget {
  @override
  _DebugColorsPageState createState() => _DebugColorsPageState();
}

class _DebugColorsPageState extends State<DebugColorsPage> {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: SingleChildScrollView(
        child: Column(
          children: listColors(context),
        ),
      ),
    );
  }

  List<Widget> listColors(BuildContext context) {
    return [
      buildCol("primaryColor", Theme.of(context).primaryColor),
      buildCol("primaryColorDark", Theme.of(context).primaryColorDark),
      buildCol("primaryColorLight", Theme.of(context).primaryColorLight),
      buildCol("secondaryHeaderColor", Theme.of(context).secondaryHeaderColor),
      buildCol("disabledColor", Theme.of(context).disabledColor),
      buildCol("splashColor", Theme.of(context).splashColor),
      Divider(),
      buildCol("canvasColor", Theme.of(context).canvasColor),
      buildCol("cardColor", Theme.of(context).cardColor),
      buildCol("dialogBackgroundColor", Theme.of(context).dialogBackgroundColor),
      buildCol("dividerColor", Theme.of(context).dividerColor),
      buildCol("focusColor", Theme.of(context).focusColor),
      buildCol("highlightColor", Theme.of(context).highlightColor),
      buildCol("hintColor", Theme.of(context).hintColor),
      buildCol("hoverColor", Theme.of(context).hoverColor),
      buildCol("indicatorColor", Theme.of(context).indicatorColor),
      buildCol("scaffoldBackgroundColor", Theme.of(context).scaffoldBackgroundColor),
      buildCol("shadowColor", Theme.of(context).shadowColor),
      buildCol("unselectedWidgetColor", Theme.of(context).unselectedWidgetColor),
      Divider(),
      buildCol("colorScheme.primary", Theme.of(context).colorScheme.primary),
      buildCol("colorScheme.onPrimary", Theme.of(context).colorScheme.onPrimary),
      buildCol("colorScheme.primaryContainer", Theme.of(context).colorScheme.primaryContainer),
      buildCol("colorScheme.onPrimaryContainer", Theme.of(context).colorScheme.onPrimaryContainer),
      buildCol("colorScheme.inversePrimary", Theme.of(context).colorScheme.inversePrimary),
      buildCol("colorScheme.secondary", Theme.of(context).colorScheme.secondary),
      buildCol("colorScheme.onSecondary", Theme.of(context).colorScheme.onSecondary),
      buildCol("colorScheme.secondaryContainer", Theme.of(context).colorScheme.secondaryContainer),
      buildCol("colorScheme.onSecondaryContainer", Theme.of(context).colorScheme.onSecondaryContainer),
      buildCol("colorScheme.tertiary", Theme.of(context).colorScheme.tertiary),
      buildCol("colorScheme.onTertiary", Theme.of(context).colorScheme.onTertiary),
      buildCol("colorScheme.tertiaryContainer", Theme.of(context).colorScheme.tertiaryContainer),
      buildCol("colorScheme.onTertiaryContainer", Theme.of(context).colorScheme.onTertiaryContainer),
      buildCol("colorScheme.surface", Theme.of(context).colorScheme.surface),
      buildCol("colorScheme.onSurface", Theme.of(context).colorScheme.onSurface),
      buildCol("colorScheme.surfaceTint", Theme.of(context).colorScheme.surfaceTint),
      buildCol("colorScheme.surfaceVariant", Theme.of(context).colorScheme.surfaceVariant),
      buildCol("colorScheme.inverseSurface", Theme.of(context).colorScheme.inverseSurface),
      buildCol("colorScheme.onInverseSurface", Theme.of(context).colorScheme.onInverseSurface),
      buildCol("colorScheme.background", Theme.of(context).colorScheme.background),
      buildCol("colorScheme.onBackground", Theme.of(context).colorScheme.onBackground),
      buildCol("colorScheme.error", Theme.of(context).colorScheme.error),
      buildCol("colorScheme.onError", Theme.of(context).colorScheme.onError),
      buildCol("colorScheme.errorContainer", Theme.of(context).colorScheme.errorContainer),
      buildCol("colorScheme.onErrorContainer", Theme.of(context).colorScheme.onErrorContainer),
      buildCol("colorScheme.outline", Theme.of(context).colorScheme.outline),
      buildCol("colorScheme.outlineVariant", Theme.of(context).colorScheme.outlineVariant),
      buildCol("colorScheme.shadow", Theme.of(context).colorScheme.shadow),
      buildCol("colorScheme.scrim", Theme.of(context).colorScheme.scrim),
      Divider(),
      buildCol("primaryTextTheme.bodyLarge.backgroundColor", Theme.of(context).primaryTextTheme.bodyLarge?.backgroundColor),
      buildCol("primaryTextTheme.bodyLarge.color", Theme.of(context).primaryTextTheme.bodyLarge?.color),
      buildCol("primaryTextTheme.displayLarge.backgroundColor", Theme.of(context).primaryTextTheme.displayLarge?.backgroundColor),
      buildCol("primaryTextTheme.displayLarge.color", Theme.of(context).primaryTextTheme.displayLarge?.color),
      buildCol("primaryTextTheme.headlineLarge.backgroundColor", Theme.of(context).primaryTextTheme.headlineLarge?.backgroundColor),
      buildCol("primaryTextTheme.headlineLarge.color", Theme.of(context).primaryTextTheme.headlineLarge?.color),
      buildCol("primaryTextTheme.labelLarge.backgroundColor", Theme.of(context).primaryTextTheme.labelLarge?.backgroundColor),
      buildCol("primaryTextTheme.labelLarge.color", Theme.of(context).primaryTextTheme.labelLarge?.color),
      buildCol("primaryTextTheme.titleLarge.backgroundColor", Theme.of(context).primaryTextTheme.titleLarge?.backgroundColor),
      buildCol("primaryTextTheme.titleLarge.color", Theme.of(context).primaryTextTheme.titleLarge?.color),
      buildCol("textTheme.bodyLarge.backgroundColor", Theme.of(context).textTheme.bodyLarge?.backgroundColor),
      buildCol("textTheme.bodyLarge.color", Theme.of(context).textTheme.bodyLarge?.color),
      buildCol("textTheme.displayLarge.backgroundColor", Theme.of(context).textTheme.displayLarge?.backgroundColor),
      buildCol("textTheme.displayLarge.color", Theme.of(context).textTheme.displayLarge?.color),
      buildCol("textTheme.headlineLarge.backgroundColor", Theme.of(context).textTheme.headlineLarge?.backgroundColor),
      buildCol("textTheme.headlineLarge.color", Theme.of(context).textTheme.headlineLarge?.color),
      buildCol("textTheme.labelLarge.backgroundColor", Theme.of(context).textTheme.labelLarge?.backgroundColor),
      buildCol("textTheme.labelLarge.color", Theme.of(context).textTheme.labelLarge?.color),
      buildCol("textTheme.titleLarge.backgroundColor", Theme.of(context).textTheme.titleLarge?.backgroundColor),
      buildCol("textTheme.titleLarge.color", Theme.of(context).textTheme.titleLarge?.color),
      Divider(),
      buildCol("iconTheme.color", Theme.of(context).iconTheme.color),
      buildCol("primaryIconTheme.color", Theme.of(context).primaryIconTheme.color),
      buildCol("appBarTheme.foregroundColor", Theme.of(context).appBarTheme.foregroundColor),
      buildCol("appBarTheme.backgroundColor", Theme.of(context).appBarTheme.backgroundColor),
      buildCol("badgeTheme.textColor", Theme.of(context).badgeTheme.textColor),
      buildCol("badgeTheme.backgroundColor", Theme.of(context).badgeTheme.backgroundColor),
      buildCol("bannerTheme.backgroundColor", Theme.of(context).bannerTheme.backgroundColor),
      buildCol("bottomAppBarTheme.color", Theme.of(context).bottomAppBarTheme.color),
      buildCol("buttonTheme.colorScheme.background", Theme.of(context).buttonTheme.colorScheme?.background),
      buildCol("buttonTheme.colorScheme.primary", Theme.of(context).buttonTheme.colorScheme?.primary),
      buildCol("buttonTheme.colorScheme.secondary", Theme.of(context).buttonTheme.colorScheme?.secondary),
      buildCol("cardTheme.color", Theme.of(context).cardTheme.color),
    ];
  }

  Widget buildCol(String key, Color? value) {
    return Row(
      children: [
        Padding(
          padding: EdgeInsets.all(4),
          child: Container(
            width: 20,
            decoration: BoxDecoration(
              border: Border.all(color: Colors.black),
              color: value ?? Color.fromARGB(0, 0, 0, 0),
            ),
            height: 20,
          ),
        ),
        Expanded(child: Text(key))
      ],
    );
  }
}
