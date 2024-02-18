import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';

import '../../models/message.dart';
import '../../state/user_account.dart';
import 'message_list_item.dart';

class MessageListPage extends StatefulWidget {
  const MessageListPage({super.key});

  @override
  State<MessageListPage> createState() => _MessageListPageState();
}

class _MessageListPageState extends State<MessageListPage> {
  static const _pageSize = 20; //TODO

  final PagingController<String, Message> _pagingController = PagingController(firstPageKey: '@start');

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

  Future<void> _fetchPage(String thisPageToken) async {
    final acc = Provider.of<UserAccount>(context, listen: false);

    if (acc.auth == null) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      final [npt, newItems] = await APIClient.getMessageList(acc.auth!, thisPageToken, _pageSize);

      if (npt == '@end') {
        _pagingController.appendLastPage(newItems);
      } else {
        _pagingController.appendPage(newItems, npt);
      }
    } catch (error) {
      _pagingController.error = error;
    }
  }

  @override
  Widget build(BuildContext context) {
    return PagedListView<String, Message>(
      pagingController: _pagingController,
      builderDelegate: PagedChildBuilderDelegate<Message>(
        itemBuilder: (context, item, index) => MessageListItem(
          message: item,
        ),
      ),
    );
  }

  void _createChannel() {
    //TODO
  }
}
