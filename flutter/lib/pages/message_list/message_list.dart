import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/pages/message_list/message_filter_chiplet.dart';
import 'package:simplecloudnotifier/pages/message_view/message_view.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/app_events.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list_item.dart';
import 'package:simplecloudnotifier/state/scn_data_cache.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class MessageListPage extends StatefulWidget {
  const MessageListPage({super.key, required this.isVisiblePage});

  final bool isVisiblePage;

  @override
  State<MessageListPage> createState() => _MessageListPageState();
}

class _MessageListPageState extends State<MessageListPage> with RouteAware {
  late final AppLifecycleListener _lifecyleListener;

  PagingController<String, SCNMessage> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: '@start');

  Map<String, Channel>? _channels = null;

  bool _isInitialized = false;

  List<MessageFilterChiplet> _filterChiplets = [];

  @override
  void initState() {
    super.initState();

    AppEvents().subscribeFilterListener(_onAddFilter);
    AppEvents().subscribeMessageReceivedListener(_onMessageReceivedViaNotification);

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

    if (SCNDataCache().hasMessagesAndChannels()) {
      // ==== Use cache values - and refresh in background

      _channels = SCNDataCache().getChannelMap();

      _pagingController.value = PagingState(nextPageKey: null, itemList: SCNDataCache().getMessagesSorted(), error: null);

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
    AppEvents().unsubscribeFilterListener(_onAddFilter);
    AppEvents().unsubscribeMessageReceivedListener(_onMessageReceivedViaNotification);
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
    if (AppSettings().backgroundRefreshMessageListOnPop) {
      ApplicationLog.debug('[MessageList::RouteObserver] --> didPopNext (will background-refresh)');
      _backgroundRefresh(false);
    }
  }

  void _onLifecycleResume() {
    if (AppSettings().alwaysBackgroundRefreshMessageListOnLifecycleResume && widget.isVisiblePage) {
      ApplicationLog.debug('[MessageList::_onLifecycleResume] --> (will background-refresh)');
      _backgroundRefresh(false);
    }
  }

  Future<void> _fetchPage(String thisPageToken) async {
    final acc = Provider.of<AppAuth>(context, listen: false);
    final cfg = Provider.of<AppSettings>(context, listen: false);

    ApplicationLog.debug('Start MessageList::_pagingController::_fetchPage [ ${thisPageToken} ]');

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      if (_channels == null) {
        final channels = await APIClient.getChannelList(acc, ChannelSelector.allAny);
        _channels = <String, Channel>{for (var v in channels) v.channel.channelID: v.channel};

        SCNDataCache().setChannelCache(channels); // no await
      }

      final (npt, newItems) = await APIClient.getMessageList(acc, thisPageToken, pageSize: cfg.messagePageSize, filter: _getFilter());

      SCNDataCache().addToMessageCache(newItems); // no await

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
    final cfg = Provider.of<AppSettings>(context, listen: false);

    ApplicationLog.debug('Start background refresh of message list (fullReplaceState: $fullReplaceState)');

    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      AppBarState().setLoadingIndeterminate(true);

      if (_channels == null || fullReplaceState) {
        final channels = await APIClient.getChannelList(acc, ChannelSelector.allAny);
        setState(() {
          _channels = <String, Channel>{for (var v in channels) v.channel.channelID: v.channel};
        });
        SCNDataCache().setChannelCache(channels); // no await
      }

      final (npt, newItems) = await APIClient.getMessageList(acc, '@start', pageSize: cfg.messagePageSize);

      SCNDataCache().addToMessageCache(newItems); // no await

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
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          if (_filterChiplets.isNotEmpty)
            Wrap(
              alignment: WrapAlignment.start,
              spacing: 5.0,
              children: [
                for (var chiplet in _filterChiplets) _buildFilterChip(context, chiplet),
              ],
            ),
          Expanded(
            child: RefreshIndicator(
              onRefresh: () => Future.sync(
                () => _pagingController.refresh(),
              ),
              child: PagedListView<String, SCNMessage>(
                pagingController: _pagingController,
                builderDelegate: PagedChildBuilderDelegate<SCNMessage>(
                  itemBuilder: (context, item, index) => MessageListItem(
                    message: item,
                    allChannels: _channels ?? {},
                    onPressed: () {
                      Navi.push(context, () => MessageViewPage(messageID: item.messageID, preloadedData: (item,)));
                    },
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilterChip(BuildContext context, MessageFilterChiplet chiplet) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 2, 0, 2),
      child: InputChip(
        avatar: Icon(chiplet.icon()),
        label: Text(chiplet.label),
        onDeleted: () => _onRemFilter(chiplet),
        onPressed: () {/* TODO idk what to do here ? */},
        visualDensity: VisualDensity(horizontal: -4, vertical: -4),
      ),
    );
  }

  void _onAddFilter(List<MessageFilterChipletType> remTypeList, List<MessageFilterChiplet> chiplets) {
    setState(() {
      final remTypes = remTypeList.toSet();

      _filterChiplets = _filterChiplets.where((element) => !remTypes.contains(element.type)).toList() + chiplets;

      _pagingController.refresh();
    });
  }

  void _onRemFilter(MessageFilterChiplet chiplet) {
    setState(() {
      _filterChiplets.remove(chiplet);

      _pagingController.refresh();
    });
  }

  void _onMessageReceivedViaNotification(SCNMessage msg) {
    setState(() {
      _pagingController.itemList = [msg] + (_pagingController.itemList ?? []);
    });
  }

  MessageFilter _getFilter() {
    var filter = MessageFilter();

    var chipletsChannel = _filterChiplets.where((p) => p.type == MessageFilterChipletType.channel).toList();
    if (chipletsChannel.isNotEmpty) {
      filter.channelIDs = chipletsChannel.map((p) => p.value as String).toList();
    }

    var chipletsSearch = _filterChiplets.where((p) => p.type == MessageFilterChipletType.search).toList();
    if (chipletsSearch.isNotEmpty) {
      filter.searchFilter = chipletsSearch.map((p) => p.value as String).toList();
    }

    var chipletsKeyTokens = _filterChiplets.where((p) => p.type == MessageFilterChipletType.sendkey).toList();
    if (chipletsKeyTokens.isNotEmpty) {
      filter.usedKeys = chipletsKeyTokens.map((p) => p.value as String).toList();
    }

    var chipletPriority = _filterChiplets.where((p) => p.type == MessageFilterChipletType.priority).toList();
    if (chipletPriority.isNotEmpty) {
      filter.priority = chipletPriority.map((p) => p.value as int).toList();
    }

    var chipletSender = _filterChiplets.where((p) => p.type == MessageFilterChipletType.sender).toList();
    if (chipletSender.isNotEmpty) {
      filter.senderNames = chipletSender.map((p) => p.value as String).toList();
    }

    return filter;
  }
}
