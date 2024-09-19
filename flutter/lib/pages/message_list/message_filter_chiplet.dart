import 'package:flutter/src/widgets/icon_data.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

enum MessageFilterChipletType {
  search,
  channel,
  sender,
  timeRange,
  priority,
  sendkey,
}

class MessageFilterChiplet {
  final String label; // display value
  final dynamic value; // search/api value
  final MessageFilterChipletType type;

  MessageFilterChiplet({required this.label, required this.value, required this.type});

  IconData? icon() {
    switch (type) {
      case MessageFilterChipletType.search:
        return FontAwesomeIcons.magnifyingGlass;
      case MessageFilterChipletType.channel:
        return FontAwesomeIcons.snake;
      case MessageFilterChipletType.sender:
        return FontAwesomeIcons.signature;
      case MessageFilterChipletType.timeRange:
        return FontAwesomeIcons.timer;
      case MessageFilterChipletType.priority:
        return FontAwesomeIcons.bolt;
      case MessageFilterChipletType.sendkey:
        return FontAwesomeIcons.gearCode;
    }
  }
}
