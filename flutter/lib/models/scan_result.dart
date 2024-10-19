import 'package:simplecloudnotifier/models/channel.dart';

enum ScanResultMode { ChannelSubscribe, MessageSend, Channel }

abstract class ScanResult {
  ScanResultMode get mode;

  static ScanResult? parse(String v) {
    var lines = v.split('\n');

    if (lines.length == 1 && lines[0].startsWith('https://simplecloudnotifier.de?')) {
      final v = Uri.tryParse(lines[0]);

      if (v != null && v.queryParameters.containsKey('preset_user_id') && v.queryParameters.containsKey('preset_user_key')) {
        return ScanResultMessageSend(userID: v.queryParameters['preset_user_id']!, userKey: v.queryParameters['preset_user_key']);
      }
      if (v != null && v.queryParameters.containsKey('preset_user_id') && v.queryParameters.containsKey('preset_user_key')) {
        return ScanResultMessageSend(userID: v.queryParameters['preset_user_id']!, userKey: null);
      }
    }

    if (lines.length == 6 && lines[0] == '@scn.channel.subscribe' && lines[1] == 'v1') {
      return ScanResultChannelSubscribe(channelDisplayName: lines[2], ownerUserID: lines[3], channelID: lines[4], subscribeKey: lines[5]);
    }

    if (lines.length == 5 && lines[0] == '@scn.channel' && lines[1] == 'v1') {
      if (lines.length != 4) return null;

      return ScanResultChannel(channelDisplayName: lines[2], ownerUserID: lines[3], channelID: lines[4]);
    }

    return null;
  }

  static String createChannelQR(Channel channel) {
    return '@scn.channel' + '\n' + "v1" + '\n' + channel.displayName + '\n' + channel.ownerUserID + '\n' + channel.channelID;
  }

  static String createChannelSubscribeQR(Channel channel, String subscribeKey) {
    return '@scn.channel.subscribe' + '\n' + "v1" + '\n' + channel.displayName + '\n' + channel.ownerUserID + '\n' + channel.channelID + '\n' + subscribeKey;
  }
}

class ScanResultMessageSend extends ScanResult {
  final String userID;
  final String? userKey;

  ScanResultMessageSend({required this.userID, required this.userKey});

  @override
  ScanResultMode get mode => ScanResultMode.MessageSend;
}

class ScanResultChannel extends ScanResult {
  final String channelDisplayName;
  final String ownerUserID;
  final String channelID;

  ScanResultChannel({required this.channelDisplayName, required this.ownerUserID, required this.channelID});

  @override
  ScanResultMode get mode => ScanResultMode.Channel;
}

class ScanResultChannelSubscribe extends ScanResult {
  final String channelDisplayName;
  final String ownerUserID;
  final String channelID;
  final String subscribeKey;

  ScanResultChannelSubscribe({required this.channelDisplayName, required this.ownerUserID, required this.channelID, required this.subscribeKey});

  @override
  ScanResultMode get mode => ScanResultMode.ChannelSubscribe;
}
