import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:provider/provider.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:simplecloudnotifier/pages/message_view/message_view.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

import 'message_list_item.dart';

class MessageListPage extends StatefulWidget {
  const MessageListPage({super.key});

  @override
  State<MessageListPage> createState() => _MessageListPageState();
}

class _MessageListPageState extends State<MessageListPage> {
  static const _pageSize = 128;

  final PagingController<String, Message> _pagingController = PagingController(firstPageKey: '@start');

  Map<String, Channel>? _channels = null;

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
      if (_channels == null) {
        final channels = await APIClient.getChannelList(acc.auth!, ChannelSelector.allAny);
        _channels = Map.fromIterable(channels, key: (e) => e.channelID);
      }

      final (npt, newItems) = await APIClient.getMessageList(acc.auth!, thisPageToken, _pageSize);

      if (npt == '@end') {
        _pagingController.appendLastPage(newItems);
      } else {
        _pagingController.appendPage(newItems, npt);
      }
    } catch (error) {
      print("API-Error: "); //TODO remove me, proper error handling
      print(error); //TODO remove me, proper error handling
      _pagingController.error = error;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.fromLTRB(8, 4, 8, 4),
      child: PagedListView<String, Message>(
        pagingController: _pagingController,
        builderDelegate: PagedChildBuilderDelegate<Message>(
          itemBuilder: (context, item, index) => MessageListItem(
            message: item,
            allChannels: _channels ?? {},
            onPressed: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => MessageViewPage(message: item)));
            },
          ),
        ),
      ),
    );
  }
}
