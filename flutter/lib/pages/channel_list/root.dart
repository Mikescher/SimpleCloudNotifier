import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';

import '../../state/user_account.dart';
import 'channel_list_item.dart';

class ChannelRootPage extends StatefulWidget {
  const ChannelRootPage({super.key});

  @override
  State<ChannelRootPage> createState() => _ChannelRootPageState();
}

class _ChannelRootPageState extends State<ChannelRootPage> {
  final PagingController<int, Channel> _pagingController = PagingController(firstPageKey: 0);

  late UserAccount userAcc;

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
    final acc = Provider.of<UserAccount>(context, listen: false);

    if (acc.auth == null) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      final items = await APIClient.getChannelList(acc.auth!, ChannelSelector.all);

      _pagingController.appendLastPage(items);
    } catch (error) {
      _pagingController.error = error;
    }
  }

  @override
  Widget build(BuildContext context) {
    return PagedListView<int, Channel>(
      pagingController: _pagingController,
      builderDelegate: PagedChildBuilderDelegate<Channel>(
        itemBuilder: (context, item, index) => ChannelListItem(
          channel: item,
        ),
      ),
    );
  }
}
