import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/scn_message.dart';
import 'package:simplecloudnotifier/settings/app_settings.dart';

class SCNDataCache {
  SCNDataCache._internal();
  static final SCNDataCache _instance = SCNDataCache._internal();
  factory SCNDataCache() => _instance;

  Future<void> addToMessageCache(List<SCNMessage> newItems) async {
    final cfg = AppSettings();

    final cache = Hive.box<SCNMessage>('scn-message-cache');

    for (var msg in newItems) await cache.put(msg.messageID, msg);

    // delete all but the newest 128 messages

    if (cache.length < cfg.messagePageSize) return;

    final allValues = cache.values.toList();

    allValues.sort((a, b) => -1 * a.timestamp.compareTo(b.timestamp));

    for (var val in allValues.sublist(cfg.messagePageSize)) {
      await cache.delete(val.messageID);
    }
  }

  Future<void> setChannelCache(List<ChannelWithSubscription> channels) async {
    final cache = Hive.box<Channel>('scn-channel-cache');

    if (cache.length != channels.length) await cache.clear();

    for (var chn in channels) await cache.put(chn.channel.channelID, chn.channel);
  }

  bool hasMessagesAndChannels() {
    final chnCache = Hive.box<Channel>('scn-channel-cache');
    final msgCache = Hive.box<SCNMessage>('scn-message-cache');

    return chnCache.isNotEmpty && msgCache.isNotEmpty;
  }

  Map<String, Channel> getChannelMap() {
    final chnCache = Hive.box<Channel>('scn-channel-cache');

    return <String, Channel>{for (var v in chnCache.values) v.channelID: v};
  }

  List<SCNMessage> getMessagesSorted() {
    final msgCache = Hive.box<SCNMessage>('scn-message-cache');

    final cacheMessages = msgCache.values.toList();
    cacheMessages.sort((a, b) => -1 * a.timestamp.compareTo(b.timestamp));

    return cacheMessages;
  }
}
