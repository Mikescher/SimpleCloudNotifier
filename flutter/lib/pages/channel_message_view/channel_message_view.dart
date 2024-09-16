import 'package:flutter/material.dart';
import 'package:infinite_scroll_pagination/infinite_scroll_pagination.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/pages/message_list/message_list_item.dart';
import 'package:simplecloudnotifier/pages/message_view/message_view.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/state/scn_data_cache.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:provider/provider.dart';

class ChannelMessageViewPage extends StatefulWidget {
  const ChannelMessageViewPage({
    required this.channel,
    super.key,
  });

  final Channel channel;

  @override
  State<ChannelMessageViewPage> createState() => _ChannelMessageViewPageState();
}

class _ChannelMessageViewPageState extends State<ChannelMessageViewPage> {
  PagingController<String, SCNMessage> _pagingController = PagingController.fromValue(PagingState(nextPageKey: null, itemList: [], error: null), firstPageKey: '@start');

  @override
  void initState() {
    super.initState();

    _pagingController.addPageRequestListener(_fetchPage);

    _pagingController.refresh();
  }

  @override
  void dispose() {
    _pagingController.dispose();
    super.dispose();
  }

  Future<void> _fetchPage(String thisPageToken) async {
    final acc = Provider.of<AppAuth>(context, listen: false);
    final cfg = Provider.of<AppSettings>(context, listen: false);

    ApplicationLog.debug('Start ChannelMessageViewPage::_pagingController::_fetchPage [ ${thisPageToken} ]');

    if (!acc.isAuth()) {
      _pagingController.error = 'Not logged in';
      return;
    }

    try {
      final (npt, newItems) = await APIClient.getMessageList(acc, thisPageToken, pageSize: cfg.messagePageSize, channelIDs: [this.widget.channel.channelID]);

      SCNDataCache().addToMessageCache(newItems); // no await

      if (npt == '@end') {
        _pagingController.appendLastPage(newItems);
      } else {
        _pagingController.appendPage(newItems, npt);
      }
    } catch (exc, trace) {
      _pagingController.error = exc.toString();
      ApplicationLog.error('Failed to list channel-messages: ' + exc.toString(), trace: trace);
    }
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: this.widget.channel.displayName,
      showSearch: false,
      showShare: false,
      child: _buildMessageList(context),
    );
  }

  Widget _buildMessageList(BuildContext context) {
    return Padding(
      padding: EdgeInsets.fromLTRB(8, 4, 8, 4),
      child: RefreshIndicator(
        onRefresh: () => Future.sync(
          () => _pagingController.refresh(),
        ),
        child: PagedListView<String, SCNMessage>(
          pagingController: _pagingController,
          builderDelegate: PagedChildBuilderDelegate<SCNMessage>(
            itemBuilder: (context, item, index) => MessageListItem(
              message: item,
              allChannels: {this.widget.channel.channelID: this.widget.channel},
              onPressed: () {
                Navi.push(context, () => MessageViewPage(messageID: item.messageID, preloadedData: (item,)));
              },
            ),
          ),
        ),
      ),
    );
  }
}
