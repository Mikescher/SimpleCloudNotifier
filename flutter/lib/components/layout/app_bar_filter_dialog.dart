import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:simplecloudnotifier/utils/navi.dart';

class AppBarFilterDialog extends StatefulWidget {
  @override
  _AppBarFilterDialogState createState() => _AppBarFilterDialogState();
}

class _AppBarFilterDialogState extends State<AppBarFilterDialog> {
  double _height = 0;

  double _targetHeight = 4 + (48 * 6) + (16 * 5) + 4;

  @override
  void initState() {
    super.initState();
    Future.delayed(Duration.zero, () {
      setState(() {
        _height = _targetHeight;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    double vpWidth = MediaQuery.sizeOf(context).width;

    return Container(
      margin: const EdgeInsets.all(0),
      width: vpWidth,
      color: Colors.transparent,
      child: Column(
        children: [
          Container(
            color: Theme.of(context).secondaryHeaderColor,
            child: AnimatedContainer(
              duration: Duration(milliseconds: 350),
              curve: Curves.easeInCubic,
              height: _height,
              child: ClipRect(
                child: OverflowBox(
                  alignment: Alignment.topCenter,
                  maxWidth: vpWidth,
                  minWidth: vpWidth,
                  minHeight: 0,
                  maxHeight: _targetHeight,
                  child: Column(
                    children: [
                      SizedBox(height: 4),
                      _buildFilterItem(context, FontAwesomeIcons.magnifyingGlass, 'Search'),
                      Divider(),
                      _buildFilterItem(context, FontAwesomeIcons.snake, 'Channel'),
                      Divider(),
                      _buildFilterItem(context, FontAwesomeIcons.signature, 'Sender'),
                      Divider(),
                      _buildFilterItem(context, FontAwesomeIcons.timer, 'Time'),
                      Divider(),
                      _buildFilterItem(context, FontAwesomeIcons.bolt, 'Priority'),
                      Divider(),
                      _buildFilterItem(context, FontAwesomeIcons.gearCode, 'Key'),
                      SizedBox(height: 4),
                    ],
                  ),
                ),
              ),
            ),
          ),
          Expanded(child: GestureDetector(child: Container(width: vpWidth, color: Color(0x88000000)), onTap: () => Navi.popDialog(context))),
        ],
      ),
    );
  }

  Widget _buildFilterItem(BuildContext context, IconData icon, String label) {
    return ListTile(
      visualDensity: VisualDensity.compact,
      title: Text(label),
      leading: Icon(icon),
      onTap: () {
        Navi.popDialog(context);
        //TOOD show more...
      },
    );
  }
}
