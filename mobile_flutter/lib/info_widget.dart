import 'package:flutter/material.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';
import 'dart:io' show Platform;
import 'package:shared_preferences/shared_preferences.dart';

class Info extends StatefulWidget {
  @override
  _InfoState createState() => _InfoState();
}

class _InfoState extends State<Info> {
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');
  Barcode result;
  QRViewController controller;
  // obtain shared preferences
  SharedPreferences sharedPrefs;
  String machineID = "";
  bool startScan = false;

  @override
  void initState() {
    super.initState();
    SharedPreferences.getInstance().then((prefs) {
      setState(() {
        sharedPrefs = prefs;
        machineID = prefs.getString('machine_id') ?? "";
      });
    });
  }

  // In order to get hot reload to work we need to pause the camera if the platform
  // is android, or resume the camera if the platform is iOS.
  @override
  void reassemble() {
    super.reassemble();
    if (Platform.isAndroid) {
      try {
        controller.pauseCamera();
      } catch (error) {}
    } else if (Platform.isIOS) {
      controller.resumeCamera();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          Text(
            machineID,
            style: TextStyle(fontSize: 20),
          ),
          if (startScan)
            Expanded(
              flex: 5,
              child: QRView(
                key: qrKey,
                onQRViewCreated: _onQRViewCreated,
              ),
            ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              RaisedButton(
                child: startScan == false
                    ? Text("SCAN QR-CODE")
                    : Text("CLOSE SCANNER"),
                onPressed: () {
                  setState(() {
                    if (startScan) {
                      startScan = false;
                    } else {
                      startScan = true;
                    }
                  });
                },
              ),
            ],
          )
        ],
      ),
    );
  }

  void _onQRViewCreated(QRViewController controller) {
    this.controller = controller;
    controller.scannedDataStream.listen((scanData) {
      setState(() {
        result = scanData;
        machineID = result.code;
        startScan = false;
        sharedPrefs.setString("machine_id", machineID);
      });
    });
  }

  @override
  void dispose() {
    controller?.dispose();
    super.dispose();
  }
}
