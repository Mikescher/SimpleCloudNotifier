import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:share_plus/share_plus.dart';
import 'package:simplecloudnotifier/api/api_client.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/subscription.dart';
import 'package:simplecloudnotifier/models/user.dart';
import 'package:simplecloudnotifier/pages/channel_message_view/channel_message_view.dart';
import 'package:simplecloudnotifier/state/app_auth.dart';
import 'package:simplecloudnotifier/state/app_bar_state.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:simplecloudnotifier/types/immediate_future.dart';
import 'package:simplecloudnotifier/utils/navi.dart';
import 'package:simplecloudnotifier/utils/toaster.dart';
import 'package:simplecloudnotifier/utils/ui.dart';
import 'package:provider/provider.dart';

class ChannelViewPage extends StatefulWidget {
  const ChannelViewPage({
    required this.channelID,
    required this.preloadedData,
    required this.needsReload,
    super.key,
  });

  final String channelID;
  final (Channel, Subscription?)? preloadedData;

  final void Function()? needsReload;

  @override
  State<ChannelViewPage> createState() => _ChannelViewPageState();
}

enum EditState { none, editing, saving }

enum ChannelViewPageInitState { loading, okay, error }

class _ChannelViewPageState extends State<ChannelViewPage> {
  late ImmediateFuture<String?> _futureSubscribeKey;
  late ImmediateFuture<List<(Subscription, UserPreview?)>> _futureSubscriptions;
  late ImmediateFuture<UserPreview> _futureOwner;

  final TextEditingController _ctrlDisplayName = TextEditingController();
  final TextEditingController _ctrlDescriptionName = TextEditingController();

  int _loadingIndeterminateCounter = 0;

  EditState _editDisplayName = EditState.none;
  String? _displayNameOverride = null;

  EditState _editDescriptionName = EditState.none;
  String? _descriptionNameOverride = null;

  ChannelPreview? channelPreview;
  Channel? channel;
  Subscription? subscription;

  ChannelViewPageInitState loadingState = ChannelViewPageInitState.loading;
  String errorMessage = '';

  @override
  void initState() {
    _initStateAsync();

    super.initState();
  }

  @override
  void _initStateAsync() async {
    final userAcc = Provider.of<AppAuth>(context, listen: false);

    if (widget.preloadedData != null) {
      channelPreview = widget.preloadedData!.$1.toPreview();
      channel = widget.preloadedData!.$1;
      subscription = widget.preloadedData!.$2;
    } else {
      try {
        var p = await APIClient.getChannelPreview(userAcc, widget.channelID);
        channelPreview = p;
        if (p.ownerUserID == userAcc.userID) {
          var r = await APIClient.getChannel(userAcc, widget.channelID);
          channel = r.channel;
          subscription = r.subscription;
        } else {
          channel = null;
          subscription = null; //TODO get own subscription on this channel, even though its foreign channel
        }
      } catch (exc, trace) {
        ApplicationLog.error('Failed to load data: ' + exc.toString(), trace: trace);
        Toaster.error("Error", 'Failed to load data');
        this.errorMessage = 'Failed to load data: ' + exc.toString();
        this.loadingState = ChannelViewPageInitState.error;
        return;
      }
    }

    this.loadingState = ChannelViewPageInitState.okay;

    assert(channelPreview != null);

    if (this.channelPreview!.ownerUserID == userAcc.userID) {
      if (this.channel != null && this.channel!.subscribeKey != null) {
        _futureSubscribeKey = ImmediateFuture<String?>.ofValue(this.channel!.subscribeKey);
      } else {
        _futureSubscribeKey = ImmediateFuture<String?>.ofFuture(_getSubscribeKey(userAcc));
      }
      _futureSubscriptions = ImmediateFuture<List<(Subscription, UserPreview?)>>.ofFuture(_listSubscriptions(userAcc));
    } else {
      _futureSubscribeKey = ImmediateFuture<String?>.ofValue(null);
      _futureSubscriptions = ImmediateFuture<List<(Subscription, UserPreview?)>>.ofValue([]);
    }

    if (this.channelPreview!.ownerUserID == userAcc.userID) {
      var cacheUser = userAcc.getUserOrNull();
      if (cacheUser != null) {
        _futureOwner = ImmediateFuture<UserPreview>.ofValue(cacheUser.toPreview());
      } else {
        _futureOwner = ImmediateFuture<UserPreview>.ofFuture(_getOwner(userAcc));
      }
    } else {
      _futureOwner = ImmediateFuture<UserPreview>.ofFuture(APIClient.getUserPreview(userAcc, this.channelPreview!.ownerUserID));
    }
  }

