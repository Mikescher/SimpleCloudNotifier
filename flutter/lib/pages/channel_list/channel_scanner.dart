import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:simplecloudnotifier/components/layout/scaffold.dart';
import 'package:simplecloudnotifier/models/scan_result.dart';
import 'package:simplecloudnotifier/utils/ui.dart';

class ChannelScannerPage extends StatefulWidget {
  const ChannelScannerPage({super.key});

  @override
  State<ChannelScannerPage> createState() => _ChannelScannerPageState();
}

class _ChannelScannerPageState extends State<ChannelScannerPage> {
  final MobileScannerController _controller = MobileScannerController(
    formats: const [BarcodeFormat.qrCode],
  );

  ScanResult? scanResult = null;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return SCNScaffold(
      title: "Scanner",
      showSearch: false,
      showShare: false,
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(24, 16, 24, 16),
          child: Column(
            children: [
              SizedBox(height: 16),
              Center(
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(UI.DefaultBorderRadius),
                  ),
                  clipBehavior: Clip.hardEdge,
                  child: SizedBox(
                    height: 300,
                    width: 300,
                    child: MobileScanner(
                      fit: BoxFit.cover,
                      controller: _controller,
                      onDetect: _handleBarcode,
                    ),
                  ),
                ),
              ),
              SizedBox(height: 16),
              _buildScanResult(context),
            ],
          ),
        ),
      ),
    );
  }

  void _handleBarcode(BarcodeCapture barcodes) {
    setState(() {
      if (barcodes.barcodes.isEmpty) {
        scanResult = null;
      } else {
        print('parsed: ${barcodes.barcodes[0].rawValue}');
        scanResult = ScanResult.parse(barcodes.barcodes[0].rawValue ?? '');
      }
    });
  }

  Widget _buildScanResult(BuildContext context) {
    if (scanResult == null) {
      return UI.box(
        padding: EdgeInsets.fromLTRB(16, 2, 4, 2), //TODO
        context: context,
        child: Center(
          child: Icon(FontAwesomeIcons.solidEmptySet, size: 64, color: Theme.of(context).colorScheme.onPrimaryContainer.withAlpha(128)),
        ),
      );
    }

    if (scanResult! is ScanResultMessageSend) {
      return UI.box(
        padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
        context: context,
        child: Text("TODO -- ScanResultMessageSend"), //TODO
      );
    }

    if (scanResult! is ScanResultChannel) {
      return UI.box(
        padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
        context: context,
        child: Text("TODO -- ScanResultChannel"), //TODO
      );
    }

    if (scanResult! is ScanResultChannelSubscribe) {
      return UI.box(
        padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
        context: context,
        child: Text("TODO -- ScanResultChannelSubscribe"), //TODO
      );
    }

    return UI.box(
      padding: EdgeInsets.fromLTRB(16, 2, 4, 2),
      context: context,
      child: Text("TODO -- ERROR"), //TODO
    );
  }
}
