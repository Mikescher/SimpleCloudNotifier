import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/pages/message_view/message_view.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list_item.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class MessageListPage extends StatefulWidget {
  const MessageListPage({super.key, required this.isVisiblePage});

  final bool isVisiblePage;

  //TODO reload on switch to tab
  //TODO reload on app to foreground

  @override
  State<MessageListPage> createState() => _MessageListPageState();
}

class _MessageListPageState extends State<MessageListPage> with RouteAware {
  static const _pageSize = 128;

  late final AppLifecycleListener _lifecyleListener;

  PagingController<String, Message> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: '@start');

  Map<String, Channel>? _channels = null;

  bool _isInitialized = false;

  @override
  void initState() {
    super.initState();

    _pagingController.addPageRequestListener(_fetchPage);

    if (widget.isVisiblePage && !_isInitialized) _realInitState();

    _lifecyleListener = AppLifecycleListener(
      onResume: _onLifecycleResume,
    );
  }

  @override
  void didUpdateWidget(MessageListPage oldWidget) {
    super.didUpdateWidget(oldWidget);

    if (oldWidget.isVisiblePage != widget.isVisiblePage && widget.isVisiblePage) {
      if (!_isInitialized) {
        _realInitState();
      } else {
        _backgroundRefresh(false);
      }
    }
  }

  void _realInitState() {
    ApplicationLog.debug('MessageListPage::_realInitState');

    final chnCache = Hive.box<Channel>('scn-channel-cache');
    final msgCache = Hive.box<Message>('scn-message-cache');

    if (chnCache.isNotEmpty && msgCache.isNotEmpty) {
      // ==== Use cache values - and refresh in background

      _channels = <String, Channel>{for (var v in chnCache.values) v.channelID: v};

      final cacheMessages = msgCache.values.toList();
      cacheMessages.sort((a, b) => -1 * a.timestamp.compareTo(b.timestamp));

      _pagingController.value = PagingState(nextPageKey: null, itemList: cacheMessages, error: null);

      _backgroundRefresh(true);
    } else {
      // ==== Full refresh - no cache available
      _pagingController.refresh();
    }

    _isInitialized = true;
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    Navi.modalRouteObserver.subscribe(this, ModalRoute.of(context)!);
  }

  @override
  void dispose() {
    ApplicationLog.debug('MessageListPage::dispose');
    Navi.modalRouteObserver.unsubscribe(this);
    _pagingController.dispose();
    _lifecyleListener.dispose();
    super.dispose();
  }

  @override
  void didPush() {
    // ...
  }

  @override
  void didPopNext() {
    ApplicationLog.debug('[MessageList::RouteObserver] --> didPopNext (will background-refresh)');
    _backgroundRefresh(false);
  }

  void _onLifecycleResume() {
    ApplicationLog.debug('[MessageList::_onLifecycleResume] --> (will background-refresh)');
    _backgroundRefresh(false);
  }

  Future<void> _fetchPage(String thisPageToken) async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    ApplicationLog.debug('Start MessageList::_pagingController::_fetchPage [ ${thisPageToken} ]');

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      if (_channels == null) {
        final channels = await APIClient.getChannelList(acc, ChannelSelector.allAny);
        _channels = <String, Channel>{for (var v in channels) v.channel.channelID: v.channel};

        _setChannelCache(channels); // no await
      }

      final (npt, newItems) = await APIClient.getMessageList(acc, thisPageToken, pageSize: _pageSize);

      _addToMessageCache(newItems); // no await

      ApplicationLog.debug('Finished MessageList::_pagingController::_fetchPage [ ${newItems.length} items and npt: ${thisPageToken} --> ${npt} ]');

      if (npt == '@end') {
        _pagingController.appendLastPage(newItems);
      } else {
        _pagingController.appendPage(newItems, npt);
      }
    } catch (exc, trace) {
      _pagingController.error = exc.toString();
      ApplicationLog.error('Failed to list messages: ' + exc.toString(), trace: trace);
    }
  }

  Future<void> _backgroundRefresh(bool fullReplaceState) async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    ApplicationLog.debug('Start background refresh of message list (fullReplaceState: $fullReplaceState)');

    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      AppBarState().setLoadingIndeterminate(true);

      if (_channels == null || fullReplaceState) {
        final channels = await APIClient.getChannelList(acc, ChannelSelector.allAny);
        setState(() {
          _channels = <String, Channel>{for (var v in channels) v.channel.channelID: v.channel};
        });
        _setChannelCache(channels); // no await
      }

      final (npt, newItems) = await APIClient.getMessageList(acc, '@start', pageSize: _pageSize);

      _addToMessageCache(newItems); // no await

      if (fullReplaceState) {
        // fully replace/reset state
        ApplicationLog.debug('Background-refresh finished (fullReplaceState) - replace state with ${newItems.length} items and npt: [ $npt ]');
        setState(() {
          if (npt == '@end')
            _pagingController.value = PagingState(nextPageKey: null, itemList: newItems, error: null);
          else
            _pagingController.value = PagingState(nextPageKey: npt, itemList: newItems, error: null);
        });
      } else {
        final itemsToBeAdded = newItems.where((p1) => !(_pagingController.itemList ?? []).any((p2) => p1.messageID == p2.messageID)).toList();
        if (itemsToBeAdded.isEmpty) {
          // nothing to do - no new items...
          // ....
          ApplicationLog.debug('Background-refresh returned no new items - nothing to do.');
        } else if (itemsToBeAdded.length == newItems.length) {
          // all items are new ?!?, the current state is completely fucked - full replace
          ApplicationLog.debug('Background-refresh found only new items ?!? - fully replace state with ${newItems.length} items');
          setState(() {
            if (npt == '@end')
              _pagingController.value = PagingState(nextPageKey: null, itemList: newItems, error: null);
            else
              _pagingController.value = PagingState(nextPageKey: npt, itemList: newItems, error: null);
            _pagingController.itemList = null;
          });
        } else {
          // add new items to the front
          ApplicationLog.debug('Background-refresh found ${newItems.length} new items - add to front');
          setState(() {
            _pagingController.itemList = itemsToBeAdded + (_pagingController.itemList ?? []);
          });
        }
      }
    } catch (exc, trace) {
      setState(() {
        _pagingController.error = exc.toString();
      });
      ApplicationLog.error('Failed to list messages: ' + exc.toString(), trace: trace);
    } finally {
      AppBarState().setLoadingIndeterminate(false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.fromLTRB(8, 4, 8, 4),
      child: RefreshIndicator(
        onRefresh: () => Future.sync(
          () => _pagingController.refresh(),
        ),
        child: PagedListView<String, Message>(
          pagingController: _pagingController,
          builderDelegate: PagedChildBuilderDelegate<Message>(
            itemBuilder: (context, item, index) => MessageListItem(
              message: item,
              allChannels: _channels ?? {},
              onPressed: () {
                Navi.push(context, () => MessageViewPage(message: item));
              },
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _setChannelCache(List<ChannelWithSubscription> channels) async {
    final cache = Hive.box<Channel>('scn-channel-cache');

    if (cache.length != channels.length) await cache.clear();

    for (var chn in channels) await cache.put(chn.channel.channelID, chn.channel);
  }

  Future<void> _addToMessageCache(List<Message> newItems) async {
    final cache = Hive.box<Message>('scn-message-cache');

    for (var msg in newItems) await cache.put(msg.messageID, msg);

    // delete all but the newest 128 messages

    if (cache.length < _pageSize) return;

    final allValues = cache.values.toList();

    allValues.sort((a, b) => -1 * a.timestamp.compareTo(b.timestamp));

    for (var val in allValues.sublist(_pageSize)) {
      await cache.delete(val.messageID);
    }
  }
}
