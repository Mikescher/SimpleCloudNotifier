import 'package:flutter/material.dart';

import '../../models/channel.dart';

class ChannelListItem extends StatelessWidget {
  const ChannelListItem({
    required this.channel,
    super.key,
  });

  final Channel channel;

  @override
  Widget build(BuildContext context) => ListTile(
        leading: const SizedBox(width: 40, height: 40, child: const Placeholder()),
        title: Text(channel.internalName),
      );
}
