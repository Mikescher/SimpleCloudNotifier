import 'dart:math';

import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/models/channel.dart';
import 'package:simplecloudnotifier/models/message.dart';
import 'package:intl/intl.dart';

class MessageListItem extends StatelessWidget {
  static final _dateFormat = DateFormat('yyyy-MM-dd kk:mm');
  static final _lineCount = 3; //TODO setting

  const MessageListItem({
    required this.message,
    required this.allChannels,
    required this.onPressed,
    super.key,
  });

  final Message message;
  final Map<String, Channel> allChannels;
  final Null Function() onPressed;

  @override
  Widget build(BuildContext context) {
    if (showChannel(message)) {
      return buildWithChannel(context);
    } else {
      return buildWithoutChannel(context);
    }
  }

  Card buildWithoutChannel(BuildContext context) {
    return Card.filled(
      margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
      shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
      color: (message.priority == 2) ? Theme.of(context).colorScheme.errorContainer : Theme.of(context).cardTheme.color,
      //clipBehavior: Clip.hardEdge,  // nto needed, because our borderRadius is 0 anyway
      child: InkWell(
        splashColor: Theme.of(context).splashColor,
        onTap: onPressed,
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (message.priority == 2) FaIcon(FontAwesomeIcons.solidTriangleExclamation, size: 16, color: Colors.red[900]),
                  if (message.priority == 2) SizedBox(width: 4),
                  if (message.priority == 0) FaIcon(FontAwesomeIcons.solidDown, size: 16, color: Colors.lightBlue[900]),
                  if (message.priority == 0) SizedBox(width: 4),
                  Expanded(
                    child: Text(
                      processTitle(message.title),
                      style: const TextStyle(fontWeight: FontWeight.bold),
                      overflow: TextOverflow.ellipsis,
                      maxLines: 3,
                    ),
                  ),
                  Text(
                    _dateFormat.format(DateTime.parse(message.timestamp).toLocal()),
                    style: const TextStyle(fontWeight: FontWeight.normal, fontSize: 11),
                    overflow: TextOverflow.clip,
                    maxLines: 1,
                  ),
                ],
              ),
              SizedBox(height: 4),
              Text(
                processContent(message.content),
                style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)),
                overflow: TextOverflow.ellipsis,
                maxLines: _lineCount,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Card buildWithChannel(BuildContext context) {
    return Card.filled(
      margin: EdgeInsets.fromLTRB(0, 4, 0, 4),
      shape: BeveledRectangleBorder(borderRadius: BorderRadius.circular(0)),
      //clipBehavior: Clip.hardEdge,  // nto needed, because our borderRadius is 0 anyway
      color: (message.priority == 2) ? Theme.of(context).colorScheme.errorContainer : Theme.of(context).cardTheme.color,
      child: InkWell(
        splashColor: Theme.of(context).splashColor,
        onTap: onPressed,
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (message.priority == 2) FaIcon(FontAwesomeIcons.solidTriangleExclamation, size: 16, color: Colors.red[900]),
                  if (message.priority == 2) SizedBox(width: 4),
                  if (message.priority == 0) FaIcon(FontAwesomeIcons.solidDown, size: 16, color: Colors.lightBlue[900]),
                  if (message.priority == 0) SizedBox(width: 4),
                  Container(
                    padding: const EdgeInsets.fromLTRB(4, 0, 4, 0),
                    margin: const EdgeInsets.fromLTRB(0, 0, 4, 0),
                    decoration: BoxDecoration(
                      color: Theme.of(context).hintColor,
                      borderRadius: BorderRadius.all(Radius.circular(4)),
                    ),
                    child: Text(
                      resolveChannelName(message),
                      style: TextStyle(fontWeight: FontWeight.bold, color: Theme.of(context).cardColor, fontSize: 12),
                      overflow: TextOverflow.clip,
                      maxLines: 1,
                    ),
                  ),
                  Expanded(child: SizedBox()),
                  Text(
                    _dateFormat.format(DateTime.parse(message.timestamp).toLocal()),
                    style: const TextStyle(fontWeight: FontWeight.normal, fontSize: 11),
                    overflow: TextOverflow.clip,
                    maxLines: 1,
                  ),
                ],
              ),
              SizedBox(height: 4),
              Text(
                processTitle(message.title),
                style: const TextStyle(fontWeight: FontWeight.bold),
                overflow: TextOverflow.ellipsis,
                maxLines: 3,
              ),
              Text(
                processContent(message.content),
                style: TextStyle(color: Theme.of(context).textTheme.bodyLarge?.color?.withAlpha(160)),
                overflow: TextOverflow.ellipsis,
                maxLines: _lineCount,
              ),
            ],
          ),
        ),
      ),
    );
  }

  String processContent(String? v) {
    if (v == null) {
      return '';
    }

    var lines = v.split('\n');
    if (lines.isEmpty) {
      return '';
    }

    return lines.sublist(0, min(_lineCount, lines.length)).join("\n").trim();
  }

  String processTitle(String? v) {
    if (v == null) {
      return '';
    }

    v = v.replaceAll("\n", " ");
    v = v.replaceAll("\t", " ");
    v = v.replaceAll("\r", "");

    return v;
  }

  String resolveChannelName(Message message) {
    return allChannels[message.channelID]?.displayName ?? message.channelInternalName;
  }

  bool showChannel(Message message) {
    return message.channelInternalName != 'main';
  }
}
