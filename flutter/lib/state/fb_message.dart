import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:simplecloudnotifier/state/interfaces.dart';

part 'fb_message.g.dart';

class FBMessageLog {
  //TODO max size, auto clear old

  static void insert(RemoteMessage msg) {
    Hive.box<FBMessage>('scn-fb-messages').add(FBMessage.fromRemoteMessage(msg));
  }
}

@HiveType(typeId: 106)
class FBMessage extends HiveObject implements FieldDebuggable {
  @HiveField(0)
  final String? senderId;
  @HiveField(1)
  final String? category;
  @HiveField(2)
  final String? collapseKey;
  @HiveField(3)
  final bool contentAvailable;
  @HiveField(4)
  final Map<String, String> data;
  @HiveField(5)
  final String? from;
  @HiveField(6)
  final String? messageId;
  @HiveField(7)
  final String? messageType;
  @HiveField(8)
  final bool mutableContent;
  @HiveField(9)
  final RemoteNotification? notification;
  @HiveField(10)
  final DateTime? sentTime;
  @HiveField(11)
  final String? threadId;
  @HiveField(12)
  final int? ttl;

  @HiveField(20)
  final String? notificationAndroidChannelId;
  @HiveField(21)
  final String? notificationAndroidClickAction;
  @HiveField(22)
  final String? notificationAndroidColor;
  @HiveField(23)
  final int? notificationAndroidCount;
  @HiveField(24)
  final String? notificationAndroidImageUrl;
  @HiveField(25)
  final String? notificationAndroidLink;
  @HiveField(26)
  final AndroidNotificationPriority? notificationAndroidPriority;
  @HiveField(27)
  final String? notificationAndroidSmallIcon;
  @HiveField(28)
  final String? notificationAndroidSound;
  @HiveField(29)
  final String? notificationAndroidTicker;
  @HiveField(30)
  final AndroidNotificationVisibility? notificationAndroidVisibility;
  @HiveField(31)
  final String? notificationAndroidTag;

  @HiveField(40)
  final String? notificationAppleBadge;
  @HiveField(41)
  final AppleNotificationSound? notificationAppleSound;
  @HiveField(42)
  final String? notificationAppleImageUrl;
  @HiveField(43)
  final String? notificationAppleSubtitle;
  @HiveField(44)
  final List<String>? notificationAppleSubtitleLocArgs;
  @HiveField(45)
  final String? notificationAppleSubtitleLocKey;

  @HiveField(50)
  final String? notificationWebAnalyticsLabel;
  @HiveField(51)
  final String? notificationWebImage;
  @HiveField(52)
  final String? notificationWebLink;

  @HiveField(60)
  final String? notificationTitle;
  @HiveField(61)
  final List<String>? notificationTitleLocArgs;
  @HiveField(62)
  final String? notificationTitleLocKey;
  @HiveField(63)
  final String? notificationBody;
  @HiveField(64)
  final List<String>? notificationBodyLocArgs;
  @HiveField(65)
  final String? notificationBodyLocKey;

  FBMessage({
    required this.senderId,
    required this.category,
    required this.collapseKey,
    required this.contentAvailable,
    required this.data,
    required this.from,
    required this.messageId,
    required this.messageType,
    required this.mutableContent,
    required this.notification,
    required this.sentTime,
    required this.threadId,
    required this.ttl,
    required this.notificationAndroidChannelId,
    required this.notificationAndroidClickAction,
    required this.notificationAndroidColor,
    required this.notificationAndroidCount,
    required this.notificationAndroidImageUrl,
    required this.notificationAndroidLink,
    required this.notificationAndroidPriority,
    required this.notificationAndroidSmallIcon,
    required this.notificationAndroidSound,
    required this.notificationAndroidTicker,
    required this.notificationAndroidVisibility,
    required this.notificationAndroidTag,
    required this.notificationAppleBadge,
    required this.notificationAppleSound,
    required this.notificationAppleImageUrl,
    required this.notificationAppleSubtitle,
    required this.notificationAppleSubtitleLocArgs,
    required this.notificationAppleSubtitleLocKey,
    required this.notificationWebAnalyticsLabel,
    required this.notificationWebImage,
    required this.notificationWebLink,
    required this.notificationTitle,
    required this.notificationTitleLocArgs,
    required this.notificationTitleLocKey,
    required this.notificationBody,
    required this.notificationBodyLocArgs,
    required this.notificationBodyLocKey,
  });

