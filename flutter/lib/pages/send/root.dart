import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:simplecloudnotifier/state/application_log.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:simplecloudnotifier/state/user_account.dart';

class SendRootPage extends StatefulWidget {
  const SendRootPage({super.key});

  @override
  State<SendRootPage> createState() => _SendRootPageState();
}

class _SendRootPageState extends State<SendRootPage> {
  late TextEditingController _msgTitle;
  late TextEditingController _msgContent;

  @override
  void initState() {
    super.initState();
    _msgTitle = TextEditingController();
    _msgContent = TextEditingController();
  }

  @override
  void dispose() {
    _msgTitle.dispose();
    _msgContent.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<UserAccount>(
      builder: (context, acc, child) {
        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _buildQRCode(context, acc),
              const SizedBox(height: 16),
              FractionallySizedBox(
                widthFactor: 1.0,
                child: TextField(
                  controller: _msgTitle,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Title',
                  ),
                ),
              ),
              const SizedBox(height: 16),
              FractionallySizedBox(
                widthFactor: 1.0,
                child: TextField(
                  controller: _msgContent,
                  decoration: const InputDecoration(
                    border: OutlineInputBorder(),
                    labelText: 'Text',
                  ),
                ),
              ),
              const SizedBox(height: 16),
              FilledButton(
                style: FilledButton.styleFrom(textStyle: const TextStyle(fontSize: 20)),
                onPressed: _send,
                child: const Text('Send'),
              ),
            ],
          ),
        );
      },
    );
  }

  void _send() {
    //...
  }

  Widget _buildQRCode(BuildContext context, UserAccount acc) {
    if (acc.auth == null) {
      return const Placeholder();
    }

    if (acc.user == null) {
      return FutureBuilder(
        future: acc.loadUser(false),
        builder: ((context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            if (snapshot.hasError) {
              return Text('Error: ${snapshot.error}'); //TODO better error display
            }
            var url = 'https://simplecloudnotifier.de?preset_user_id=${acc.user!.userID}&preset_user_key=${acc.auth!.tokenSend}';
            return GestureDetector(
              onTap: () {
                _openWeb(url);
              },
              child: QrImageView(
                data: url,
                version: QrVersions.auto,
                size: 400.0,
                eyeStyle: QrEyeStyle(
                  eyeShape: QrEyeShape.square,
                  color: Theme.of(context).textTheme.bodyLarge?.color,
                ),
                dataModuleStyle: QrDataModuleStyle(
                  dataModuleShape: QrDataModuleShape.square,
                  color: Theme.of(context).textTheme.bodyLarge?.color,
                ),
              ),
            );
          }
          return const SizedBox(
            width: 400.0,
            height: 400.0,
            child: Center(child: CircularProgressIndicator()),
          );
        }),
      );
    }

    var url = 'https://simplecloudnotifier.de?preset_user_id=${acc.user!.userID}&preset_user_key=${acc.auth!.tokenSend}';

    return GestureDetector(
      onTap: () {
        _openWeb(url);
      },
      child: QrImageView(
        data: url,
        version: QrVersions.auto,
        size: 400.0,
        eyeStyle: QrEyeStyle(
          eyeShape: QrEyeShape.square,
          color: Theme.of(context).textTheme.bodyLarge?.color,
        ),
        dataModuleStyle: QrDataModuleStyle(
          dataModuleShape: QrDataModuleShape.square,
          color: Theme.of(context).textTheme.bodyLarge?.color,
        ),
      ),
    );
  }

  void _openWeb(String url) async {
    try {
      final Uri uri = Uri.parse(url);

      if (await canLaunchUrl(uri)) {
        await launchUrl(uri);
      } else {
        // TODO ("Cannot open URL");
      }
    } catch (exc, trace) {
      ApplicationLog.error('Failed to open URL: ' + exc.toString(), additional: 'URL: ${url}', trace: trace);
    }
  }
}
