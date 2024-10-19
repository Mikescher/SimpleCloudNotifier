import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list_item.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class ChannelListExtendedPage extends StatefulWidget {
  const ChannelListExtendedPage({super.key});

  @override
  State<ChannelListExtendedPage> createState() => _ChannelListExtendedPageState();
}

class _ChannelListExtendedPageState extends State<ChannelListExtendedPage> with RouteAware {
  final PagingController<int, ChannelWithSubscription> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: 0);

  bool _reloadEnqueued = false;

  @override
  void initState() {
    super.initState();

    _pagingController.addPageRequestListener(_fetchPage);

    _pagingController.refresh();
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
  void didPopNext() {
    if (_reloadEnqueued) {
      ApplicationLog.debug('[ChannelList::RouteObserver] --> didPopNext (will background-refresh) (_reloadEnqueued == true)');
      () async {
        _reloadEnqueued = false;
        await Future.delayed(const Duration(milliseconds: 500), () {}); // prevents flutter bug where the whole process crashes ?!?
        await _backgroundRefresh();
      }();
    }
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
    return SCNScaffold(
      title: "Channels",
      showSearch: false,
      showShare: false,
      child: Padding(
        padding: EdgeInsets.fromLTRB(8, 4, 8, 4),
        child: RefreshIndicator(
          onRefresh: () => Future.sync(
            () => _pagingController.refresh(),
          ),
          child: PagedListView<int, ChannelWithSubscription>(
            pagingController: _pagingController,
            builderDelegate: PagedChildBuilderDelegate<ChannelWithSubscription>(
              itemBuilder: (context, item, index) => ChannelListItem(
                channel: item.channel,
                subscription: item.subscription,
                mode: ChannelListItemMode.Extended,
                onChannelListReloadTrigger: _enqueueReload,
                onSubscriptionChanged: (channelID, subscription) {
                  setState(() {
                    final idx = _pagingController.itemList?.indexWhere((p) => p.channel.channelID == channelID);
                    if (idx != null && idx >= 0) _pagingController.itemList![idx] = ChannelWithSubscription(channel: _pagingController.itemList![idx].channel, subscription: subscription);
                  });
                },
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _enqueueReload() {
    _reloadEnqueued = true;
  }
}