  FBMessage.fromRemoteMessage(RemoteMessage rmsg)
      : this.senderId = rmsg.senderId,
        this.category = rmsg.category,
        this.collapseKey = rmsg.collapseKey,
        this.contentAvailable = rmsg.contentAvailable,
        this.data = rmsg.data.map((key, value) => MapEntry(key, value?.toString() ?? '')),
        this.from = rmsg.from,
        this.messageId = rmsg.messageId,
        this.messageType = rmsg.messageType,
        this.mutableContent = rmsg.mutableContent,
        this.notification = rmsg.notification,
        this.sentTime = rmsg.sentTime,
        this.threadId = rmsg.threadId,
        this.ttl = rmsg.ttl,
        this.notificationAndroidChannelId = rmsg.notification?.android?.channelId,
        this.notificationAndroidClickAction = rmsg.notification?.android?.clickAction,
        this.notificationAndroidColor = rmsg.notification?.android?.color,
        this.notificationAndroidCount = rmsg.notification?.android?.count,
        this.notificationAndroidImageUrl = rmsg.notification?.android?.imageUrl,
        this.notificationAndroidLink = rmsg.notification?.android?.link,
        this.notificationAndroidPriority = rmsg.notification?.android?.priority,
        this.notificationAndroidSmallIcon = rmsg.notification?.android?.smallIcon,
        this.notificationAndroidSound = rmsg.notification?.android?.sound,
        this.notificationAndroidTicker = rmsg.notification?.android?.ticker,
        this.notificationAndroidVisibility = rmsg.notification?.android?.visibility,
        this.notificationAndroidTag = rmsg.notification?.android?.tag,
        this.notificationAppleBadge = rmsg.notification?.apple?.badge,
        this.notificationAppleSound = rmsg.notification?.apple?.sound,
        this.notificationAppleImageUrl = rmsg.notification?.apple?.imageUrl,
        this.notificationAppleSubtitle = rmsg.notification?.apple?.subtitle,
        this.notificationAppleSubtitleLocArgs = rmsg.notification?.apple?.subtitleLocArgs,
        this.notificationAppleSubtitleLocKey = rmsg.notification?.apple?.subtitleLocKey,
        this.notificationWebAnalyticsLabel = rmsg.notification?.web?.analyticsLabel,
        this.notificationWebImage = rmsg.notification?.web?.image,
        this.notificationWebLink = rmsg.notification?.web?.link,
        this.notificationTitle = rmsg.notification?.title,
        this.notificationTitleLocArgs = rmsg.notification?.titleLocArgs,
        this.notificationTitleLocKey = rmsg.notification?.titleLocKey,
        this.notificationBody = rmsg.notification?.body,
        this.notificationBodyLocArgs = rmsg.notification?.bodyLocArgs,
        this.notificationBodyLocKey = rmsg.notification?.bodyLocKey {}

  @override
  String toString() {
    return 'FBMessage[${this.messageId ?? 'NULL'}]';
  }

  List<(String, String)> debugFieldList() {
    return [
      ('senderId', this.senderId ?? ''),
      ('category', this.category ?? ''),
      ('collapseKey', this.collapseKey ?? ''),
      ('contentAvailable', this.contentAvailable.toString()),
      ('data', this.data.toString()),
      ('from', this.from ?? ''),
      ('messageId', this.messageId ?? ''),
      ('messageType', this.messageType ?? ''),
      ('mutableContent', this.mutableContent.toString()),
      ('notification', this.notification?.toString() ?? ''),
      ('sentTime', this.sentTime?.toString() ?? ''),
      ('threadId', this.threadId ?? ''),
      ('ttl', this.ttl?.toString() ?? ''),
      ('notification.Android.ChannelId', this.notificationAndroidChannelId ?? ''),
      ('notification.Android.ClickAction', this.notificationAndroidClickAction ?? ''),
      ('notification.Android.Color', this.notificationAndroidColor ?? ''),
      ('notification.Android.Count', this.notificationAndroidCount?.toString() ?? ''),
      ('notification.Android.ImageUrl', this.notificationAndroidImageUrl ?? ''),
      ('notification.Android.Link', this.notificationAndroidLink ?? ''),
      ('notification.Android.Priority', this.notificationAndroidPriority?.toString() ?? ''),
      ('notification.Android.SmallIcon', this.notificationAndroidSmallIcon ?? ''),
      ('notification.Android.Sound', this.notificationAndroidSound ?? ''),
      ('notification.Android.Ticker', this.notificationAndroidTicker ?? ''),
      ('notification.Android.Visibility', this.notificationAndroidVisibility?.toString() ?? ''),
      ('notification.Android.Tag', this.notificationAndroidTag ?? ''),
      ('notification.Apple.Badge', this.notificationAppleBadge ?? ''),
      ('notification.Apple.Sound', this.notificationAppleSound?.toString() ?? ''),
      ('notification.Apple.ImageUrl', this.notificationAppleImageUrl ?? ''),
      ('notification.Apple.Subtitle', this.notificationAppleSubtitle ?? ''),
      ('notification.Apple.SubtitleLocArgs', this.notificationAppleSubtitleLocArgs?.toString() ?? ''),
      ('notification.Apple.SubtitleLocKey', this.notificationAppleSubtitleLocKey ?? ''),
      ('notification.Web.AnalyticsLabel', this.notificationWebAnalyticsLabel ?? ''),
      ('notification.Web.Image', this.notificationWebImage ?? ''),
      ('notification.Web.Link', this.notificationWebLink ?? ''),
      ('notification.Title', this.notificationTitle ?? ''),
      ('notification.TitleLocArgs', this.notificationTitleLocArgs?.toString() ?? ''),
      ('notification.TitleLocKey', this.notificationTitleLocKey ?? ''),
      ('notification.Body', this.notificationBody ?? ''),
      ('notification.BodyLocArgs', this.notificationBodyLocArgs?.toString() ?? ''),
      ('notification.BodyLocKey', this.notificationBodyLocKey ?? ''),
    ];
  }
}
