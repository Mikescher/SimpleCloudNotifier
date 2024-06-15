import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list_item.dart';

class ChannelRootPage extends StatefulWidget {
  const ChannelRootPage({super.key, required this.isVisiblePage});

  final bool isVisiblePage;

  @override
  State<ChannelRootPage> createState() => _ChannelRootPageState();
}

class _ChannelRootPageState extends State<ChannelRootPage> {
  final PagingController<int, Channel> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: 0);

  bool _isInitialized = false;

  @override
  void initState() {
    super.initState();

    _pagingController.addPageRequestListener(_fetchPage);

    if (widget.isVisiblePage && !_isInitialized) _realInitState();
  }

  @override
  void dispose() {
    ApplicationLog.debug('ChannelRootPage::dispose');
    _pagingController.dispose();
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
      final items = (await APIClient.getChannelList(acc, ChannelSelector.all)).map((p) => p.channel).toList();

      items.sort((a, b) => -1 * (a.timestampLastSent ?? '').compareTo(b.timestampLastSent ?? ''));

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

      final items = (await APIClient.getChannelList(acc, ChannelSelector.all)).map((p) => p.channel).toList();

      items.sort((a, b) => -1 * (a.timestampLastSent ?? '').compareTo(b.timestampLastSent ?? ''));

      _pagingController.value = PagingState(nextPageKey: null, itemList: items, error: null);
    } catch (exc, trace) {
      _pagingController.error = exc.toString();
      ApplicationLog.error('Failed to list channels: ' + exc.toString(), trace: trace);
    } finally {
      AppBarState().setLoadingIndeterminate(false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () => Future.sync(
        () => _pagingController.refresh(),
      ),
      child: PagedListView<int, Channel>(
        pagingController: _pagingController,
        builderDelegate: PagedChildBuilderDelegate<Channel>(
          itemBuilder: (context, item, index) => ChannelListItem(
            channel: item,
            onPressed: () {/*TODO*/},
          ),
        ),
      ),
    );
  }
}
