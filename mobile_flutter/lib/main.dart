import 'package:flutter/material.dart';
import 'logs_widget.dart';
import 'info_widget.dart';
import 'package:onesignal_flutter/onesignal_flutter.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  int _currentIndex = 0;
  final List<Widget> _children = [
    Logs(),
    Info(),
  ];

  void onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    bool _requireConsent = true;

    OneSignal.shared.init("<one-signal-app-id>", iOSSettings: {
      OSiOSSettings.autoPrompt: false,
      OSiOSSettings.inAppLaunchUrl: false
    });
    OneSignal.shared
        .setInFocusDisplayType(OSNotificationDisplayType.notification);

    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: Text("PC Auth Notifier"),
        ),
        body: _children[_currentIndex],
        bottomNavigationBar: BottomNavigationBar(
          onTap: onTabTapped, // new
          currentIndex: _currentIndex, // new
          items: [
            BottomNavigationBarItem(
              icon: new Icon(Icons.list),
              label: "Logs",
            ),
            BottomNavigationBarItem(
              icon: new Icon(Icons.computer),
              label: "Info",
            ),
          ],
        ),
      ),
    );
  }
}
