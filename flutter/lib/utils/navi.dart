import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';

class Navi {
  static final SCNRouteObserver routeObserver = SCNRouteObserver();
  static final RouteObserver<ModalRoute<void>> modalRouteObserver = RouteObserver<ModalRoute<void>>();

  static void push<T extends Widget>(BuildContext context, T Function() builder) {
    Provider.of<AppBarState>(context, listen: false).setLoadingIndeterminate(false);
    Provider.of<AppBarState>(context, listen: false).setShowSearchField(false);

    Navigator.push(context, MaterialPageRoute<T>(builder: (context) => builder()));
  }

  static void popToRoot(BuildContext context) {
    Provider.of<AppBarState>(context, listen: false).setLoadingIndeterminate(false);
    Provider.of<AppBarState>(context, listen: false).setShowSearchField(false);

    Navigator.popUntil(context, (route) => route.isFirst);
  }

  static void popDialog(BuildContext dialogContext) {
    Navigator.pop(dialogContext);
  }
}

class SCNRouteObserver extends RouteObserver<PageRoute<dynamic>> {
  @override
  void didPush(Route<dynamic> route, Route<dynamic>? previousRoute) {
    super.didPush(route, previousRoute);
    if (route is PageRoute) {
      AppBarState().setLoadingIndeterminate(false);
      AppBarState().setShowSearchField(false);

      print('[SCNRouteObserver] .didPush()');
    }
  }

  @override
  void didReplace({Route<dynamic>? newRoute, Route<dynamic>? oldRoute}) {
    super.didReplace(newRoute: newRoute, oldRoute: oldRoute);
    if (newRoute is PageRoute) {
      AppBarState().setLoadingIndeterminate(false);
      AppBarState().setShowSearchField(false);

      print('[SCNRouteObserver] .didReplace()');
    }
  }

  @override
  void didPop(Route<dynamic> route, Route<dynamic>? previousRoute) {
    super.didPop(route, previousRoute);
    if (previousRoute is PageRoute && route is PageRoute) {
      AppBarState().setLoadingIndeterminate(false);
      AppBarState().setShowSearchField(false);

      print('[SCNRouteObserver] .didPop()');
    }
  }
}
