import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/pages/channel_list/channel_list_item.dart';

class ChannelRootPage extends StatefulWidget {
  const ChannelRootPage({super.key});

  @override
  State<ChannelRootPage> createState() => _ChannelRootPageState();
}

class _ChannelRootPageState extends State<ChannelRootPage> {
  final PagingController<int, Channel> _pagingController = PagingController(firstPageKey: 0);

  @override
  void initState() {
    _pagingController.addPageRequestListener((pageKey) {
      _fetchPage(pageKey);
    });
    super.initState();
  }

  @override
  void dispose() {
    _pagingController.dispose();
    super.dispose();
  }

  Future<void> _fetchPage(int pageKey) async {
    final acc = Provider.of<AppAuth>(context, listen: false);

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      final items = (await APIClient.getChannelList(acc, ChannelSelector.all)).map((p) => p.channel).toList();

      items.sort((a, b) => -1 * (a.timestampLastSent ?? '').compareTo(b.timestampLastSent ?? ''));

      _pagingController.appendLastPage(items);
    } catch (exc, trace) {
      _pagingController.error = exc.toString();
      ApplicationLog.error('Failed to list channels: ' + exc.toString(), trace: trace);
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
