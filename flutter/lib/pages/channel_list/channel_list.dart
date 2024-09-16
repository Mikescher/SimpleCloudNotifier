import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/pages/channel_view/channel_view.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list_item.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class ChannelRootPage extends StatefulWidget {
  const ChannelRootPage({super.key, required this.isVisiblePage});

  final bool isVisiblePage;

  @override
  State<ChannelRootPage> createState() => _ChannelRootPageState();
}

class _ChannelRootPageState extends State<ChannelRootPage> with RouteAware {
  final PagingController<int, ChannelWithSubscription> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: 0);

  bool _isInitialized = false;

  bool _reloadEnqueued = false;

  @override
  void initState() {
    super.initState();

    _pagingController.addPageRequestListener(_fetchPage);

    if (widget.isVisiblePage && !_isInitialized) _realInitState();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    Navi.modalRouteObserver.subscribe(this, ModalRoute.of(context)!);
  }

  @override
  void dispose() {
    ApplicationLog.debug('ChannelRootPage::dispose');
    _pagingController.dispose();
    Navi.modalRouteObserver.unsubscribe(this);
    super.dispose();
  }

  @override
  void didUpdateWidget(ChannelRootPage oldWidget) {
    super.didUpdateWidget(oldWidget);

    if (oldWidget.isVisiblePage != widget.isVisiblePage && widget.isVisiblePage) {
      if (!_isInitialized) {
        _realInitState();
      } else {
        _backgroundRefresh();
      }
    }
  }

  @override
  void didPush() {
    // ...
  }

  @override
  void didPopNext() {
    if (_reloadEnqueued) {
      ApplicationLog.debug('[ChannelList::RouteObserver] --> didPopNext (will background-refresh) (_reloadEnqueued == true)');
      () async {
        _reloadEnqueued = false;
        AppBarState().setLoadingIndeterminate(true);
        await Future.delayed(const Duration(milliseconds: 500)); // prevents flutter bug where the whole process crashes ?!?
        await _backgroundRefresh();
      }();
    }
  }

  void _realInitState() {
    ApplicationLog.debug('ChannelRootPage::_realInitState');
    _pagingController.refresh();
    _isInitialized = true;
  }

  Future<void> _fetchPage(int pageKey) async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    ApplicationLog.debug('Start ChannelList::_pagingController::_fetchPage [ ${pageKey} ]');

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      final items = (await APIClient.getChannelList(acc, ChannelSelector.all)).toList();

      items.sort((a, b) => -1 * (a.channel.timestampLastSent ?? '').compareTo(b.channel.timestampLastSent ?? ''));

      _pagingController.value = PagingState(nextPageKey: null, itemList: items, error: null);
    } catch (exc, trace) {
      _pagingController.error = exc.toString();
      ApplicationLog.error('Failed to list channels: ' + exc.toString(), trace: trace);
    }
  }

  Future<void> _backgroundRefresh() async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    ApplicationLog.debug('Start background refresh of channel list');

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      AppBarState().setLoadingIndeterminate(true);

      final items = (await APIClient.getChannelList(acc, ChannelSelector.all)).toList();

      items.sort((a, b) => -1 * (a.channel.timestampLastSent ?? '').compareTo(b.channel.timestampLastSent ?? ''));

      setState(() {
        _pagingController.value = PagingState(nextPageKey: null, itemList: items, error: null);
      });
    } catch (exc, trace) {
      setState(() {
        _pagingController.error = exc.toString();
      });
      ApplicationLog.error('Failed to list channels: ' + exc.toString(), trace: trace);
    } finally {
      AppBarState().setLoadingIndeterminate(false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () => Future.sync(
          () => _pagingController.refresh(),
        ),
        child: PagedListView<int, ChannelWithSubscription>(
          pagingController: _pagingController,
          builderDelegate: PagedChildBuilderDelegate<ChannelWithSubscription>(
            itemBuilder: (context, item, index) => ChannelListItem(
              channel: item.channel,
              subscription: item.subscription,
              onPressed: () {
                Navi.push(context, () => ChannelViewPage(channelID: item.channel.channelID, preloadedData: (item.channel, item.subscription), needsReload: _enqueueReload));
              },
            ),
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'fab_channel_list_qr',
        onPressed: () {
          //TODO scan qr code to subscribe channel
        },
        child: const Icon(FontAwesomeIcons.qrcode),
      ),
    );
  }

  void _enqueueReload() {
    _reloadEnqueued = true;
  }
}
