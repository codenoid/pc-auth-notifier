import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'model/log.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class Logs extends StatefulWidget {
  @override
  _LogsState createState() => _LogsState();
}

class _LogsState extends State<Logs> {
  var logs = [];

  SharedPreferences sharedPrefs;
  String machineID = "";

  @override
  void initState() {
    super.initState();
    SharedPreferences.getInstance().then((prefs) {
      setState(() {
        sharedPrefs = prefs;
        machineID = prefs.getString('machine_id') ?? "";
        fetchLogs(machineID).then((value) {
          setState(() {
            if (!value.error) {
              for (var log in value.logs) {
                logs.add(log.raw);
              }
            }
          });
        });
      });
    });
  }

  Future<LogModel> fetchLogs(mid) async {
    final response =
        await http.get('http://192.168.8.117:8080/notification/list?id=$mid');

    if (response.statusCode == 200) {
      // If the server did return a 200 OK response,
      // then parse the JSON.
      return LogModel.fromJson(jsonDecode(response.body));
    } else {
      // If the server did not return a 200 OK response,
      // then throw an exception.
      throw Exception('Failed to load album');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      child: ListView.separated(
        separatorBuilder: (context, index) {
          return Divider(
            color: Colors.grey,
          );
        },
        itemBuilder: (context, index) {
          return Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(logs[index]),
          );
        },
        itemCount: logs.length,
      ),
    );
  }
}