  @override
  void dispose() {
    _ctrlDisplayName.dispose();
    _ctrlDescriptionName.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final userAcc = Provider.of<AppAuth>(context, listen: false);

    Widget child;

    if (loadingState == ChannelViewPageInitState.loading) {
      child = Center(child: CircularProgressIndicator());
    } else if (loadingState == ChannelViewPageInitState.error) {
      child = Center(child: Text('Error: ' + errorMessage)); //TODO better error
    } else if (loadingState == ChannelViewPageInitState.okay && channelPreview!.ownerUserID == userAcc.userID) {
      child = _buildOwnedChannelView(context, this.channel!);
    } else {
      child = _buildForeignChannelView(context, this.channelPreview!);
    }

    return SCNScaffold(
      title: 'Channel',
      showSearch: false,
      showShare: false,
      child: child,
    );
  }

  Widget _buildOwnedChannelView(BuildContext context, Channel channel) {
    final isSubscribed = (subscription != null && subscription!.confirmed);

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _buildQRCode(context),
            SizedBox(height: 8),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidIdCardClip,
              title: 'ChannelID',
              values: [channel.channelID],
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidInputNumeric,
              title: 'InternalName',
              values: [channel.internalName],
            ),
            _buildDisplayNameCard(context, true),
            _buildDescriptionNameCard(context, true),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidDiagramSubtask,
              title: 'Subscription (own)',
              values: [_formatSubscriptionStatus(this.subscription)],
              iconActions: isSubscribed ? [(FontAwesomeIcons.solidSquareXmark, _unsubscribe)] : [(FontAwesomeIcons.solidSquareRss, _subscribe)],
            ),
            _buildForeignSubscriptions(context),
            _buildOwnerCard(context, true),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidEnvelope,
              title: 'Messages',
              values: [channel.messagesSent.toString()],
              mainAction: () {
                Navi.push(context, () => ChannelMessageViewPage(channel: channel));
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildForeignChannelView(BuildContext context, ChannelPreview channel) {
    final isSubscribed = (subscription != null && subscription!.confirmed);

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            SizedBox(height: 8),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidIdCardClip,
              title: 'ChannelID',
              values: [channel.channelID],
            ),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidInputNumeric,
              title: 'InternalName',
              values: [channel.internalName],
            ),
            _buildDisplayNameCard(context, false),
            _buildDescriptionNameCard(context, false),
            UI.metaCard(
              context: context,
              icon: FontAwesomeIcons.solidDiagramSubtask,
              title: 'Subscription (own)',
              values: [_formatSubscriptionStatus(subscription)],
              iconActions: isSubscribed ? [(FontAwesomeIcons.solidSquareXmark, _unsubscribe)] : [(FontAwesomeIcons.solidSquareRss, _subscribe)],
            ),
            _buildForeignSubscriptions(context),
            _buildOwnerCard(context, false),
          ],
        ),
      ),
    );
  }

  Widget _buildForeignSubscriptions(BuildContext context) {
    return FutureBuilder(
      future: _futureSubscriptions.future,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          return Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              for (final (sub, user) in snapshot.data!.where((v) => v.$1.subscriptionID != subscription?.subscriptionID))
                UI.metaCard(
                  context: context,
                  icon: FontAwesomeIcons.solidDiagramSuccessor,
                  title: 'Subscription (' + (user?.username ?? user?.userID ?? 'other') + ')',
                  values: [_formatSubscriptionStatus(sub)],
                  iconActions: _getForeignSubActions(sub),
                ),
            ],
          );
        } else {
          return SizedBox();
        }
      },
    );
  }

  Widget _buildOwnerCard(BuildContext context, bool isOwned) {
    return FutureBuilder(
      future: _futureOwner.future,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          return UI.metaCard(
            context: context,
            icon: FontAwesomeIcons.solidUser,
            title: 'Owner',
            values: [channelPreview!.ownerUserID + (isOwned ? ' (you)' : ''), if (snapshot.data?.username != null) snapshot.data!.username!],
          );
        } else {
          return UI.metaCard(
            context: context,
            icon: FontAwesomeIcons.solidUser,
            title: 'Owner',
            values: [channelPreview!.ownerUserID + (isOwned ? ' (you)' : '')],
          );
        }
      },
    );
  }

  Widget _buildQRCode(BuildContext context) {
    return FutureBuilder(
      future: _futureSubscribeKey.future,
      builder: (context, snapshot) {
        if (snapshot.hasData && snapshot.data != null) {
          var text = 'TODO' + '\n' + channel!.channelID + '\n' + snapshot.data!; //TODO deeplink-y (also perhaps just bas64 everything together?)
          return GestureDetector(
            onTap: () {
              Share.share(text, subject: _displayNameOverride ?? channel!.displayName);
            },
            child: Center(
              child: QrImageView(
                data: text,
                version: QrVersions.auto,
                size: 300.0,
                eyeStyle: QrEyeStyle(
                  eyeShape: QrEyeShape.square,
                  color: Theme.of(context).textTheme.bodyLarge?.color,
                ),
                dataModuleStyle: QrDataModuleStyle(
                  dataModuleShape: QrDataModuleShape.square,
                  color: Theme.of(context).textTheme.bodyLarge?.color,
                ),
              ),
            ),
          );
        } else if (snapshot.hasData && snapshot.data == null) {
          return const SizedBox(
            width: 300.0,
            height: 300.0,
            child: Center(child: Icon(FontAwesomeIcons.solidSnake, size: 64)),
          );
        } else {
          return const SizedBox(
            width: 300.0,
            height: 300.0,
            child: Center(child: CircularProgressIndicator()),
          );
        }
      },
    );
  }

  Widget _buildDisplayNameCard(BuildContext context, bool isOwned) {
    if (_editDisplayName == EditState.editing) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: UI.box(
          context: context,
          padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
          child: Row(
            children: [
              Container(child: Center(child: FaIcon(FontAwesomeIcons.solidInputText, size: 18)), height: 43),
              SizedBox(width: 16),
              Expanded(
                child: TextField(
                  autofocus: true,
                  controller: _ctrlDisplayName,
                  decoration: new InputDecoration.collapsed(hintText: 'DisplayName'),
                ),
              ),
              SizedBox(width: 12),
              SizedBox(width: 4),
              IconButton(icon: FaIcon(FontAwesomeIcons.solidFloppyDisk), onPressed: _saveDisplayName),
            ],
          ),
        ),
      );
    } else if (_editDisplayName == EditState.none) {
      return UI.metaCard(
        context: context,
        icon: FontAwesomeIcons.solidInputText,
        title: 'DisplayName',
        values: [_displayNameOverride ?? channelPreview!.displayName],
        iconActions: isOwned ? [(FontAwesomeIcons.penToSquare, _showEditDisplayName)] : [],
      );
    } else if (_editDisplayName == EditState.saving) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: UI.box(
          context: context,
          padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
          child: Row(
            children: [
              Container(child: Center(child: FaIcon(FontAwesomeIcons.solidInputText, size: 18)), height: 43),
              SizedBox(width: 16),
              Expanded(child: SizedBox()),
              SizedBox(width: 12),
              SizedBox(width: 4),
              Padding(padding: const EdgeInsets.all(8.0), child: SizedBox(width: 18, height: 18, child: CircularProgressIndicator())),
            ],
          ),
        ),
      );
    } else {
      throw 'Invalid EditDisplayNameState: $_editDisplayName';
    }
  }

  Widget _buildDescriptionNameCard(BuildContext context, bool isOwned) {
    if (_editDescriptionName == EditState.editing) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: UI.box(
          context: context,
          padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
          child: Row(
            children: [
              Container(child: Center(child: FaIcon(FontAwesomeIcons.solidInputPipe, size: 18)), height: 43),
              SizedBox(width: 16),
              Expanded(
                child: TextField(
                  autofocus: true,
                  controller: _ctrlDescriptionName,
                  decoration: new InputDecoration.collapsed(hintText: 'Description'),
                ),
              ),
              SizedBox(width: 12),
              SizedBox(width: 4),
              IconButton(icon: FaIcon(FontAwesomeIcons.solidFloppyDisk), onPressed: _saveDescriptionName),
            ],
          ),
        ),
      );
    } else if (_editDescriptionName == EditState.none) {
      return UI.metaCard(
        context: context,
        icon: FontAwesomeIcons.solidInputPipe,
        title: 'Description',
        values: [_descriptionNameOverride ?? channelPreview?.descriptionName ?? ''],
        iconActions: isOwned ? [(FontAwesomeIcons.penToSquare, _showEditDescriptionName)] : [],
      );
    } else if (_editDescriptionName == EditState.saving) {
      return Padding(
        padding: EdgeInsets.symmetric(vertical: 4, horizontal: 0),
        child: UI.box(
          context: context,
          padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
          child: Row(
            children: [
              Container(child: Center(child: FaIcon(FontAwesomeIcons.solidInputPipe, size: 18)), height: 43),
              SizedBox(width: 16),
              Expanded(child: SizedBox()),
              SizedBox(width: 12),
              SizedBox(width: 4),
              Padding(padding: const EdgeInsets.all(8.0), child: SizedBox(width: 18, height: 18, child: CircularProgressIndicator())),
            ],
          ),
        ),
      );
    } else {
      throw 'Invalid EditDescriptionNameState: $_editDescriptionName';
    }
  }

  void _subscribe() {
    //TODO
  }

  void _unsubscribe() {
    //TODO
  }

  void _showEditDisplayName() {
    setState(() {
      _ctrlDisplayName.text = _displayNameOverride ?? channelPreview?.displayName ?? '';
      _editDisplayName = EditState.editing;
      if (_editDescriptionName == EditState.editing) _editDescriptionName = EditState.none;
    });
  }

  void _saveDisplayName() async {
    final userAcc = Provider.of<AppAuth>(context, listen: false);

    final newName = _ctrlDisplayName.text;

    try {
      setState(() {
        _editDisplayName = EditState.saving;
      });

      final newChannel = await APIClient.updateChannel(userAcc, widget.channelID, displayName: newName);

      setState(() {
        _editDisplayName = EditState.none;
        _displayNameOverride = newChannel.channel.displayName;
      });

      widget.needsReload?.call();
    } catch (exc, trace) {
      ApplicationLog.error('Failed to save DisplayName: ' + exc.toString(), trace: trace);
      Toaster.error("Error", 'Failed to save DisplayName');
    }
  }

  void _showEditDescriptionName() {
    setState(() {
      _ctrlDescriptionName.text = _descriptionNameOverride ?? channelPreview?.descriptionName ?? '';
      _editDescriptionName = EditState.editing;
      if (_editDisplayName == EditState.editing) _editDisplayName = EditState.none;
    });
  }

  void _saveDescriptionName() async {
    final userAcc = Provider.of<AppAuth>(context, listen: false);

    final newName = _ctrlDescriptionName.text;

    try {
      setState(() {
        _editDescriptionName = EditState.saving;
      });

      final newChannel = await APIClient.updateChannel(userAcc, widget.channelID, descriptionName: newName);

      setState(() {
        _editDescriptionName = EditState.none;
        _descriptionNameOverride = newChannel.channel.descriptionName ?? '';
      });

      widget.needsReload?.call();
    } catch (exc, trace) {
      ApplicationLog.error('Failed to save DescriptionName: ' + exc.toString(), trace: trace);
      Toaster.error("Error", 'Failed to save DescriptionName');
    }
  }

  void _cancelForeignSubscription(Subscription sub) {
    //TODO
  }

  void _confirmForeignSubscription(Subscription sub) {
    //TODO
  }

  void _denyForeignSubscription(Subscription sub) {
    //TODO
  }

  String _formatSubscriptionStatus(Subscription? subscription) {
    if (subscription == null) {
      return 'Not Subscribed';
    } else if (subscription.confirmed) {
      return 'Subscribed';
    } else {
      return 'Requested';
    }
  }

  Future<String?> _getSubscribeKey(AppAuth auth) async {
    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      _incLoadingIndeterminateCounter(1);

      var channel = await APIClient.getChannel(auth, widget.channelID);

      //await Future.delayed(const Duration(seconds: 10), () {});

      return channel.channel.subscribeKey;
    } finally {
      _incLoadingIndeterminateCounter(-1);
    }
  }

  Future<List<(Subscription, UserPreview?)>> _listSubscriptions(AppAuth auth) async {
    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      _incLoadingIndeterminateCounter(1);

      var subs = await APIClient.getChannelSubscriptions(auth, widget.channelID);

      var userMap = {for (var v in (await Future.wait(subs.map((e) => e.subscriberUserID).toSet().map((e) => APIClient.getUserPreview(auth, e)).toList()))) v.userID: v};

      //await Future.delayed(const Duration(seconds: 10), () {});

      return subs.map((e) => (e, userMap[e.subscriberUserID] ?? null)).toList();
    } finally {
      _incLoadingIndeterminateCounter(-1);
    }
  }

  Future<UserPreview> _getOwner(AppAuth auth) async {
    try {
      await Future.delayed(const Duration(seconds: 0), () {}); // this is annoyingly important - otherwise we call setLoadingIndeterminate directly in initStat() and get an exception....

      _incLoadingIndeterminateCounter(1);

      final owner = APIClient.getUserPreview(auth, channelPreview!.ownerUserID);

      //await Future.delayed(const Duration(seconds: 10), () {});

      return owner;
    } finally {
      _incLoadingIndeterminateCounter(-1);
    }
  }

  List<(IconData, void Function())> _getForeignSubActions(Subscription sub) {
    if (sub.confirmed) {
      return [(FontAwesomeIcons.solidSquareXmark, () => _cancelForeignSubscription(sub))];
    } else {
      return [
        (FontAwesomeIcons.solidSquareCheck, () => _confirmForeignSubscription(sub)),
        (FontAwesomeIcons.solidSquareXmark, () => _denyForeignSubscription(sub)),
      ];
    }
  }

  void _incLoadingIndeterminateCounter(int delta) {
    setState(() {
      _loadingIndeterminateCounter += delta;
      AppBarState().setLoadingIndeterminate(_loadingIndeterminateCounter > 0);
    });
  }
}
